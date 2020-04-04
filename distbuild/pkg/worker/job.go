package worker

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"go.uber.org/zap"

	"gitlab.com/slon/shad-go/distbuild/pkg/api"
	"gitlab.com/slon/shad-go/distbuild/pkg/artifact"
	"gitlab.com/slon/shad-go/distbuild/pkg/build"
)

const (
	outputDirName    = "output"
	srcDirName       = "src"
	exitCodeFileName = "exit_code"
	stdoutFileName   = "stdout"
	stderrFileName   = "stderr"
)

func (w *Worker) getJobFromCache(jobID build.ID) (*api.JobResult, error) {
	aRoot, unlock, err := w.artifacts.Get(jobID)
	if err != nil {
		return nil, err
	}
	defer unlock()

	res := &api.JobResult{
		ID: jobID,
	}

	exitCodeStr, err := ioutil.ReadFile(filepath.Join(aRoot, exitCodeFileName))
	if err != nil {
		return nil, err
	}

	res.ExitCode, err = strconv.Atoi(string(exitCodeStr))
	if err != nil {
		return nil, err
	}

	res.Stdout, err = ioutil.ReadFile(filepath.Join(aRoot, stdoutFileName))
	if err != nil {
		return nil, err
	}

	res.Stderr, err = ioutil.ReadFile(filepath.Join(aRoot, stderrFileName))
	if err != nil {
		return nil, err
	}

	return res, nil
}

func executeCmd(ctx context.Context, cmd *build.Cmd) (stdout, stderr []byte, exitCode int, err error) {
	var stdoutBuf, stderrBuf bytes.Buffer

	if cmd.CatOutput != "" {
		err = ioutil.WriteFile(cmd.CatOutput, []byte(cmd.CatTemplate), 0666)
		return
	}

	p := exec.CommandContext(ctx, cmd.Exec[0], cmd.Exec[1:]...)
	p.Dir = cmd.WorkingDirectory
	p.Env = cmd.Environ
	p.Stdout = &stdoutBuf
	p.Stderr = &stderrBuf

	err = p.Run()

	stdout = stdoutBuf.Bytes()
	stderr = stderrBuf.Bytes()

	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			exitCode = exitErr.ExitCode()
			err = nil
		}
	}
	return
}

func (w *Worker) prepareSourceDir(sourceDir string, sourceFiles map[build.ID]string) (unlockSources func(), err error) {
	var unlocks []func()
	doUnlock := func() {
		for _, u := range unlocks {
			u()
		}
	}

	defer func() {
		if doUnlock != nil {
			doUnlock()
		}
	}()

	for id, path := range sourceFiles {
		dir, _ := filepath.Split(path)
		if dir != "" {
			if err := os.MkdirAll(filepath.Join(sourceDir, dir), 0777); err != nil {
				return nil, err
			}
		}

		cached, unlock, err := w.fileCache.Get(id)
		if err != nil {
			return nil, err
		}
		unlocks = append(unlocks, unlock)

		if err := os.Link(cached, filepath.Join(sourceDir, path)); err != nil {
			return nil, err
		}
	}

	unlockSources = doUnlock
	doUnlock = nil
	return
}

func (w *Worker) lockDeps(deps []build.ID) (paths map[build.ID]string, unlockDeps func(), err error) {
	var unlocks []func()
	doUnlock := func() {
		for _, u := range unlocks {
			u()
		}
	}

	defer func() {
		if doUnlock != nil {
			doUnlock()
		}
	}()

	paths = make(map[build.ID]string)

	for _, id := range deps {
		path, unlock, err := w.artifacts.Get(id)
		if err != nil {
			return nil, nil, err
		}
		unlocks = append(unlocks, unlock)

		paths[id] = filepath.Join(path, outputDirName)
	}

	unlockDeps = doUnlock
	doUnlock = nil
	return
}

func (w *Worker) runJob(ctx context.Context, spec *api.JobSpec) (*api.JobResult, error) {
	res, err := w.getJobFromCache(spec.Job.ID)
	if err != nil && !errors.Is(err, artifact.ErrNotFound) {
		return nil, err
	} else if err == nil {
		return res, nil
	}

	if err = w.pullFiles(ctx, spec.SourceFiles); err != nil {
		return nil, err
	}

	aRoot, commit, abort, err := w.artifacts.Create(spec.Job.ID)
	if err != nil {
		return nil, err
	}

	defer func() {
		if abort == nil {
			return
		}

		if err = abort(); err != nil {
			w.log.Warn("error aborting job", zap.Any("job_id", spec.Job.ID), zap.Error(err))
		}
	}()

	outputDir := filepath.Join(aRoot, outputDirName)
	if err = os.Mkdir(outputDir, 0777); err != nil {
		return nil, err
	}

	sourceDir := filepath.Join(aRoot, srcDirName)
	if err = os.Mkdir(sourceDir, 0777); err != nil {
		return nil, err
	}

	stdoutFile, err := os.Create(filepath.Join(aRoot, stdoutFileName))
	if err != nil {
		return nil, err
	}
	defer stdoutFile.Close()

	stderrFile, err := os.Create(filepath.Join(aRoot, stderrFileName))
	if err != nil {
		return nil, err
	}
	defer stderrFile.Close()

	jobContext := build.JobContext{
		OutputDir: outputDir,
		SourceDir: sourceDir,
	}

	var unlock []func()
	defer func() {
		for _, u := range unlock {
			u()
		}
	}()

	unlockSourceFiles, err := w.prepareSourceDir(sourceDir, spec.SourceFiles)
	if err != nil {
		return nil, err
	}
	unlock = append(unlock, unlockSourceFiles)

	deps, unlockDeps, err := w.lockDeps(spec.Job.Deps)
	if err != nil {
		return nil, err
	}
	unlock = append(unlock, unlockDeps)
	jobContext.Deps = deps

	res = &api.JobResult{
		ID: spec.Job.ID,
	}

	for _, cmd := range spec.Job.Cmds {
		cmd, err := cmd.Render(jobContext)
		if err != nil {
			return nil, err
		}

		stdout, stderr, exitCode, err := executeCmd(ctx, cmd)
		if err != nil {
			return nil, err
		}

		res.Stdout = append(res.Stdout, stdout...)
		_, err = stdoutFile.Write(stdout)
		if err != nil {
			return nil, err
		}

		res.Stderr = append(res.Stderr, stderr...)
		_, err = stderrFile.Write(stderr)
		if err != nil {
			return nil, err
		}

		if exitCode != 0 {
			res.ExitCode = exitCode
			break
		}
	}

	if err := ioutil.WriteFile(filepath.Join(aRoot, exitCodeFileName), []byte(strconv.Itoa(res.ExitCode)), 0666); err != nil {
		return nil, err
	}

	abort = nil
	if err := commit(); err != nil {
		return nil, err
	}

	return res, nil
}
