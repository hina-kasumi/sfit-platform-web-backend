package services

import (
	"errors"
	"sfit-platform-web-backend/entities"
	"sfit-platform-web-backend/repositories"
	"time"
)

type TaskService struct {
	taskRepo *repositories.TaskRepository
}

func NewTaskService(repo *repositories.TaskRepository) *TaskService {
	return &TaskService{
		taskRepo: repo,
	}
}

func (ts *TaskService) GetTasks(page, pageSize int, name, eventID string) ([]*entities.Task, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	return ts.taskRepo.GetTasks(page, pageSize, name, eventID)
}

func (ts *TaskService) CreateTask(createdBy, name, description, eventID string, startDate, deadline time.Time) (*entities.Task, error) {
	return ts.taskRepo.CreateTask(createdBy, name, description, eventID, startDate, deadline)
}

func (ts *TaskService) GetTaskByID(id string) (*entities.Task, error) {
	return ts.taskRepo.GetTaskByID(id)
}

func (ts *TaskService) UpdateTask(id, name, description string, percent_complete float32, startTime, deadline time.Time) error {
	if id == "" {
		return errors.New("invalid task ID")
	}
	task, err := ts.GetTaskByID(id)
	if err != nil || task == nil {
		return err
	}

	return ts.taskRepo.UpdateTask(id, name, description, percent_complete, startTime, deadline)
}

func (ts *TaskService) DeleteTask(id string) error {
	if id == "" {
		return errors.New("invalid task ID")
	}
	return ts.taskRepo.DeleteTask(id)
}
