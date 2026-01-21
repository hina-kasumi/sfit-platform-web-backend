package dtos

type UpdateUserDto struct {
	Email       string `json:"email"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type UserListItem struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}
