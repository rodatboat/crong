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
