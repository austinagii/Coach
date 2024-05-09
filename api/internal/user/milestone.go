package user

import "time"

type Milestone struct {
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	TargetDate     time.Time `json:"target_date,omitempty"`
	IsComplete     bool      `json:"is_complete"`
	IsDeleted      bool      `json:"is_deleted,omitempty"`
	CompletionDate time.Time `json:"completion_date,omitempty"`
}
