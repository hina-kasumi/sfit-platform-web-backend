package dtos

type AddTeamMemberRequest struct {
	UserID string `json:"user_id" binding:"required"`
	Role   string `json:"role" binding:"required"`
}

type AddTeamMemberResponse struct {
	ID        string `json:"id"`
	CreatedAt string `json:"create_at"`
}
