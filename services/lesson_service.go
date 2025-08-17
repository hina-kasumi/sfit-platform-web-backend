package services

import (
	"encoding/json"
	"errors"
	"os"
	"sfit-platform-web-backend/dtos"
	"sfit-platform-web-backend/entities"
	"sfit-platform-web-backend/repositories"
	"sfit-platform-web-backend/utils/caller"
	"sfit-platform-web-backend/utils/converter"
	"strings"

	"github.com/google/uuid"
)

type LessonService struct {
	courseSer     *CourseService
	lessonRepo    *repositories.LessonRepository
	youtubeApiUrl string
	youtubeApiKey string
	youtubePart   string
}

func NewLessonService(lessonRepo *repositories.LessonRepository, courseSer *CourseService) *LessonService {
	return &LessonService{
		lessonRepo:    lessonRepo,
		courseSer:     courseSer,
		youtubeApiUrl: os.Getenv("YOUTUBE_API_URL"),
		youtubeApiKey: os.Getenv("YOUTUBE_API_KEY"),
		youtubePart:   os.Getenv("YOUTUBE_PART"),
	}
}
func (s *LessonService) createLesson(moduleID string, req dtos.LessonRequest) (*entities.Lesson, error) {
	var lesson *entities.Lesson
	var err error
	parsedModuleID, err := uuid.Parse(moduleID)
	if err != nil {
		return nil, err
	}
	switch req.Type {
	case string(entities.QuizLesson):
		lesson = entities.NewQuizLesson(
			parsedModuleID,
			req.Title,
			req.Description,
			req.Duration,
			req.QuizContent,
		)
	case string(entities.OnlineLesson):
		// Check if the video URL is from YouTube
		if strings.Contains(req.VideoURL, "youtube.com") {
			videoID := strings.Split(req.VideoURL, "v=")[1]
			callerRp, _ := caller.GetRequest(s.youtubeApiUrl, map[string]string{
				"part": s.youtubePart,
				"id":   videoID,
				"key":  s.youtubeApiKey,
			})
			var videoListResp dtos.VideoListResponse
			_ = json.Unmarshal(callerRp.Body, &videoListResp)
			timeDuration, _ := converter.ISO8601ToNumber(videoListResp.Items[0].ContentDetails.Duration)
			req.Duration = int(timeDuration.Seconds())
		}
		if req.VideoURL == "" {
			return nil, errors.New("video URL is required for online lessons")
		}
		lesson = entities.NewOnlineLesson(
			parsedModuleID,
			req.Title,
			req.Description,
			req.Duration,
			req.VideoURL,
		)
	case string(entities.OfflineLesson):
		if req.Location == "" || req.Date.IsZero() {
			return nil, errors.New("location and date are required for offline lessons")
		}
		if req.Duration <= 0 {
			req.Duration = 7200
		}
		lesson = entities.NewOfflineLesson(
			parsedModuleID,
			req.Title,
			req.Description,
			req.Duration,
			req.Location,
			req.Date,
		)
	case string(entities.ReadingLesson):
		lesson = entities.NewReadingLesson(
			parsedModuleID,
			req.Title,
			req.Description,
			req.Duration,
			req.ReadingContent,
		)
	default:
		return nil, errors.New("invalid lesson type")
	}
	if lesson.Duration <= 0 {
		return nil, errors.New("invalid duration for online lessons")
	}
	return lesson, nil
}

func (s *LessonService) CreateNewLesson(moduleID string, req dtos.LessonRequest) (*entities.Lesson, error) {
	lesson, err := s.createLesson(moduleID, req)
	if err != nil {
		return nil, err
	}
	err = s.courseSer.UpdateTotalTime(moduleID, lesson.Duration)
	if err != nil {
		return nil, err
	}
	return s.lessonRepo.CreateLesson(lesson)
}

func (s *LessonService) UpdateLesson(moduleID string, lessonID string, req dtos.LessonRequest) error {
	lesson, err := s.createLesson(moduleID, req)
	if err != nil {
		return err
	}
	parsedLessonID, err := uuid.Parse(lessonID)
	if err != nil {
		return err
	}
	lesson.ID = parsedLessonID
	return s.lessonRepo.UpdateLesson(lesson)
}

func (s *LessonService) GetLessonByID(id string) (*entities.Lesson, error) {
	lesson, err := s.lessonRepo.GetLessonByID(id)
	if err != nil {
		return nil, err
	}
	return lesson, nil
}

func (s *LessonService) DeleteLessonByID(id string) error {
	lesson, err := s.lessonRepo.GetLessonByID(id)
	if err != nil {
		return err
	}
	err = s.courseSer.UpdateTotalTime(lesson.ModuleID.String(), -lesson.Duration)
	if err != nil {
		return err
	}
	return s.lessonRepo.DeleteLessonByID(id)
}
