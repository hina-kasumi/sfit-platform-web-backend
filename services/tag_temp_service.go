package services

import (
	"sfit-platform-web-backend/entities"
	"sfit-platform-web-backend/repositories"

	"github.com/google/uuid"
)

type TagTempService struct {
    tagTemp_repo *repositories.TagTempRepository
}

func NewTagTempService(tagTemp_repo *repositories.TagTempRepository) *TagTempService {
    return &TagTempService{tagTemp_repo: tagTemp_repo}
}

func (s *TagTempService) CreateTagTemp(tagID string, courseID uuid.UUID) (*entities.TagTemp, error) {
    return s.tagTemp_repo.CreateNewTagTemp(tagID, courseID)
}