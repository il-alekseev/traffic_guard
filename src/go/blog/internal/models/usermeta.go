package models

type UserMeta struct {
	UUID       string `json:"uuid"`
	Username   string `json:"username"`
	ClientRole string `json:"client_role"`
	ShortRole  string `json:"short_role"`
	ContextID  string `json:"context_id"`
}
