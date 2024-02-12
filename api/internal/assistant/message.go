package assistant 

import "fmt"

type Sender int

const (
  SenderUser Sender = iota
  SenderAssistant
  SenderSystem
)

var senderRoles = []string{"user", "assistant", "system"}
var senderNames = []string{"User", "ANa", "System"}

func (sender Sender) getRole() string {
  var role string = ""
  
  if int(sender) >= 0 && int(sender) < len(senderRoles) {
    role = senderRoles[sender]
  }
  return role
}

func (sender Sender) getName() string {
  var name string = ""
  
  if int(sender) >= 0 && int(sender) < len(senderNames) {
    name = senderNames[sender]
  }
  return name
}

type Message struct {
  Sender Sender 
  Text string   `json:"text"`
}

func NewUserMessage(text string) *Message {
  return &Message{Sender: SenderUser, Text: text}
}

func NewAssistantMessage(text string) *Message {
  return &Message{Sender: SenderAssistant, Text: text}
}

func NewSystemMessage(text string) *Message {
  return &Message{Sender: SenderSystem, Text: text}
}

func (m Message) String() string {
  return fmt.Sprintf("%v: %v", m.Sender.getName(), m.Text)
}
