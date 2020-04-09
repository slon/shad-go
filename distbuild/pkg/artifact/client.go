// +build !solution

package artifact

import (
	"context"

	"gitlab.com/slon/shad-go/distbuild/pkg/build"
)

// Download artifact from remote cache into local cache.
func Download(ctx context.Context, endpoint string, c *Cache, artifactID build.ID) error {
	panic("implement me")
}
