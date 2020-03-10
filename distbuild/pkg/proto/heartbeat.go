package proto

import (
	"gitlab.com/slon/shad-go/distbuild/pkg/build"
)

// CompleteJob описывает результат работы джоба.
type FinishedJob struct {
	ID build.ID

	Stdout, Stderr []byte

	ExitCode int

	// Error описывает сообщение об ошибке, из-за которого джоб не удалось выполнить.
	//
	// Если Error == nil, значит джоб завершился успешно.
	Error *string
}

type HeartbeatRequest struct {
	// WorkerID задаёт персистентный идентификатор данного воркера.
	//
	// WorkerID так же выступает в качестве endpoint-а, к которому можно подключиться по HTTP.
	//
	// В наших тестов, идентификатор будет иметь вид "localhost:%d".
	WorkerID string

	// ProcessID задаёт эфемерный идентификатор текущего процесса воркера.
	//
	// Координатор запоминает ProcessID для каждого воркера.
	//
	// Измение ProcessID значит, что воркер перезапустился.
	ProcessID string

	// RunningJobs перечисляет список джобов, которые выполняются на этом воркере
	// в данный момент.
	RunningJobs []build.ID

	DownloadingSources []build.ID

	DownloadingArtifacts []build.ID

	// FreeSlots сообщаяет, сколько еще процессов можно запустить на этом воркере.
	FreeSlots int

	// FinishedJob сообщает координатору, какие джобы завершили исполнение на этом воркере
	// на этой итерации цикла.
	FinishedJob []FinishedJob

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

// ArtifactSpec описывает артефакт, который нужно скачать с другого воркера.
type ArtifactSpec struct {
}

// SourceFileSpec описывает файл с исходным кодом, который нужно скачать с координатора.
type SourceFileSpec struct {
}

type HeartbeatResponse struct {
	JobsToRun map[build.ID]JobSpec

	ArtifactsToDownload map[build.ID]ArtifactSpec

	ArtifactsToRemove []build.ID

	SourceFilesToDownload map[build.ID]SourceFileSpec

	SourceFilesToRemove []build.ID
}
