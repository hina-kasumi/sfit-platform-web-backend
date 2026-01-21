package services

import (
	"sfit-platform-web-backend/internal/model"
	"sfit-platform-web-backend/internal/repositories"

	"gorm.io/gorm"
)

type TagService struct {
	tag_repo *repositories.TagRepository
}

func NewTagService(tag_repo *repositories.TagRepository) *TagService {
	return &TagService{tag_repo: tag_repo}
}

func (service *TagService) GetAll() ([]model.Tag, error) {
	return service.tag_repo.FindAll()
}

func (s *TagService) EnsureTags(tags []string) ([]model.Tag, error) {
	var result []model.Tag
	for _, tagName := range tags {
		tag, err := s.tag_repo.FindByID(tagName)
		if err == gorm.ErrRecordNotFound {
			tag = &model.Tag{ID: tagName}
			if err := s.tag_repo.CreateNewTag(tag); err != nil {
				return nil, err
			}
		} else if err != nil {
			return nil, err
		}
		result = append(result, *tag)
	}
	return result, nil
}
