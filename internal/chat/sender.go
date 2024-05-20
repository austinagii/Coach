package chat

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type Sender string

const (
	SenderUser      Sender = "user"
	SenderAssistant Sender = "assistant"
	SenderSystem    Sender = "system"
)

var ErrUnknownSender = errors.New("Unknown sender")

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

	switch strings.ToLower(s) { // Perform a case insenstitive match.
	case "user":
		*sender = SenderUser
	case "assistant":
		*sender = SenderAssistant
	case "system":
		*sender = SenderSystem
	default:
		return fmt.Errorf("%w '%s'", ErrUnknownSender, s)
	}
	return nil
}
