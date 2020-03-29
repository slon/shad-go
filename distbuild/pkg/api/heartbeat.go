package api

import (
	"context"
	"net/http"

	"go.uber.org/zap"

	"gitlab.com/slon/shad-go/distbuild/pkg/build"
)

// JobResult описывает результат работы джоба.
type JobResult struct {
	ID build.ID

	Stdout, Stderr []byte

	ExitCode int

	// Error описывает сообщение об ошибке, из-за которого джоб не удалось выполнить.
	//
	// Если Error == nil, значит джоб завершился успешно.
	Error *string
}

type WorkerID string

func (w WorkerID) String() string {
	return string(w)
}

type HeartbeatRequest struct {
	// WorkerID задаёт персистентный идентификатор данного воркера.
	//
	// WorkerID так же выступает в качестве endpoint-а, к которому можно подключиться по HTTP.
	//
	// В наших тестов, идентификатор будет иметь вид "localhost:%d".
	WorkerID WorkerID

	// RunningJobs перечисляет список джобов, которые выполняются на этом воркере
	// в данный момент.
	RunningJobs []build.ID

	DownloadingSources []build.ID

	DownloadingArtifacts []build.ID

	// FreeSlots сообщаяет, сколько еще процессов можно запустить на этом воркере.
	FreeSlots int

	// JobResult сообщает координатору, какие джобы завершили исполнение на этом воркере
	// на этой итерации цикла.
	FinishedJob []JobResult

	// AddedArtifacts говорит, какие артефакты появились в кеше на этой итерации цикла.
	AddedArtifacts []build.ID

	// AddedSourceFiles говорит, какие файлы появились в кеше на этой итерации цикла.
	AddedSourceFiles []build.ID
}

// JobSpec описывает джоб, который нужно запустить.
type JobSpec struct {
	SourceFiles map[build.ID]string

	Job build.Job
}

type HeartbeatResponse struct {
	JobsToRun map[build.ID]JobSpec
}

type HeartbeatService interface {
	Heartbeat(ctx context.Context, req *HeartbeatRequest) (*HeartbeatResponse, error)
}

type HeartbeatClient struct {
	Endpoint string
}

func (c *HeartbeatClient) Heartbeat(ctx context.Context, req *HeartbeatRequest) (*HeartbeatResponse, error) {
	panic("implement me")
}

type HeartbeatHandler struct {
	l *zap.Logger
	s HeartbeatService
}

func (h *HeartbeatHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/heartbeat", h.heartbeat)
}

func (h *HeartbeatHandler) heartbeat(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}
