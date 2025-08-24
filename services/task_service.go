package services

import (
	"errors"
	"sfit-platform-web-backend/dtos"
	"sfit-platform-web-backend/entities"
	"sfit-platform-web-backend/repositories"
	"time"

	"github.com/google/uuid"
)

type TaskService struct {
	taskRepo *repositories.TaskRepository
}

func NewTaskService(repo *repositories.TaskRepository) *TaskService {
	return &TaskService{
		taskRepo: repo,
	}
}

func (ts *TaskService) GetTasks(page, pageSize int, name, eventID string, isCompleted *bool) ([]*entities.Task, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	return ts.taskRepo.GetTasks(page, pageSize, name, eventID, isCompleted)
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

func (ts *TaskService) ListTasksByUserID(userID string, page, pageSize int, isCompleted *bool) ([]*dtos.ResponseTasksOfUser, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	userIDParsed, err := uuid.Parse(userID)
	if err != nil {
		return nil, 0, err
	}
	return ts.taskRepo.ListTasksByUserID(userIDParsed, page, pageSize, isCompleted)
}

func (ts *TaskService) AddUserTask(userID, taskID string) (*entities.TaskAssignments, error) {
	userIDParsed, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}
	taskIDParsed, err := uuid.Parse(taskID)
	if err != nil {
		return nil, err
	}
	return ts.taskRepo.AddUserTask(taskIDParsed, userIDParsed)
}

func (ts *TaskService) UpdateTaskUserStatus(taskID, userID string, isCompleted bool) error {
	userIDParsed, err := uuid.Parse(userID)
	if err != nil {
		return err
	}
	taskIDParsed, err := uuid.Parse(taskID)
	if err != nil {
		return err
	}

	return ts.taskRepo.UpdateTaskUserStatus(taskIDParsed, userIDParsed, isCompleted)
}

func (ts *TaskService) DeleteUserTask(userID, taskID string) error {
	userIDParsed, err := uuid.Parse(userID)
	if err != nil {
		return err
	}
	taskIDParsed, err := uuid.Parse(taskID)
	if err != nil {
		return err
	}
	return ts.taskRepo.DeleteUserTask(taskIDParsed, userIDParsed)
}
