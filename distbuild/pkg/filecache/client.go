// +build !solution

package filecache

import (
	"context"

	"go.uber.org/zap"

	"gitlab.com/slon/shad-go/distbuild/pkg/build"
)

type Client struct {
}

func NewClient(l *zap.Logger, endpoint string) *Client {
	panic("implement me")
}

func (c *Client) Upload(ctx context.Context, id build.ID, localPath string) error {
	panic("implement me")
}

func (c *Client) Download(ctx context.Context, localCache *Cache, id build.ID) error {
	panic("implement me")
}
