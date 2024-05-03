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
