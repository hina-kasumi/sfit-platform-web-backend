package dtos

import "github.com/google/uuid"

type CreateTeamRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type CreateTeamResponse struct {
	ID       uuid.UUID `json:"id"`
	CreateAt string    `json:"create_at"`
}

type UpdateTeamRequest struct {
	ID          uuid.UUID `json:"id" binding:"required"`
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description" binding:"required"`
}
type UpdateTeamResponse struct {
	UpdatedAt string `json:"updateAt"`
}
