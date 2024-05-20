package chat

import (
	"encoding/json"
	"testing"
)

func TestUnmarshalledMessageAreNotNewByDefault(t *testing.T) {
	jsonString := []byte(`{"sender": "user", "text": "test"}`)

	var message Message
	if err := json.Unmarshal(jsonString, &message); err != nil {
		t.Fatalf("Failed to unmarshal JSON message '%s': %v", jsonString, err)
	}

	if message.IsNew != false {
		t.Fatalf("Incorrect value for 'IsNew': expected %v, got %v", message.IsNew, false)
	}
}
