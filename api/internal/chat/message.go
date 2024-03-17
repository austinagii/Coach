package chat

import (
	"aisu.ai/api/v2/internal/user"
	"fmt"
)

type Message struct {
	Sender Sender     `json:"-"`
	User   *user.User `json:"user,omitempty"`
	Task   *Task      `json:"task,omitempty"`
	Text   string     `json:"text"`
}

func NewEmptyUserMessage() *Message {
	return &Message{Sender: SenderUser}
}

func NewUserMessage(text string) *Message {
	return &Message{
		Sender: SenderUser,
		Text:   text,
	}
}

func NewEmptyAssistantMessage() *Message {
	return &Message{Sender: SenderAssistant}
}

func NewAssistantMessage(text string, user *user.User, task *Task) *Message {
	return &Message{
		Sender: SenderAssistant,
		User:   user,
		Task:   task,
		Text:   text,
	}
}

func (m Message) String() string {
	return fmt.Sprintf("%v: %v", m.Sender, m.Text)
}
