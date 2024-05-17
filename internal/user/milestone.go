package user

import "time"

// Milestone represents an identified step on the way to achieving a goal.
type Milestone struct {
	Title          string    `json:"title" bson:"title"`
	Description    string    `json:"description" bson:"description"`
	TargetDate     time.Time `json:"target_date,omitempty" bson:"target_date,omitempty"`
	IsComplete     bool      `json:"is_complete" bson:"is_complete"`
	CompletionDate time.Time `json:"completion_date,omitempty" bson:"completion_date,omitempty"`
	IsDeleted      bool      `json:"is_deleted" bson:"is_deleted"`
}
