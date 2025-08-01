package dtos

type UpdateUserProfileRequest struct {
	FullName     string            `json:"full_name"`
	ClassName    string            `json:"class_name"`
	Khoa         string            `json:"khoa"`
	Phone        string            `json:"phone"`
	Introduction string            `json:"introduction"`
	SocialLink   map[string]string `json:"social_link"`
}
