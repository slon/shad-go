package dist

import (
	"gitlab.com/slon/shad-go/distbuild/pkg/build"
	"gitlab.com/slon/shad-go/distbuild/pkg/proto"
)

type Cluster struct {
	sourceFiles map[build.ID]map[proto.WorkerID]struct{}
	artifacts   map[build.ID]map[proto.WorkerID]struct{}
}

func (c *Cluster) FindOptimalWorkers(task build.ID, sources, deps []build.ID) []proto.WorkerID {

}
