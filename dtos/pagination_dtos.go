package dtos

type PageListQuery struct {
	Page     int `form:"page" json:"page" binding:"required,min=1"`
	PageSize int `form:"page_size" json:"page_size" binding:"required,min=1,max=100"`
}

type PageListResp struct {
	TotalCount int64 `json:"total_count"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Items      any   `json:"items"`
}
