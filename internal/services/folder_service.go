package services

import (
	"errors"

	"github.com/gofiber/fiber/v3/log"
	"github.com/rodatboat/crong/internal/entities"
	"github.com/rodatboat/crong/internal/models"
	"github.com/rodatboat/crong/internal/repositories"
	"github.com/rodatboat/crong/internal/resp"
	"gorm.io/gorm"
)

type FolderService struct {
	folderRepo *repositories.FolderRepository
}

func NewFolderService(folderRepo *repositories.FolderRepository) *FolderService {
	return &FolderService{
		folderRepo: folderRepo,
	}
}

func (f *FolderService) FolderExists(folderID uint, userID uint) (bool, error) {
	log.Infof("Fetching folder with id %v for user %v", folderID, userID)
	folder, err := f.folderRepo.FindByFolderIDAndUserID(folderID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Errorf("Folder %v not found for user %v", folderID, userID)
			return false, nil
		}
		return false, err
	}

	if folder != nil {
		return true, nil
	}

	return false, nil
}

func (f *FolderService) GetFoldersByUser(userID uint) ([]*models.Folder, error) {
	log.Infof("Listing folders for user %v", userID)

	folderEntities, err := f.folderRepo.ListByUser(userID)
	if err != nil {
		return nil, err
	}

	folders := make([]*models.Folder, len(folderEntities))
	for idx, folderEntity := range folderEntities {
		folders[idx] = f.mapFolderEntityToFolderModel(folderEntity)
	}
	return folders, nil
}

func (f *FolderService) GetFolderDetailsByID(folderID uint, userID uint) (*models.Folder, error) {
	log.Infof("Fetching folder with id %v for user %v", folderID, userID)

	folderEntity, err := f.folderRepo.FindByFolderIDAndUserID(folderID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, resp.ErrNotFound
		}
		return nil, err
	}

	return f.mapFolderEntityToFolderModel(folderEntity), nil
}

func (f *FolderService) CreateFolder(userID uint, req *models.FolderCreate) (*models.Folder, error) {
	log.Infof("Creating folder for user %v with payload %+v", userID, req)

	// Create folder
	folderEntity, err := f.folderRepo.Create(req.Name, userID)
	if err != nil {
		return nil, err
	}

	return f.mapFolderEntityToFolderModel(folderEntity), nil
}

func (f *FolderService) UpdateFolder(folderID uint, userID uint, req *models.FolderUpdate) (*models.Folder, error) {
	log.Infof("Updating folder with id %v for user %v with payload %+v", folderID, userID, req)

	// Update folder
	folderEntity, err := f.folderRepo.Update(folderID, userID, req.Name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, resp.ErrNotFound
		}
		return nil, err
	}

	return f.mapFolderEntityToFolderModel(folderEntity), nil
}

func (f *FolderService) DeleteFolder(folderID uint, userID uint) error {
	log.Infof("Deleting folder with id %v for user %v", folderID, userID)

	// Delete folder
	if err := f.folderRepo.Delete(folderID, userID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return resp.ErrNotFound
		}
		return err
	}

	return nil
}

// ========== UTILITIES ==========

func (f *FolderService) mapFolderEntityToFolderModel(folderEntity *entities.Folder) *models.Folder {
	return &models.Folder{
		ID:        folderEntity.ID,
		Name:      folderEntity.Name,
		UserID:    folderEntity.UserID,
		CreatedAt: folderEntity.CreatedAt,
		UpdatedAt: folderEntity.UpdatedAt,
	}
}
