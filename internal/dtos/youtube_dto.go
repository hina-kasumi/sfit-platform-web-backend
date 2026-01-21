package dtos

type VideoListResponse struct {
	Kind     string   `json:"kind"`
	Etag     string   `json:"etag"`
	Items    []Video  `json:"items"`
	PageInfo PageInfo `json:"pageInfo"`
}

type Video struct {
	Kind           string         `json:"kind"`
	Etag           string         `json:"etag"`
	ID             string         `json:"id"`
	ContentDetails ContentDetails `json:"contentDetails"`
}

type ContentDetails struct {
	Duration        string        `json:"duration"`
	Dimension       string        `json:"dimension"`
	Definition      string        `json:"definition"`
	Caption         string        `json:"caption"`
	LicensedContent bool          `json:"licensedContent"`
	ContentRating   ContentRating `json:"contentRating"`
	Projection      string        `json:"projection"`
}

type ContentRating struct {
	// Add fields if needed
}

type PageInfo struct {
	TotalResults   int `json:"totalResults"`
	ResultsPerPage int `json:"resultsPerPage"`
}
