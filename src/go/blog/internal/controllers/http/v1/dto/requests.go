package dto

type FilterLogs struct {
	Page      int    `json:"page"`
	Limit     int    `json:"limit"`
	Role      string `json:"role"`
	ContextID string `json:"context_id"`
	Search    string `json:"search"`
}

type BusinessLog struct {
	EventType   string `json:"event_type,omitempty"` // CREATE, UPDATE, DELETE
	Entity      string `json:"entity,omitempty"`     // user, context
	Username    string `json:"user_name,omitempty"`
	UserRole    string `json:"user_role,omitempty"` // SA, CA
	Context     string `json:"context,omitempty"`
	EntityID    string `json:"entity_id,omitempty"` // userID or ContextID
	OldValue    any    `json:"old_value,omitempty"`
	NewValue    any    `json:"new_value,omitempty"`
	Description string `json:"description,omitempty"`
}
