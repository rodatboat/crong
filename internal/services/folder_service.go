package services

import (
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
	// TODO: Add userID to check
	folder, err := f.folderRepo.FindByFolder(folderID)
	if err != nil {
		return false
	}

	if folder != nil {
		return true
	}

	return false
}

// TODO: Add folder service methods
