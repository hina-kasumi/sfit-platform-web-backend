package services

import (
    "sfit-platform-web-backend/entities"
    "sfit-platform-web-backend/repositories"
    "time"

    "github.com/google/uuid"
)

type CourseService struct {
    course_repo *repositories.CourseRepository
}

func NewCourseService(course_repo *repositories.CourseRepository) *CourseService {
    return &CourseService{course_repo: course_repo}
}

func (s *CourseService) CreateCourse(
    title, description, courseType string,
    targets, requires, teachers []string,
    language string,
    certificate bool,
    level string,
    // tags []entities.Tag,
) (uuid.UUID, time.Time, error) {
    course := entities.Course{
        ID:          uuid.New(),
        Title:       title,
        Description: description,
        Type:        courseType,
        Target:      targets,
        Require:     requires,
        Teachers:    teachers,
        Language:    language,
        Certificate: certificate,
        Level:       level,
        CreatedAt:   time.Now(),
    }
    if err := s.course_repo.CreateNewCourse(&course); err != nil {
        return uuid.Nil, time.Time{}, err
    }
    return course.ID, course.CreatedAt, nil
}