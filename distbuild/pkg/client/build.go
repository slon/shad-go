package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"gitlab.com/slon/shad-go/distbuild/pkg/build"
	"gitlab.com/slon/shad-go/distbuild/pkg/proto"
)

type Client struct {
	CoordinatorEndpoint string

	SourceDir string
}

type BuildListener interface {
	OnJobStdout(jobID build.ID, stdout []byte) error
	OnJobStderr(jobID build.ID, stdout []byte) error

	OnJobFinished(jobID build.ID) error
	OnJobFailed(jobID build.ID, code int, error string) error
}

func (c *Client) uploadSources(ctx context.Context, src proto.MissingSources) error {

}

func (c *Client) Build(ctx context.Context, graph build.Graph, lsn BuildListener) error {
	graphJS, err := json.Marshal(graph)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", c.CoordinatorEndpoint+"/build", bytes.NewBuffer(graphJS))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req = req.WithContext(ctx)

	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("build failed: %w", err)
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != 200 {
		errorMsg, _ := ioutil.ReadAll(rsp.Body)
		return fmt.Errorf("build failed: %s", errorMsg)
	}

	d := json.NewDecoder(rsp.Body)

	var missing proto.MissingSources
	if err := d.Decode(&missing); err != nil {
		return err
	}

	if err := c.uploadSources(ctx, missing); err != nil {
		return err
	}

	for {
		var update proto.StatusUpdate
		if err := d.Decode(&update); err != nil {
			return err
		}

		switch {
		case update.BuildFailed != nil:
			return fmt.Errorf("build failed: %s", update.BuildFailed.Error)

		case update.JobFinished != nil:
			jf := update.JobFinished

			if jf.Stdout != nil {
				if err := lsn.OnJobStdout(jf.ID, jf.Stdout); err != nil {
					return err
				}
			}

			if jf.Stderr != nil {
				if err := lsn.OnJobStderr(jf.ID, jf.Stderr); err != nil {
					return err
				}
			}

			if jf.Error != nil {
				if err := lsn.OnJobFailed(jf.ID, jf.ExitCode, *jf.Error); err != nil {
					return err
				}
			} else {
				if err := lsn.OnJobFinished(jf.ID); err != nil {
					return err
				}
			}

		default:
			return fmt.Errorf("build failed: unexpected status update")
		}
	}
}
