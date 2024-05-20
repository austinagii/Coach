package chat

import (
	"errors"
	"testing"
)

func TestUnmarshalJSON(t *testing.T) {
	testCases := map[string]Sender{
		`"user"`:      SenderUser,
		`"assistant"`: SenderAssistant,
		`"system"`:    SenderSystem,
		`"USER"`:      SenderUser,
	}

	var actualSender Sender
	for jsonString, expectedSender := range testCases {
		err := actualSender.UnmarshalJSON([]byte(jsonString))
		if err != nil {
			t.Fatalf("Failed to unmarshall JSON '%s': %v", jsonString, err)
		}

		if actualSender != expectedSender {
			t.Fatalf("Expected sender '%s', got '%s'", expectedSender, actualSender)
		}
	}
}

func TestUnmarshalJSONReturnsErrorForUnknownSender(t *testing.T) {
	var sender Sender

	expectedError := ErrUnknownSender
	actualError := sender.UnmarshalJSON([]byte(`"unknown"`))
	if actualError == nil || !errors.Is(actualError, expectedError) {
		t.Fatalf("Expected error '%v', got '%v'", expectedError, actualError)
	}
}

func TestUnmarshalJSONReturnsErrorForEmptyString(t *testing.T) {
	var sender Sender

	expectedError := ErrUnknownSender
	actualError := sender.UnmarshalJSON([]byte(`""`))
	if actualError == nil || !errors.Is(actualError, expectedError) {
		t.Fatalf("Expected error '%v', got '%v'", expectedError, actualError)
	}
}
