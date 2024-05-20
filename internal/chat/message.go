package chat

import (
	"fmt"
	"time"
)

type Message struct {
	Sender     Sender `json:"sender" bson:"sender"`
	Text       string `json:"text" bson:"text"`
	CreatedAt  int64  `json:"-" bson:"created_at"`
	ExchangeId string `json:"-" bson:"exchange_id,omitempty"`
	IsNew      bool   `json:"-" bson:"-"`
}

func newMessage(sender Sender, text string) *Message {
	return &Message{
		Sender:    sender,
		Text:      text,
		CreatedAt: time.Now().UnixMilli(),
		IsNew:     true,
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
	return fmt.Sprintf("%v: %v", m.Sender, m.Text)
}
