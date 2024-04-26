package chat

const DefaultChatMessageLimit int = 10

type Chat struct {
	Id           string     `json:"-" bson:"_id"`
	Messages     []*Message `json:"messages" bson:"messages"`
	messageLimit int        `json:"-" bson:"-"`
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
