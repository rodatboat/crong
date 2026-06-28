package scheduler

import (
	"time"

	"github.com/rodatboat/crong/internal/entities"
	"github.com/rodatboat/crong/internal/repositories"
)

// loadDueJobs extracts the current UTC time components and queries the DB for
// all enabled jobs whose decomposed schedule matches this exact minute.
func loadDueJobs(scheduleRepo *repositories.ScheduleRepository) ([]*entities.Job, error) {
	now := time.Now().UTC()

	minute := uint(now.Minute())
	hour := uint(now.Hour())
	mday := uint(now.Day())
	month := uint(now.Month())
	wday := uint(now.Weekday())

	return scheduleRepo.ListJobsBySchedule(minute, hour, mday, month, wday)
}
