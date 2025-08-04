package services

import (
	"sfit-platform-web-backend/entities"
	"sfit-platform-web-backend/repositories"

	"gorm.io/gorm"
)

type TagService struct {
	tag_repo *repositories.TagRepository
}

func NewTagService(tag_repo *repositories.TagRepository) *TagService {
	return &TagService{tag_repo: tag_repo}
}

func (service *TagService) GetAll() ([]entities.Tag, error) {
	return service.tag_repo.FindAll()
}

func (s *TagService) EnsureTags(tags []string) ([]entities.Tag, error) {
    var result []entities.Tag
    for _, tagName := range tags {
        tag, err := s.tag_repo.FindByID(tagName)
        if err == gorm.ErrRecordNotFound {
            tag = &entities.Tag{ID: tagName}
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