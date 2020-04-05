package api

import (
	"context"

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

	// FreeSlots сообщает, сколько еще процессов можно запустить на этом воркере.
	FreeSlots int

	// JobResult сообщает координатору, какие джобы завершили исполнение на этом воркере
	// на этой итерации цикла.
	FinishedJob []JobResult

	// AddedArtifacts говорит, какие артефакты появились в кеше на этой итерации цикла.
	AddedArtifacts []build.ID
}

// JobSpec описывает джоб, который нужно запустить.
type JobSpec struct {
	// SourceFiles задаёт список файлов, который должны присутствовать в директории с исходным кодом при запуске этого джоба.
	SourceFiles map[build.ID]string

	// Artifacts задаёт воркеров, с которых можно скачать артефакты необходимые этом джобу.
	Artifacts map[build.ID]WorkerID

	build.Job
}

type HeartbeatResponse struct {
	JobsToRun map[build.ID]JobSpec
}

type HeartbeatService interface {
	Heartbeat(ctx context.Context, req *HeartbeatRequest) (*HeartbeatResponse, error)
}
