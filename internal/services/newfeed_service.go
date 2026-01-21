package services

import (
	"sfit-platform-web-backend/internal/dtos"
	"sfit-platform-web-backend/internal/model"

	"github.com/google/uuid"
)

type NewFeedService struct {
	courseSer *CourseService
	userSer   *UserService
	taskSer   *TaskService
	eventSer  *EventService
}

func NewNewFeedService(courseSer *CourseService, userSer *UserService, taskSer *TaskService, evenetSer *EventService) *NewFeedService {
	return &NewFeedService{
		courseSer: courseSer,
		userSer:   userSer,
		taskSer:   taskSer,
		eventSer:  evenetSer,
	}
}

func (n *NewFeedService) GetNewFeed(userID uuid.UUID) (*dtos.NewFeedResponse, error) {
	events, totalEvents, err := n.eventSer.GetEvents(1, 3, "", "", string(model.StatusUpcoming), "", userID.String())
	if err != nil {
		return nil, err
	}

	totalLearningCourses, err := n.courseSer.GetTotalCourseOfUser(userID, model.UserCourseStatusLearn)
	if err != nil {
		return nil, err
	}

	isCompleted := false
	tasks, totalPeddingTasks, err := n.taskSer.ListTasksByUserID(userID.String(), 1, 3, &isCompleted)
	if err != nil {
		return nil, err
	}

	return &dtos.NewFeedResponse{
		TotalLearningCourses: totalLearningCourses,
		TotalEvents:          totalEvents,
		TotalPeddingTasks:    totalPeddingTasks,
		Events:               events,
		Tasks:                tasks,
	}, nil
}
