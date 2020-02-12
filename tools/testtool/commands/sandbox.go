package commands

import (
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"syscall"
)

func currentUserIsRoot() bool {
	return os.Getuid() == 0
}

func sandbox(cmd *exec.Cmd) error {
	nobody, err := user.Lookup("nobody")
	if err != nil {
		return err
	}

	uid, _ := strconv.Atoi(nobody.Uid)
	gid, _ := strconv.Atoi(nobody.Gid)

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Credential: &syscall.Credential{
			Uid: uint32(uid),
			Gid: uint32(gid),
		},
	}

	cmd.Env = []string{}

	return nil
}
