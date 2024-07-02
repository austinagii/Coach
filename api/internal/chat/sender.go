package chat

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// Sender represents the origin of a message.
type Sender string

// Possible values for Sender.
const (
	SenderUser      Sender = "user"
	SenderAssistant Sender = "assistant"
)

// ErrUnknownSender is returned when an unrecognized sender is unmarshaled.
var ErrUnknownSender = errors.New("Unknown sender")

// String returns the string representation of the Sender.
func (sender Sender) String() string {
	return string(sender)
}

// MarshalJSON converts the Sender to its JSON representation.
func (sender Sender) MarshalJSON() ([]byte, error) {
	return json.Marshal(sender.String())
}

// UnmarshalJSON parses JSON data into a Sender, performing a case-insensitive match.
func (sender *Sender) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	switch strings.ToLower(s) {
	case "user":
		*sender = SenderUser
	case "assistant":
		*sender = SenderAssistant
	default:
		return fmt.Errorf("%w '%s'", ErrUnknownSender, s)
	}
	return nil
}
