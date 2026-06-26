package repositories

import (
	"github.com/rodatboat/crong/internal/entities"
	"gorm.io/gorm"
)

type JobRepository struct {
	db           *gorm.DB
	scheduleRepo *ScheduleRepository
}

func NewJobRepository(db *gorm.DB, scheduleRepo *ScheduleRepository) *JobRepository {
	return &JobRepository{
		db:           db,
		scheduleRepo: scheduleRepo,
	}
}

func (r *JobRepository) ListByUser(userID uint) ([]*entities.Job, error) {
	var jobs []*entities.Job
	if err := r.db.Where("user_id = ?", userID).Find(&jobs).Error; err != nil {
		return nil, err
	}

	return jobs, nil
}

func (r *JobRepository) FindByJobID(jobID uint, userID uint) (*entities.Job, error) {
	var job entities.Job
	if err := r.db.Where("id = ? AND user_id = ?", jobID, userID).First(&job).Error; err != nil {
		return nil, err
	}

	return &job, nil
}

// WithTransaction runs a callback function within a database transaction
// This allows services to orchestrate multiple repository operations atomically
func (r *JobRepository) WithTransaction(fn func(*gorm.DB) error) error {
	return r.db.Transaction(fn)
}

// Create saves a single job without schedule
func (r *JobRepository) Create(tx *gorm.DB, job *entities.Job) error {
	return tx.Create(job).Error
}

func (r *JobRepository) Update(tx *gorm.DB, job *entities.Job) error {
	return tx.Save(job).Error
}

func (r *JobRepository) Delete(jobID uint) error {
	return r.db.Delete(&entities.Job{}, jobID).Error
}
