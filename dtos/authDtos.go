package dtos

type LoginRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Username  string `json:"username" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,password"`
	FullName  string `json:"full_name" binding:"required"`
	Phone     string `json:"phone" binding:"required"`
	ClassName string `json:"class_name" binding:"required"`
	Khoa      string `json:"khoa" binding:"required"`
	MSV       string `json:"msv" binding:"required"`
}
