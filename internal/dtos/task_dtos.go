package dtos

import (
	"time"
)

type CreateTaskReq struct {
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description" binding:"required"`
	EventID     string    `json:"event_id"`
	StartDate   time.Time `json:"start_date" binding:"required"`
	Deadline    time.Time `json:"deadline" binding:"required"`
}

type UpdateTaskReq struct {
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	PercentComplete float32   `json:"percent_complete"`
	StartDate       time.Time `json:"start_date"`
	Deadline        time.Time `json:"deadline"`
}

type ListTaskQuery struct {
	PageListQuery
	Name        string `form:"name"`
	EventID     string `form:"event_id"`
	IsCompleted *bool  `form:"is_completed"`
}

type AddUserTaskReq struct {
	TaskID string `json:"task_id" binding:"required"`
}

type ListTaskOfUserReq struct {
	PageListQuery
	IsCompleted *bool `form:"is_completed"`
}

type ListTasksByEventID struct {
	PageListQuery
	IsCompleted *bool `form:"is_completed"`
}

type UpdateTaskUserStatusReq struct {
	IsCompleted bool `json:"is_completed"`
}

type ResponseTasksOfUser struct {
	ID              string    `json:"id"`
	EventID         *string   `json:"event_id,omitempty"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	StartDate       time.Time `json:"start_date"`
	Deadline        time.Time `json:"deadline"`
	PercentComplete float32   `json:"percent_complete"`
	CreateBy        string    `json:"create_by"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	IsCompleted     bool      `json:"is_completed"`
}
