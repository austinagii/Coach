package chat

import (
	"fmt"
	"time"
)

type Message struct {
	Sender     Sender `json:"sender" bson:"sender"`
	Content    string `json:"content" bson:"content"`
	CreatedAt  int64  `json:"-" bson:"created_at"`
	ExchangeId string `json:"-" bson:"exchange_id,omitempty"`
}

func newMessage(sender Sender, content string) *Message {
	return &Message{
		Sender:    sender,
		Content:   content,
		CreatedAt: time.Now().UnixMilli(),
	}
}

func NewEmptyUserMessage() *Message {
	return newMessage(SenderUser, "")
}

func NewUserMessage(text string) *Message {
	return newMessage(SenderUser, text)
}

func NewEmptyAssistantMessage() *Message {
	return newMessage(SenderAssistant, "")
}

func NewAssistantMessage(text string, exchangeId string) *Message {
	message := newMessage(SenderAssistant, text)
	message.ExchangeId = exchangeId
	return message
}

func (m Message) String() string {
	return fmt.Sprintf("%v: %v", m.Sender, m.Content)
}
