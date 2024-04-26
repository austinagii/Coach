package user

import "time"

type Goal struct {
	Id             int          `json:"id,omitempty"`
	Title          string       `json:"title"`
	Description    string       `json:"description"`
	Milestones     []*Milestone `json:"milestones,omitempty"`
	IsComplete     bool         `json:"is_complete"`
	IsDeleted      bool         `json:"is_deleted"`
	CompletionDate time.Time    `json:"completion_date"`
}
