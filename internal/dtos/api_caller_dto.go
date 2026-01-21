package dtos

type ApiCallerRp struct {
	StatusCode int    `json:"code"`
	Body       []byte `json:"body"`
}
