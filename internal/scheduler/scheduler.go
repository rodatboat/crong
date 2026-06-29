package scheduler

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v3/log"
	"github.com/rodatboat/crong/internal/entities"
	"github.com/rodatboat/crong/internal/repositories"
)

// Scheduler queries the DB every minute for due jobs and pushes them into the
// job queue for the worker pool to consume.
type Scheduler struct {
	scheduleRepo *repositories.ScheduleRepository
	jobQueue     chan entities.Job
}

func New(scheduleRepo *repositories.ScheduleRepository, queueSize int) *Scheduler {
	return &Scheduler{
		scheduleRepo: scheduleRepo,
		jobQueue:     make(chan entities.Job, queueSize),
	}
}

// JobQueue returns a read-only view of the job channel for workers to consume.
func (s *Scheduler) JobQueue() <-chan entities.Job {
	return s.jobQueue
}

// Start runs the scheduler loop. It aligns to the next minute boundary, then
// fires every 60 seconds. Blocks until ctx is cancelled, then closes the channel.
func (s *Scheduler) Start(ctx context.Context) {
	defer close(s.jobQueue)

	// Wait until the next minute boundary before starting.
	now := time.Now()
	nextMinute := now.Truncate(time.Minute).Add(time.Minute)
	waitDuration := time.Until(nextMinute)
	log.Warnf("Scheduler: first tick in %s", waitDuration.Round(time.Second))

	select {
	case <-ctx.Done():
		return
	case <-time.After(waitDuration):
	}

	s.tick()

	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Warnf("Scheduler: stopped")
			return
		case <-ticker.C:
			s.tick()
		}
	}
}

func (s *Scheduler) tick() {
	jobs, err := loadDueJobs(s.scheduleRepo)
	if err != nil {
		log.Errorf("Scheduler: error loading due jobs: %v", err)
		return
	}

	log.Infof("Scheduler: %d job(s) due this tick", len(jobs))

	for _, job := range jobs {
		select {
		case s.jobQueue <- *job:
		default:
			log.Warnf("Scheduler: job queue full, dropping job %d", job.ID)
		}
	}

	// Check if cleanup should run
	// s.checkAndRunCleanup()
}

// checkAndRunCleanup checks if it's time to run the cleanup job based on schedule
// func (s *Scheduler) checkAndRunCleanup() {
// 	now := time.Now().UTC()

// 	// Parse the cleanup schedule (cron-like: minute hour day month weekday)
// 	if shouldRunCleanup(now, s.lastCleanupExecution, s.cleanupSchedule) {
// 		log.Infof("Scheduler: triggering cleanup job")
// 		go func() {
// 			if err := s.jobExecutionService.CleanupOldExecutions(s.cleanupRetention); err != nil {
// 				log.Errorf("Scheduler: cleanup job failed: %v", err)
// 			}
// 		}()
// 		s.lastCleanupExecution = now
// 	}
// }
