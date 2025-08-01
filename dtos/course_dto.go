package dtos

type CreateCourseRequest struct {
    Title       string   `json:"title" binding:"required"`
    Description string   `json:"description" binding:"required"`
    Type        string   `json:"type" binding:"required"`
    Target      []string `json:"targets"`
    Require     []string `json:"requires"`
    Teachers    []string `json:"teachers"`
    Language    string   `json:"language" binding:"required"`
    Certificate bool     `json:"certificate"`
    Level       string   `json:"level" binding:"required"`
    Tags        []string `json:"tags"`
}

type CreateCourseResponse struct {
    ID        string    `json:"id"`
    CreatedAt string    `json:"createdAt"`
}