package scheduler

import (
	"context"
	"log"
	"sync"

	"github.com/rodatboat/crong/internal/entities"
	"github.com/rodatboat/crong/internal/services"
)

// WorkerPool maintains a fixed number of goroutines that consume jobs from the
// queue and execute them via JobExecutionService.
type WorkerPool struct {
	jobQueue    <-chan entities.Job
	execService *services.JobExecutionService
	workerCount int
}

func NewWorkerPool(
	jobQueue <-chan entities.Job,
	execService *services.JobExecutionService,
	workerCount int,
) *WorkerPool {
	return &WorkerPool{
		jobQueue:    jobQueue,
		execService: execService,
		workerCount: workerCount,
	}
}

// Start launches all workers and blocks until they all exit. Workers stop when
// ctx is cancelled or when the job queue channel is closed.
func (wp *WorkerPool) Start(ctx context.Context) {
	var wg sync.WaitGroup

	for i := range wp.workerCount {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			wp.work(ctx, id)
		}(i)
	}

	wg.Wait()
	log.Println("Worker pool: stopped")
}

func (wp *WorkerPool) work(ctx context.Context, workerID int) {
	for {
		select {
		case <-ctx.Done():
			return
		case job, ok := <-wp.jobQueue:
			if !ok {
				return
			}
			log.Printf("Worker %d: executing job %d", workerID, job.ID)
			if _, err := wp.execService.ExecuteJob(job); err != nil {
				log.Printf("Worker %d: job %d failed: %v", workerID, job.ID, err)
			}
		}
	}
}
