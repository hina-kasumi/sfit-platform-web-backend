package dtos

import "time"

type CreateTaskReq struct {
	Name        string    `json:"name" bind:"required"`
	Description string    `json:"description" bind:"required"`
	EventID     string    `json:"event_id"`
	StartDate   time.Time `json:"start_date" bind:"required"`
	Deadline    time.Time `json:"deadline" bind:"required"`
}

type UpdateTaskReq struct {
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	PercentComplete float32   `json:"percent_complete"`
	StartDate       time.Time `json:"start_date"`
	Deadline        time.Time `json:"deadline"`
}

type ListTaskQuery struct {
	Page     int    `form:"page" binding:"required"`
	PageSize int    `form:"pageSize" binding:"required"`
	Name     string `form:"name"`
	EventID  string `form:"event_id"`
}

type AddUserTaskReq struct {
	TaskID string `json:"task_id" binding:"required"`
}
