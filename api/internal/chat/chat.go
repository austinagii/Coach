package chat

// Represents the default number of messages that a chat can contain.
const DefaultMessageLimit int = 20

// Chat represents a collection of messages exchanged between a user and an assistant.
type Chat struct {
	Messages     []*Message `json:"messages" bson:"messages"`
	messageLimit int        `json:"-" bson:"-"`
}

// NewChat creates a new chat with an optional message limit. If no limit is specified,
// then the message limit will be set to the default message limit.
func NewChat(limits ...int) *Chat {
	messageLimit := DefaultMessageLimit
	// Override the default limit if at least one limit is specified.
	if len(limits) > 0 {
		messageLimit = limits[0]
	}

	return &Chat{
		Messages:     make([]*Message, 0, messageLimit),
		messageLimit: messageLimit,
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
	// Return the messages from the last exchange.
	return c.Messages[len(c.Messages)-2:]
}
