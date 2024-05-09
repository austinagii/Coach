package user

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Gender string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
)

type User struct {
	Id      string  `json:"id,omitempty" bson:"_id,omitempty"`
	Name    string  `json:"name" bson:"name"`
	Gender  Gender  `json:"gender" bson:"gender"`
	Summary string  `json:"summary,omitempty" bson:"summary,omitempty"`
	Goals   []*Goal `json:"goals,omitempty" bson:"goals,omitempty"`
}

func NewUser(name string, gender Gender) *User {
	return &User{Name: name, Gender: gender}
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

func (g *Gender) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	s = strings.ToLower(s)
	switch s {
	case "male":
		*g = GenderMale
	case "female":
		*g = GenderFemale
	default:
		return fmt.Errorf("No gender found for string '%s'", s)
	}
	return nil
}
