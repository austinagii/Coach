package user

type User struct {
	Id      string  `json:"id,omitempty" bson:"_id"`
	Name    string  `json:"name"`
	Gender  string  `json:"gender"`
	Summary string  `json:"summary,omitempty"`
	Goals   []*Goal `json:"goals,omitempty"`
}
