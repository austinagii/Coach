package chat

import (
	"testing"
	"time"
)

func TestNewMessage(t *testing.T) {
	expectedContent := "test"
	expectedSender := SenderUser
	msg := NewMessage(expectedSender, expectedContent)
	now := time.Now().UnixMilli()

	if msg.Sender != expectedSender {
		t.Errorf("Invalid message sender; expected '%s', got '%s'", expectedSender, msg.Sender)
	}
	if msg.Content != expectedContent {
		t.Errorf("Invalid message content; expected '%s', got '%s'", expectedContent, msg.Sender)
	}
	if msg.CreatedAt > now || msg.CreatedAt < now-1000 {
		t.Errorf("CreatedAt not within reasonable range, got %v", msg.CreatedAt)
	}
}

func TestNewMessageIgnoresAdditionalContent(t *testing.T) {
	expectedContent := "test"
	additionalContent := "ignore"

	msg := NewMessage(SenderUser, expectedContent, additionalContent)
	if msg.Content != expectedContent {
		t.Errorf("Incorrect message content; expected: '%s', got '%s'", expectedContent, msg.Content)
	}
}

func TestNewUserMessage(t *testing.T) {
	msg := NewUserMessage("Hello")
	if msg.Sender != SenderUser {
		t.Errorf("Expected sender to be SenderUser, got %v", msg.Sender)
	}
	if msg.Content != "Hello" {
		t.Errorf("Expected content to be 'Hello', got %v", msg.Content)
	}
	// You can add more checks for CreatedAt if needed
}

func TestNewAssistantMessage(t *testing.T) {
	msg := NewAssistantMessage("Hi")
	if msg.Sender != SenderAssistant {
		t.Errorf("Expected sender to be SenderAssistant, got %v", msg.Sender)
	}
	if msg.Content != "Hi" {
		t.Errorf("Expected content to be 'Hi', got %v", msg.Content)
	}
	// You can add more checks for CreatedAt if needed
}

func TestMessageString(t *testing.T) {
	m := NewUserMessage("test")
	expected := "user: test"
	actual := m.String()
	if actual != expected {
		t.Fatalf("Invalid value; expected: '%s', got '%s'", expected, actual)
	}
}
