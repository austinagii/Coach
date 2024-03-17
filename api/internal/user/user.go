package user

import (
	"time"
)

type User struct {
	Summary string  `json:"summary"`
	Goals   []*Goal `json:"goals,omitempty"`
}

type Goal struct {
	Id         int          `json:"id,omitempty"`
	Title      string       `json:"summary"`
	Detail     string       `json:"detail"`
	Milestones []*Milestone `json:"milestones,omitempty"`
}

type Milestone struct {
	Id             int       `json:"id"`
	Title          string    `json:"title"`
	Detail         string    `json:"detail"`
	TargetDate     time.Time `json:"target_date,omitempty"`
	IsComplete     bool      `json:"is_complete"`
	CompletionDate time.Time `json:"completion_date,omitempty"`
}
