package chat

const DefaultChatMessageLimit int = 20

type Chat struct {
	Messages     []*Message `json:"messages" bson:"messages"`
	messageLimit int        `json:"-" bson:"-"`
}

func NewChat() *Chat {
	return &Chat{
		Messages:     make([]*Message, 0, DefaultChatMessageLimit),
		messageLimit: DefaultChatMessageLimit,
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

// TODO: Add a way to identify and return new messages
// could be as simple as adding a 'is_new' flag to each message
func (c *Chat) GetNewMessages() []*Message {
	return make([]*Message, 1)
}
