package dto

type ValidateQuery struct {
	Page      int    `json:"page"`
	Limit     int    `json:"limit"`
	Role      string `json:"role"`
	ContextID string `json:"context_id"`
	Search    string `json:"search"`
}
