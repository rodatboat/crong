package repositories

import (
	"github.com/rodatboat/crong/internal/entities"
	"gorm.io/gorm"
)

type JobRepository struct {
	db *gorm.DB
}

func NewJobRepository(db *gorm.DB) *JobRepository {
	return &JobRepository{db: db}
}

func (r *JobRepository) FindByUser(userID uint) ([]*entities.Job, error) {
	var jobs []*entities.Job
	if err := r.db.Where("user_id = ?", userID).Find(&jobs).Error; err != nil {
		return nil, err
	}

	return jobs, nil
}
