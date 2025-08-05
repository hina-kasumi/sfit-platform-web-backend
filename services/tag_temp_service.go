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

// func (s *TagTempService) GetTagTempsByCourseID(courseID uuid.UUID) ([]entities.TagTemp, error) {
// 	return s.tagTemp_repo.GetTagTempsByCourseID(courseID)
// }

// func (s *TagTempService) DeleteTagTemp(tagID string, courseID uuid.UUID) error {
// 	return s.tagTemp_repo.DeleteTagTemp(tagID, courseID)
// }

func (s *TagTempService) UpdateTagTemp(courseID string, tags []entities.Tag) error {
	return s.tagTemp_repo.UpdateTagTemp(courseID, tags)
}