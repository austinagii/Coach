package user

import (
	"fmt"
	"time"
)

type User struct {
	Id       string         `json:"id,omitempty" bson:"_id,omitempty"`
	Name     string         `json:"name" bson:"name"`
	Summary  string         `json:"summary,omitempty" bson:"summary,omitempty"`
	Goals    []*Goal        `json:"goals,omitempty" bson:"goals,omitempty"`
	Schedule *DailySchedule `json:"schedule,omitempty" bson:"schedule,omitempty"`
}

type Goal struct {
	Id             int          `json:"id,omitempty" bson:"_id,omitempty"`
	Title          string       `json:"title" bson:"title"`
	Description    string       `json:"description" bson:"description"`
	Milestones     []*Milestone `json:"milestones,omitempty" bson:"milestones,omitempty"`
	IsComplete     bool         `json:"is_complete" bson:"is_complete"`
	IsDeleted      bool         `json:"is_deleted" bson:"is_deleted"`
	CompletionDate time.Time    `json:"completion_date,omitempty" bson:"completion_date,omitempty"`
}

type Milestone struct {
	Title          string    `json:"title" bson:"title"`
	Description    string    `json:"description" bson:"description"`
	TargetDate     time.Time `json:"target_date,omitempty" bson:"target_date,omitempty"`
	IsComplete     bool      `json:"is_complete" bson:"is_complete"`
	CompletionDate time.Time `json:"completion_date,omitempty" bson:"completion_date,omitempty"`
	IsDeleted      bool      `json:"is_deleted" bson:"is_deleted"`
}

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

func NewUser(name string) *User {
	return &User{Name: name}
}

func (user *User) AddNewGoal(goal *Goal) {
	goalId := 0
	if len(user.Goals) > 0 {
		goalId = user.Goals[len(user.Goals)-1].Id + 1
	}
	goal.Id = goalId
	user.Goals = append(user.Goals, goal)
}

func (user *User) GetGoalById(id int) (*Goal, error) {
	var goal *Goal
	for id, g := range user.Goals {
		if g.Id == id {
			goal = g
			break
		}
	}
	if goal == nil {
		return nil, fmt.Errorf("No goal with id '%d' could be found", id)
	}
	return goal, nil
}
