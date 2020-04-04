package artifact

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"gitlab.com/slon/shad-go/distbuild/pkg/build"
	"gitlab.com/slon/shad-go/distbuild/pkg/tarstream"
)

// Download artifact from remote cache into local cache.
func Download(ctx context.Context, endpoint string, c *Cache, artifactID build.ID) error {
	dir, commit, abort, err := c.Create(artifactID)
	if err != nil {
		return err
	}
	defer func() { _ = abort() }()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint+"/artifact?id="+artifactID.String(), nil)
	if err != nil {
		return err
	}

	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != http.StatusOK {
		errStr, _ := ioutil.ReadAll(rsp.Body)
		return fmt.Errorf("download: %s", errStr)
	}

	if err := tarstream.Receive(dir, rsp.Body); err != nil {
		return err
	}

	return commit()
}
