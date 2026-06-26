package services

import (
	"errors"

	"github.com/gofiber/fiber/v3/log"
	"github.com/rodatboat/crong/internal/repositories"
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

// TODO: Add folder service methods
