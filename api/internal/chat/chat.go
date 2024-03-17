package chat

type Chat struct {
	Id           string     `json:"-"`
	Messages     []*Message `json:"messages"`
	messageLimit int        `json:"-"`
}

func NewChat() *Chat {
	return &Chat{
		Messages:     make([]*Message, 0),
		messageLimit: 10,
	}
}

// Append adds a new message to the chat. If the chat's message limit is
// reached the oldest message will be removed from the chat and returned
func (c *Chat) Append(message *Message) *Message {
	var oldestMessage *Message
	if len(c.Messages) == c.messageLimit {
		oldestMessage = c.Messages[0]
		c.Messages = c.Messages[1:]
	}
	c.Messages = append(c.Messages, message)
	return oldestMessage
}

// GetLastMessage returns the most recent message in the chat.
func (c *Chat) GetLastMessage() *Message {
	return c.Messages[len(c.Messages)-1]
}
