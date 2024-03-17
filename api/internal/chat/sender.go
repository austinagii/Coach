package chat

type Sender int

const (
	SenderUser Sender = iota
	SenderAssistant
	SenderSystem
)

func (sender Sender) String() string {
	var s string
	switch sender {
	case SenderUser:
		s = "User"
	case SenderAssistant:
		s = "Ana"
	case SenderSystem:
		s = "System"
	}
	return s
}
