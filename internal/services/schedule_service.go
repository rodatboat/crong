package services

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/rodatboat/crong/internal/entities"
	"github.com/rodatboat/crong/internal/models"
)

type ScheduleService struct {
	// Can add schedule-specific dependencies here if needed
}

func NewScheduleService() *ScheduleService {
	return &ScheduleService{}
}

// CronExpressionToSchedule parses a cron expression and returns a schedule model
// Cron format: "minute hour mday month wday"
// Each field can be: * (any), number (specific), or comma-separated list
// Returns models.Schedule with parsed values
func (s *ScheduleService) CronExpressionToSchedule(cronExpr string) (*models.Schedule, error) {
	// Split cron expression into fields. (e.g. "0 */3 * * *" -> ["0", "*/3", "*", "*", "*"])
	parts := strings.Fields(strings.TrimSpace(cronExpr))
	if len(parts) != 5 {
		return nil, fmt.Errorf("invalid cron expression: expected 5 fields, got %d", len(parts))
	}

	minute, err := s.parseField(parts[0], 0, 59, "minute")
	if err != nil {
		return nil, err
	}

	hour, err := s.parseField(parts[1], 0, 23, "hour")
	if err != nil {
		return nil, err
	}

	mday, err := s.parseField(parts[2], 1, 31, "mday")
	if err != nil {
		return nil, err
	}

	month, err := s.parseField(parts[3], 1, 12, "month")
	if err != nil {
		return nil, err
	}

	wday, err := s.parseField(parts[4], 0, 6, "wday")
	if err != nil {
		return nil, err
	}

	return &models.Schedule{
		Minute: minute,
		Hour:   hour,
		Mday:   mday,
		Month:  month,
		Wday:   wday,
	}, nil
}

// parseField parses a single cron field and returns a slice of values
// Returns [-1] for "*" (any), or a slice of specific values
func (s *ScheduleService) parseField(field string, min, max int, fieldName string) ([]int, error) {
	if field == "*" {
		return []int{-1}, nil
	}

	var values []int
	parts := strings.Split(field, ",")

	for _, part := range parts {
		part = strings.TrimSpace(part)
		val, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("invalid %s value: %s", fieldName, part)
		}

		if val < min || val > max {
			return nil, fmt.Errorf("%s value out of range: %d (must be %d-%d)", fieldName, val, min, max)
		}

		values = append(values, val)
	}

	return values, nil
}

// ScheduleModelToEntity converts a models.Schedule to entities.Schedule with a given job ID
func (s *ScheduleService) ScheduleModelToEntity(jobID uint, schedule *models.Schedule) *entities.Schedule {
	scheduleEntity := &entities.Schedule{}

	// Convert minute values
	if schedule.Minute != nil {
		for _, m := range schedule.Minute {
			scheduleEntity.Minute = append(scheduleEntity.Minute, entities.ScheduleMinute{
				JobID:  jobID,
				Minute: m,
			})
		}
	}

	// Convert hour values
	if schedule.Hour != nil {
		for _, h := range schedule.Hour {
			scheduleEntity.Hour = append(scheduleEntity.Hour, entities.ScheduleHour{
				JobID: jobID,
				Hour:  h,
			})
		}
	}

	// Convert mday values
	if schedule.Mday != nil {
		for _, md := range schedule.Mday {
			scheduleEntity.Mday = append(scheduleEntity.Mday, entities.ScheduleMday{
				JobID: jobID,
				Mday:  md,
			})
		}
	}

	// Convert wday values
	if schedule.Wday != nil {
		for _, w := range schedule.Wday {
			scheduleEntity.Wday = append(scheduleEntity.Wday, entities.ScheduleWday{
				JobID: jobID,
				Wday:  w,
			})
		}
	}

	// Convert month values
	if schedule.Month != nil {
		for _, mo := range schedule.Month {
			scheduleEntity.Month = append(scheduleEntity.Month, entities.ScheduleMonth{
				JobID: jobID,
				Month: mo,
			})
		}
	}

	return scheduleEntity
}
