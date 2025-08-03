package dtos

type AddTeamMemberRequest struct {
	UserID string `json:"user_id" binding:"required"`
	Role   string `json:"role" binding:"required"`
}

type AddTeamMemberResponse struct {
	ID        string `json:"id"`
	CreatedAt string `json:"create_at"`
}

type DeleteTeamMemberRequest struct {
	UserID string `json:"user_id" binding:"required"`
}

type UpdateTeamMemberRequest struct {
	UserID string `json:"user_id" binding:"required"`
	Role   string `json:"role" binding:"required"`
}

type UserJoinedTeamResponse struct {
	TeamID string `json:"team_id"`
	Name   string `json:"name"`
	Role   string `json:"role"`
}
