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

type TeamMemberUserResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type TeamMembersResponse struct {
	Users    []TeamMemberUserResponse `json:"users"`
	Page     int                      `json:"page"`
	PageSize int                      `json:"page_size"`
	Total    int64                    `json:"total"`
}
