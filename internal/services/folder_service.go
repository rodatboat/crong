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

// TODO: Add folder service methods
