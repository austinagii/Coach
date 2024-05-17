package chat

import (
	"encoding/json"
	"fmt"
)

type Sender string

const (
	SenderUser      Sender = "User"
	SenderAssistant Sender = "Assistant"
	SenderSystem    Sender = "System"
)

func (sender Sender) String() string {
	return string(sender)
}

func (sender Sender) MarshalJSON() ([]byte, error) {
	return json.Marshal(sender.String())
}

func (sender *Sender) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	switch s {
	case "User":
		*sender = SenderUser
	case "Assistant":
		*sender = SenderAssistant
	case "System":
		*sender = SenderSystem
	default:
		*sender = ""
		return fmt.Errorf("No sender for string '%s'", s)
	}
	return nil
}
