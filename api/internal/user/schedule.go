package user

import (
	"time"
)

// DailySchedule represents a collection of timeboxed tasks that a user
// has determined to accomplish.
type DailySchedule struct {
	Tasks []ScheduledTask `json:"tasks" bson:"task"`
}

// ScheduledTask represents a singular timeboxed task that a user wants
// to accomplish.
type ScheduledTask struct {
	Start       time.Time `json:"start" bson:"start"`
	End         time.Time `json:"end" bson:"end"`
	Title       string    `json:"title" bson:"title"`
	Description string    `json:"description" bson:"description"`
}
