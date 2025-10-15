package services

import (
	"sfit-platform-web-backend/dtos"
	"sfit-platform-web-backend/entities"

	"github.com/google/uuid"
)

type NewFeedService struct {
	courseSer *CourseService
	userSer   *UserService
	taskSer   *TaskService
	evenetSer *EventService
}

func NewNewFeedService(courseSer *CourseService, userSer *UserService, taskSer *TaskService, evenetSer *EventService) *NewFeedService {
	return &NewFeedService{
		courseSer: courseSer,
		userSer:   userSer,
		taskSer:   taskSer,
		evenetSer: evenetSer,
	}
}

func (n *NewFeedService) GetNewFeed(userID uuid.UUID) (*dtos.NewFeedResponse, error) {
	events, totalEvents, err := n.evenetSer.GetEvents(1, 3, "", "", string(entities.StatusUpcoming), "", userID.String())
	if err != nil {
		return nil, err
	}

	totalLearningCourses, err := n.courseSer.GetTotalCourseOfUser(userID, entities.UserCourseStatusLearn)
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
