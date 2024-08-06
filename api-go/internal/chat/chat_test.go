package chat

import (
	"testing"
)

func TestNewChatAcceptsMessageLimit(t *testing.T) {
	var expectedMessageLimit int = 30
	chat := NewChat(expectedMessageLimit)

	if chat.messageLimit != expectedMessageLimit {
		t.Fatalf("Invalid message limit; Expected %d got %d", expectedMessageLimit, chat.messageLimit)
	}
}

func TestChatUsesDefaultMessageLimitIfNoneSpecified(t *testing.T) {
	c := NewChat()

	if c.messageLimit != DefaultMessageLimit {
		t.Fatalf("Incorred message limit; expected %d, got %d", DefaultMessageLimit, c.messageLimit)
	}
}
