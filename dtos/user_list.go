package dtos

type UserListQuery struct {
	Page     int `form:"page" binding:"required"`
	PageSize int `form:"pageSize" binding:"required"`
}

type UserListItem struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type UserListResponse struct {
	Users    []UserListItem `json:"users"`
	Page     int            `json:"page"`
	PageSize int            `json:"pageSize"`
	Total    int64          `json:"total"`
}
