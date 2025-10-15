package dtos

import (
	"sfit-platform-web-backend/entities"
)

type NewFeedResponse struct {
	TotalLearningCourses int64                  `json:"total_learning_courses"`
	TotalEvents          int64                  `json:"total_events"`
	TotalPeddingTasks    int64                  `json:"total_pedding_tasks"`
	Events               []entities.Event       `json:"events"`
	Tasks                []*ResponseTasksOfUser `json:"tasks"`
}
