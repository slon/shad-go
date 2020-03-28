package dist

import (
	"context"

	"gitlab.com/slon/shad-go/distbuild/pkg/build"
	"gitlab.com/slon/shad-go/distbuild/pkg/proto"
)

type Build struct {
	ID    build.ID
	Graph *build.Graph

	coordinator    *Coordinator
	uploadComplete chan struct{}
}

func NewBuild(graph *build.Graph, coordinator *Coordinator) *Build {
	id := build.NewID()

	return &Build{
		ID:    id,
		Graph: graph,

		coordinator:    coordinator,
		uploadComplete: make(chan struct{}),
	}
}

func (b *Build) Run(ctx context.Context, onStatusUpdate func(update proto.StatusUpdate) error) error {
	panic("implement me")
}

func (b *Build) UploadComplete() {
	close(b.uploadComplete)
}
