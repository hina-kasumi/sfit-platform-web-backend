package dtos

type AddTeamMemberRequest struct {
	Role string `json:"role" binding:"required"`
}

type AddTeamMemberResponse struct {
	ID        string `json:"id"`
	CreatedAt string `json:"create_at"`
}

type UpdateTeamMemberRequest struct {
	Role string `json:"role" binding:"required"`
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
