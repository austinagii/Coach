package user

import "time"

// Goal represents a goal that can be set and tracked.
type Goal struct {
	Id             int          `json:"id,omitempty" bson:"_id,omitempty"`
	Title          string       `json:"title" bson:"title"`
	Description    string       `json:"description" bson:"description"`
	Milestones     []*Milestone `json:"milestones,omitempty" bson:"milestones,omitempty"`
	IsComplete     bool         `json:"is_complete" bson:"is_complete"`
	IsDeleted      bool         `json:"is_deleted" bson:"is_deleted"`
	CompletionDate time.Time    `json:"completion_date,omitempty" bson:"completion_date,omitempty"`
}
