package repositories

import (
	"github.com/rodatboat/crong/internal/entities"
	"gorm.io/gorm"
)

type ScheduleRepository struct {
	db *gorm.DB
}

func NewScheduleRepository(db *gorm.DB) *ScheduleRepository {
	return &ScheduleRepository{db: db}
}

// CreateSchedules creates all schedule entries (minute, hour, mday, wday, month) in the database
func (r *ScheduleRepository) CreateSchedules(tx *gorm.DB, schedule *entities.Schedule) error {
	// Create schedule minute entries
	if len(schedule.Minute) > 0 {
		if err := tx.CreateInBatches(schedule.Minute, 60).Error; err != nil {
			return err
		}
	}

	// Create schedule hour entries
	if len(schedule.Hour) > 0 {
		if err := tx.CreateInBatches(schedule.Hour, 24).Error; err != nil {
			return err
		}
	}

	// Create schedule mday entries
	if len(schedule.Mday) > 0 {
		if err := tx.CreateInBatches(schedule.Mday, 31).Error; err != nil {
			return err
		}
	}

	// Create schedule wday entries
	if len(schedule.Wday) > 0 {
		if err := tx.CreateInBatches(schedule.Wday, 7).Error; err != nil {
			return err
		}
	}

	// Create schedule month entries
	if len(schedule.Month) > 0 {
		if err := tx.CreateInBatches(schedule.Month, 12).Error; err != nil {
			return err
		}
	}

	return nil
}

// DeleteSchedulesByJobID deletes all schedule entries for a given job
func (r *ScheduleRepository) DeleteSchedulesByJobID(tx *gorm.DB, jobID uint) error {
	if err := tx.Where("job_id = ?", jobID).Delete(&entities.ScheduleMinute{}).Error; err != nil {
		return err
	}
	if err := tx.Where("job_id = ?", jobID).Delete(&entities.ScheduleHour{}).Error; err != nil {
		return err
	}
	if err := tx.Where("job_id = ?", jobID).Delete(&entities.ScheduleMday{}).Error; err != nil {
		return err
	}
	if err := tx.Where("job_id = ?", jobID).Delete(&entities.ScheduleWday{}).Error; err != nil {
		return err
	}
	if err := tx.Where("job_id = ?", jobID).Delete(&entities.ScheduleMonth{}).Error; err != nil {
		return err
	}
	return nil
}
