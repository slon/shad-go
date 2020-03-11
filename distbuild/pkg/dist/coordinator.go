package dist

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"go.uber.org/zap"

	"gitlab.com/slon/shad-go/distbuild/pkg/build"
	"gitlab.com/slon/shad-go/distbuild/pkg/filecache"
	"gitlab.com/slon/shad-go/distbuild/pkg/proto"
)

type Build struct {
}

type Coordinator struct {
	log       *zap.Logger
	mux       *http.ServeMux
	fileCache *filecache.Cache
}

func NewCoordinator(
	log *zap.Logger,
	fileCache *filecache.Cache,
) *Coordinator {
	c := &Coordinator{
		log:       log,
		mux:       http.NewServeMux(),
		fileCache: fileCache,
	}

	c.mux.HandleFunc("/build", c.Build)
	return c
}

func (c *Coordinator) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.mux.ServeHTTP(w, r)
}

func (c *Coordinator) doBuild(w http.ResponseWriter, r *http.Request) error {
	graphJS, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	var g build.Graph
	if err := json.Unmarshal(graphJS, &g); err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	enc := json.NewEncoder(w)
	if err := enc.Encode(proto.MissingSources{}); err != nil {
		return err
	}

	return fmt.Errorf("coordinator not implemented")
}

func (c *Coordinator) Build(w http.ResponseWriter, r *http.Request) {
	if err := c.doBuild(w, r); err != nil {
		c.log.Error("build failed", zap.Error(err))

		errorUpdate := proto.StatusUpdate{BuildFailed: &proto.BuildFailed{Error: err.Error()}}
		errorJS, _ := json.Marshal(errorUpdate)
		_, _ = w.Write(errorJS)
	}
}
