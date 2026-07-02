package scheduler

import (
	"time"

	"github.com/gofiber/fiber/v3/log"
	"github.com/rodatboat/crong/internal/entities"
	"github.com/rodatboat/crong/internal/repositories"
)

// loadDueJobs queries the database for all distinct timezones, then for each
// timezone, converts the current UTC time to that timezone and fetches jobs
// whose schedule matches in that timezone.
func loadDueJobs(scheduleRepo *repositories.ScheduleRepository) ([]*entities.Job, error) {
	var allJobs []*entities.Job

	// Get all distinct timezones from enabled jobs
	timezones, err := scheduleRepo.ListDistinctTimezones()
	if err != nil {
		return nil, err
	}

	if len(timezones) == 0 {
		return allJobs, nil
	}

	// Process jobs for each timezone
	for _, tzName := range timezones {
		jobs, err := loadJobsForTimezone(scheduleRepo, tzName)
		if err != nil {
			log.Errorf("loadDueJobs: error loading jobs for timezone %s: %v", tzName, err)
			continue
		}
		allJobs = append(allJobs, jobs...)
	}

	return allJobs, nil
}

// loadJobsForTimezone converts UTC time to the specified timezone and queries
// for jobs whose schedule matches in that timezone.
func loadJobsForTimezone(scheduleRepo *repositories.ScheduleRepository, tzName string) ([]*entities.Job, error) {
	// Load the timezone location
	loc, err := time.LoadLocation(tzName)
	if err != nil {
		return nil, err
	}

	// Convert UTC now to the target timezone
	nowUTC := time.Now().UTC()
	nowTZ := nowUTC.In(loc)

	// Extract time components in the target timezone
	minute := uint(nowTZ.Minute())
	hour := uint(nowTZ.Hour())
	mday := uint(nowTZ.Day())
	month := uint(nowTZ.Month())
	wday := uint(nowTZ.Weekday())

	return scheduleRepo.ListJobsBySchedule(minute, hour, mday, month, wday, tzName)
}
