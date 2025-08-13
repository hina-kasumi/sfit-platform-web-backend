package services

import (
	"sfit-platform-web-backend/dtos"
	"sfit-platform-web-backend/repositories"
)

type TaskService struct {
	taskRepo *repositories.TaskRepository
}

func NewTaskService(taskRepo *repositories.TaskRepository) *TaskService {
	return &TaskService{
		taskRepo: taskRepo,
	}
}

func (s *TaskService) GetListTasks(pageNum, pageSize int) ([]dtos.TaskGeneralInformationResponse, dtos.TaskPaginationResponse, error) {
	tasks, totalTasks, err := s.taskRepo.GetListTasks(pageNum, pageSize)
	if err != nil {
		return nil, dtos.TaskPaginationResponse{}, err
	}

	taskResponses := make([]dtos.TaskGeneralInformationResponse, len(tasks))
	for i, task := range tasks {
		taskResponses[i] = dtos.TaskGeneralInformationResponse{
			ID:              task.ID,
			EventID:         task.EventID,
			Name:            task.Name,
			Description:     task.Description,
			StartDate:       task.StartDate,
			Deadline:        task.Deadline,
			Assignee:        task.Assignee,
			PercentComplete: task.PercentComplete,
			CreatedAt:       task.CreatedAt,
			UpdatedAt:       task.UpdatedAt,
		}
	}

	pagination := dtos.TaskPaginationResponse{
		CurrentPage: pageNum,
		TotalPages:  (totalTasks + pageSize - 1) / pageSize,
		TotalTasks:  totalTasks,
	}

	return taskResponses, pagination, nil
}
