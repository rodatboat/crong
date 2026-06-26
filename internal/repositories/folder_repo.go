package repositories

import (
	"time"

	"github.com/rodatboat/crong/internal/entities"
	"gorm.io/gorm"
)

type FolderRepository struct {
	db *gorm.DB
}

func NewFolderRepository(db *gorm.DB) *FolderRepository {
	return &FolderRepository{db: db}
}

func (r *FolderRepository) ListByUser(userID uint) ([]*entities.Folder, error) {
	var folders []*entities.Folder
	if err := r.db.Where("user_id = ?", userID).Find(&folders).Error; err != nil {
		return nil, err
	}

	return folders, nil
}

func (r *FolderRepository) FindByFolderIDAndUserID(folderID uint, userID uint) (*entities.Folder, error) {
	var folder *entities.Folder
	if err := r.db.Where("id = ? AND user_id = ?", folderID, userID).First(&folder).Error; err != nil {
		return nil, err
	}

	return folder, nil
}

func (r *FolderRepository) Create(name string, userID uint) (*entities.Folder, error) {
	folder := &entities.Folder{
		Name:   name,
		UserID: userID,
	}

	if err := r.db.Create(folder).Error; err != nil {
		return nil, err
	}

	return folder, nil
}

func (r *FolderRepository) Update(folderID uint, userID uint, name string) (*entities.Folder, error) {
	var folderEntity entities.Folder
	if err := r.db.Where("id = ? AND user_id = ?", folderID, userID).First(&folderEntity).Error; err != nil {
		return nil, err
	}

	folderEntity.Name = name
	folderEntity.UpdatedAt = time.Now()

	if err := r.db.Save(&folderEntity).Error; err != nil {
		return nil, err
	}

	return &folderEntity, nil
}

func (r *FolderRepository) Delete(folderID uint, userID uint) error {
	var folder entities.Folder
	if err := r.db.Where("id = ? AND user_id = ?", folderID, userID).First(&folder).Error; err != nil {
		return err
	}

	if err := r.db.Delete(folder).Error; err != nil {
		return err
	}

	return nil
}
