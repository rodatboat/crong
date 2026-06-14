package repositories

import (
	"github.com/rodatboat/crong/internal/entities"
	"gorm.io/gorm"
)

type FolderRepository struct {
	db *gorm.DB
}

func NewFolderRepository(db *gorm.DB) *FolderRepository {
	return &FolderRepository{db: db}
}

func (r *FolderRepository) FindByUser(userID uint) ([]*entities.Folder, error) {
	var folders []*entities.Folder
	if err := r.db.Where("user_id = ?", userID).Find(&folders).Error; err != nil {
		return nil, err
	}

	return folders, nil
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

func (r *FolderRepository) Update(id uint, name string) (*entities.Folder, error) {
	var folder entities.Folder
	if err := r.db.First(&folder, id).Error; err != nil {
		return nil, err
	}

	folder.Name = name

	if err := r.db.Save(&folder).Error; err != nil {
		return nil, err
	}

	return &folder, nil
}

func (r *FolderRepository) Delete(id uint) error {
	if err := r.db.Delete(&entities.Folder{}, id).Error; err != nil {
		return err
	}

	return nil
}
