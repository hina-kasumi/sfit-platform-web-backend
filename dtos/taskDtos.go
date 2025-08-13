package dtos

import (
	"time"
)

//
// REQUEST DTOs
//

//
// RESPONSE DTOs
//

type TaskListResponse struct {
	Tasks      []TaskGeneralInformationResponse `json:"tasks"`
	Pagination TaskPaginationResponse           `json:"pagination"`
}

type TaskGeneralInformationResponse struct {
	ID              string    `json:"id"`
	EventID         string    `json:"event_id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	StartDate       time.Time `json:"start_date"`
	Deadline        time.Time `json:"deadline"`
	Assignee        string    `json:"assignee"`
	PercentComplete int       `json:"percent_complete"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type TaskPaginationResponse struct {
	CurrentPage int `json:"currentPage"`
	TotalPages  int `json:"totalPages"`
	TotalTasks  int `json:"total"`
}
