package dto

type ValidateQuery struct {
	Page      int    `json:"page"`
	Limit     int    `json:"limit"`
	Role      string `json:"role"`
	ContextID string `json:"context_id"`
	Search    string `json:"search"`
}

type DtoBusinessLog struct {
	//Timestamp   time.Time `json:"timestamp,omitzero" gorm:"column:timestamp;type:timestamp with time zone;not null"`
	EventType   string `json:"event_type,omitempty"` // CREATE, UPDATE, DELETE
	Entity      string `json:"entity,omitempty"`     // user, context
	Username    string `json:"user_name,omitempty"`
	UserRole    string `json:"user_role,omitempty"` // SA, CA
	ContextID   string `json:"context,omitempty"`
	EntityID    string `json:"entity_id,omitempty"` // userID or ContextID
	OldValue    Value  `json:"old_value,omitempty"`
	NewValue    Value  `json:"new_value,omitempty"`
	Description string `json:"description,omitempty"`
	ContextStr  string `json:"context_str,omitempty"`
}

type Value struct {
	Email      string `json:"email"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Login      string `json:"login"`
	Patronymic string `json:"patronymic"`
	Role       string `json:"role"`
	UserID     string `json:"user_id"`
}
