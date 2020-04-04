package filecache

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"go.uber.org/zap"

	"gitlab.com/slon/shad-go/distbuild/pkg/build"
)

type Client struct {
	l        *zap.Logger
	endpoint string
}

func NewClient(l *zap.Logger, endpoint string) *Client {
	return &Client{
		l:        l,
		endpoint: endpoint,
	}
}

func (c *Client) Upload(ctx context.Context, id build.ID, localPath string) error {
	f, err := os.Open(localPath)
	if err != nil {
		return err
	}
	defer f.Close()

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, c.endpoint+"/file?id="+id.String(), f)
	if err != nil {
		return err
	}

	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != 200 {
		errStr, _ := ioutil.ReadAll(rsp.Body)
		return fmt.Errorf("file upload: %s", errStr)
	}

	return nil
}

func (c *Client) Download(ctx context.Context, localCache *Cache, id build.ID) error {
	w, abort, err := localCache.Write(id)
	if err != nil {
		return err
	}
	defer abort()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.endpoint+"/file?id="+id.String(), nil)
	if err != nil {
		return err
	}

	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != 200 {
		errStr, _ := ioutil.ReadAll(rsp.Body)
		return fmt.Errorf("file upload: %s", errStr)
	}

	_, err = io.Copy(w, rsp.Body)
	if err != nil {
		return err
	}

	return w.Close()
}
