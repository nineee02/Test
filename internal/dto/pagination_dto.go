package dto

type PaginationRequest struct {
	Page    int    `json:"page" validate:"required,min=1"`
	Limit   int    `json:"limit" validate:"required,min=1"`
	Keyword string `json:"keyword", omitempty`
}
