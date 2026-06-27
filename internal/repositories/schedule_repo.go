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

func (r *ScheduleRepository) ListJobsBySchedule(minute uint, hour uint, mday uint, month uint, wday uint) ([]*entities.Job, error) {
	var jobs []*entities.Job

	err := r.db.
		Joins("INNER JOIN schedule_hours sh ON jobs.id = sh.job_id AND (sh.hour = ? OR sh.hour = -1)", hour).
		Joins("INNER JOIN schedule_minutes sm ON jobs.id = sm.job_id AND (sm.minute = ? OR sm.minute = -1)", minute).
		Joins("INNER JOIN schedule_mdays smd ON jobs.id = smd.job_id AND (smd.mday = ? OR smd.mday = -1)", mday).
		Joins("INNER JOIN schedule_wdays sw ON jobs.id = sw.job_id AND (sw.wday = ? OR sw.wday = -1)", wday).
		Joins("INNER JOIN schedule_months smo ON jobs.id = smo.job_id AND (smo.month = ? OR smo.month = -1)", month).
		Where("jobs.enabled = ?", true).
		Group("jobs.id").
		Order("jobs.id ASC").
		Find(&jobs).Error

	if err != nil {
		return nil, err
	}

	return jobs, err
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
