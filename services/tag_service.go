package services

import (
	"sfit-platform-web-backend/entities"
	"sfit-platform-web-backend/repositories"
)

type TagService struct {
	tagRepo *repositories.TagRepository
}

func NewTagService(tag_repo *repositories.TagRepository) *TagService {
	return &TagService{tagRepo: tag_repo}
}

func (service *TagService) GetAll() ([]entities.Tag, error) {
	return service.tagRepo.FindAll()
}
