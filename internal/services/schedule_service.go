package services

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3/log"
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
// Format: "minute hour mday month wday"
// Each field supports:
//   - "*" = any value (represented as -1)
//   - Single values: "5", "10"
//   - Lists: "1,3,5" (comma-separated)
//   - Ranges: "1-5" (expands to 1,2,3,4,5)
//   - Steps: "*/N" only (e.g., "*/15", "*/5")
//   - Combined: "0,15,30-45" (lists and ranges, no steps on individual items)
//
// NOTE: Steps (/) are ONLY allowed with * (e.g., "*/5").
//
//	Steps on ranges ("1-10/2") or lists ("1,2/5") are NOT supported.
//
// Examples:
//   - "0 9 * * *" = 9:00 AM every day
//   - "*/5 * * * *" = every 5 minutes
//   - "0 9,17 * * 1-5" = 9 AM and 5 PM on weekdays only
func (s *ScheduleService) CronExpressionToSchedule(cronExpr string) (*models.Schedule, error) {
	// Split cron expression into fields. (e.g. "0 */3 * * *" -> ["0", "*/3", "*", "*", "*"])
	log.Infof("Parsing cron expression: %v", cronExpr)
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

	schedule := &models.Schedule{
		Minute: minute,
		Hour:   hour,
		Mday:   mday,
		Month:  month,
		Wday:   wday,
	}

	log.Infof("Parsed schedule: %+v", schedule)
	return schedule, nil
}

// parseField parses a single cron field and returns a slice of values
// Supported formats:
//   - "*" (any) returns [-1]
//   - "5" (single value)
//   - "1,3,5" (list)
//   - "1-5" (range, expands to 1,2,3,4,5)
//   - "1-2,5,10" (combined ranges and values)
//   - "*/N" (step, every N units)
//
// NOT supported: steps on ranges "1-10/2" or steps on lists "1,5/2"
func (s *ScheduleService) parseField(field string, min, max int, fieldName string) ([]int, error) {
	if field == "*" {
		return []int{-1}, nil
	}

	// Handle "*/N" step syntax
	if strings.HasPrefix(field, "*/") {
		step, err := strconv.Atoi(field[2:])
		if err != nil || step <= 0 {
			return nil, fmt.Errorf("invalid step value in %s: %s", fieldName, field)
		}

		var values []int
		for i := min; i <= max; i += step {
			values = append(values, i)
		}
		return values, nil
	}

	// No step syntax - parse as list of ranges/values
	var allValues []int
	seen := make(map[int]bool) // Avoid duplicates
	expressions := strings.Split(field, ",")

	for _, expr := range expressions {
		expr = strings.TrimSpace(expr)

		if expr == "" {
			return nil, fmt.Errorf("empty expression in %s", fieldName)
		}

		// Check for range
		if idx := strings.Index(expr, "-"); idx != -1 {
			// Range like "1-5"
			rangeMin, err1 := strconv.Atoi(strings.TrimSpace(expr[:idx]))
			rangeMax, err2 := strconv.Atoi(strings.TrimSpace(expr[idx+1:]))
			if err1 != nil {
				return nil, fmt.Errorf("invalid range start in %s: %s", fieldName, expr[:idx])
			}
			if err2 != nil {
				return nil, fmt.Errorf("invalid range end in %s: %s", fieldName, expr[idx+1:])
			}

			if rangeMin < min || rangeMax > max {
				return nil, fmt.Errorf("%s value out of range: %d-%d (must be %d-%d)", fieldName, rangeMin, rangeMax, min, max)
			}

			if rangeMin > rangeMax {
				return nil, fmt.Errorf("invalid range in %s: start %d > end %d", fieldName, rangeMin, rangeMax)
			}

			for i := rangeMin; i <= rangeMax; i++ {
				if !seen[i] {
					allValues = append(allValues, i)
					seen[i] = true
				}
			}
		} else {
			// Single value
			val, err := strconv.Atoi(strings.TrimSpace(expr))
			if err != nil {
				return nil, fmt.Errorf("invalid %s value: %s", fieldName, expr)
			}

			if val < min || val > max {
				return nil, fmt.Errorf("%s value out of range: %d (must be %d-%d)", fieldName, val, min, max)
			}

			if !seen[val] {
				allValues = append(allValues, val)
				seen[val] = true
			}
		}
	}

	return allValues, nil
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
