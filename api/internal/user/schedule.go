package user

// import (
// 	"time"
// )

// DailySchedule represents a collection of timeboxed activities that a user
// has determined to accomplish.
type DailySchedule struct {
	Activities []ScheduledActivity `json:"activities" bson:"activities"`
}

// ScheduledActivity represents a singular timeboxed activity that a user wants
// to complete.
type ScheduledActivity struct {
	Start       string `json:"start" bson:"start"`
	End         string `json:"end" bson:"end"`
	Title       string `json:"title" bson:"title"`
	Description string `json:"description" bson:"description"`
}
