package assistant

import (
  "container/list"
  "context"
  "encoding/json"
  "fmt"
  "strings"
  openai "github.com/sashabaranov/go-openai"
  "os"
)

var APIKey string = os.Getenv("OPENAI_API_KEY")

type Assistant struct {
  Type      AssistantType
  Chat      *list.List
  client    *openai.Client
}

type AssistantResponse struct {
  Goal        string    `json:"goal"`
  IsSmartGoal bool      `json:"is_smart_goal"`
  Message     *Message  `json:"message"`
}

func NewAssistant(description string) *Assistant {
  assistant := &Assistant{
    // TODO: Load API key from config in a testable way.
    client: openai.NewClient(APIKey),
    Chat: list.New(),
  }
  assistant.Chat.PushBack(NewSystemMessage(description))
  assistant.Chat.PushBack(NewAssistantMessage("Hey! What's your goal?"))
  return assistant
}

func (assistant *Assistant) Respond(message *Message) (*AssistantResponse, error) {
  assistant.Chat.PushBack(message)

  resp, err := assistant.client.CreateChatCompletion(
    context.Background(),
    openai.ChatCompletionRequest {
      Model: openai.GPT4TurboPreview,
      Messages: toChatCompletionMessages(assistant.Chat),
      ResponseFormat: &openai.ChatCompeltionResponseFormat{
        Type: openai.ChatCompletionResponseFormatTypeJSONObject,
      },
    },
  )
  if err != nil {
    return nil, err
  }
  // fmt.Println(resp.Choices[0].Message.Content)  
  assistantResponse, err := parseAssistantResponse(resp.Choices[0].Message.Content)
  if err != nil {
    fmt.Printf("-----------------------------------------\n")
    fmt.Println("An error occurred. Printing the conversation history")
    for iterator := assistant.Chat.Front(); iterator != nil; iterator = iterator.Next() {
      fmt.Println(iterator.Value)
    }
    return nil, err
  }
  assistant.Chat.PushBack(assistantResponse.Message)
  return assistantResponse, nil
}

func toChatCompletionMessages(chat *list.List) []openai.ChatCompletionMessage {
  var messages []openai.ChatCompletionMessage

  for iterator := chat.Front(); iterator != nil; iterator = iterator.Next() {
    message := iterator.Value.(*Message)
    messages = append(messages, openai.ChatCompletionMessage{
      Role: message.Sender.getRole(),
      Content: message.Text,
    })
  }
  return messages
}

func parseAssistantResponse(jsonResponse string) (*AssistantResponse, error) {
  assistantResponse := &AssistantResponse{
    Message: &Message{
      Sender: SenderAssistant,
    }, 
  } 

  jsonResponse = strings.TrimPrefix(jsonResponse, "```json\n")
  jsonResponse = strings.TrimSuffix(jsonResponse, "\n```")
  err := json.Unmarshal([]byte(jsonResponse), assistantResponse)
  if err != nil {
    return nil, err
  }
  return assistantResponse, nil
}
