package services

import (
	"github.com/gofiber/fiber/v3/log"
	"github.com/rodatboat/crong/internal/repositories"
)

type FolderService struct {
	folderRepo *repositories.FolderRepository
}

func NewFolderService(folderRepo *repositories.FolderRepository) *FolderService {
	return &FolderService{
		folderRepo: folderRepo,
	}
}

func (f *FolderService) FolderExists(folderID uint, userID uint) bool {
	log.Infof("Fetching folder with id %v for user %v", folderID, userID)
	folder, err := f.folderRepo.FindByFolderIDAndUserID(folderID, userID)
	if err != nil {
		return false
	}

	if folder != nil {
		return true
	}

	return false
}

// TODO: Add folder service methods
