package chat

import (
	"fmt"
	"time"
)

// Message represents a one-way communication between a user and an assistant.
type Message struct {
	Sender    Sender `json:"sender" bson:"sender"`
	Content   string `json:"content" bson:"content"`
	CreatedAt int64  `json:"-" bson:"created_at"`
}

// NewUserMessage creates a new message with optional content.
func NewMessage(sender Sender, content ...string) *Message {
	var messageContent string
	if len(content) > 0 {
		messageContent = content[0]
	}

	return &Message{
		Sender:    sender,
		Content:   messageContent,
		CreatedAt: time.Now().UnixMilli(),
	}
}

// NewUserMessage creates a new message from a user with optional content.
func NewUserMessage(content ...string) *Message {
	return NewMessage(SenderUser, content...)
}

// NewUserMessage creates a new message from an assistant with optional content.
func NewAssistantMessage(content ...string) *Message {
	message := NewMessage(SenderAssistant, content...)
	return message
}

func (m Message) String() string {
	return fmt.Sprintf("%v: %v", m.Sender, m.Content)
}
