package repositories

import (
	"github.com/rodatboat/crong/internal/entities"
	"gorm.io/gorm"
)

type JobExecutionRepository struct {
	db *gorm.DB
}

func NewJobExecutionRepository(db *gorm.DB) *JobExecutionRepository {
	return &JobExecutionRepository{db: db}
}

func (r *JobExecutionRepository) Create(jobExecution *entities.JobExecution) error {

	if err := r.db.Create(jobExecution).Error; err != nil {
		return err
	}

	return nil
}

func (r *JobExecutionRepository) ListByJobID(jobID uint, userID uint, limit int, offset int) ([]*entities.JobExecution, error) {
	var jobExecutions []*entities.JobExecution
	if err := r.db.Where("job_id = ? AND user_id = ?", jobID, userID).
		Order("executed_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&jobExecutions).Error; err != nil {
		return nil, err
	}

	return jobExecutions, nil
}
