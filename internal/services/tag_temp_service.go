package services

import (
	"sfit-platform-web-backend/internal/model"
	"sfit-platform-web-backend/internal/repositories"

	"github.com/google/uuid"
)

type TagTempService struct {
	tagTemp_repo *repositories.TagTempRepository
}

func NewTagTempService(tagTemp_repo *repositories.TagTempRepository) *TagTempService {
	return &TagTempService{tagTemp_repo: tagTemp_repo}
}

func (s *TagTempService) CreateTagTemp(tagID string, courseID uuid.UUID) (*model.TagTemp, error) {
	return s.tagTemp_repo.CreateNewTagTemp(tagID, courseID)
}

func (s *TagTempService) CreateTagTempEvent(EventID uuid.UUID, tagIDs []string) ([]model.TagTemp, error) {
	return s.tagTemp_repo.CreateNewTagTempEvent(EventID, tagIDs)
}

// func (s *TagTempService) GetTagTempsByCourseID(courseID uuid.UUID) ([]model.TagTemp, error) {
// 	return s.tagTemp_repo.GetTagTempsByCourseID(courseID)
// }

// func (s *TagTempService) DeleteTagTemp(tagID string, courseID uuid.UUID) error {
// 	return s.tagTemp_repo.DeleteTagTemp(tagID, courseID)
// }

func (s *TagTempService) UpdateTagTemp(courseID string, tags []model.Tag) error {
	return s.tagTemp_repo.UpdateTagTemp(courseID, tags)
}

func (s *TagTempService) GetByEventOrCourse(eventID, courseID string) []string {
	tags := []string{}
	if eventID != "" {
		t, _ := s.tagTemp_repo.FindByEvent(eventID)
		for _, tag := range t {
			tags = append(tags, tag.TagID)
		}
	}
	if courseID != "" {
		t, _ := s.tagTemp_repo.FindByCourse(courseID)
		for _, tag := range t {
			tags = append(tags, tag.TagID)
		}
	}
	return tags
}
