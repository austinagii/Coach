package assistant

import (
	"aisu.ai/api/v2/internal/chat"
	"context"
	"encoding/json"
	"fmt"
	openai "github.com/sashabaranov/go-openai"
	"log"
	"os"
	"strings"
)

var APIKey string = os.Getenv("OPENAI_API_KEY")

type Assistant struct {
	client      *openai.Client
	Description string
	Chat        *chat.Chat
	repository  *chat.ChatRepository
}

func NewAssistant(task *chat.Task, chatId string) *Assistant {
	assistant := &Assistant{
		client:     openai.NewClient(APIKey),
		repository: chat.NewChatRepository(),
	}

	if assistantMessage.Task.IsComplete {
		if assistantMessage.Task.TaskType == TaskTypeGoalCreation {
			goalId = goalRepository.Save(Task.Output)
			assistantMessage.Task.TaskType = TaskTypeMilestoneCreation
			assistantMessage.Task.TargetEntityId = goalId
		} else if assistantMessage.Task.TaskType == TaskTypeMilestoneCreation {

		}
	}
	log.Print("Assistant defined")
	assistant.Chat = assistant.repository.Get(chatId)
	assistant.Chat.Id = chatId
	log.Printf("Loaded chat: %v", assistant.Chat)
	if err := assistant.SetTask(task); err != nil {
		log.Print("Here")
		log.Fatal(err)
	}
	return assistant
}

func LoadAssistant(chatId string) *Assistant {
	assistant := &Assistant{
		client:     openai.NewClient(APIKey),
		repository: chat.NewChatRepository(),
	}
	assistant.Chat = assistant.repository.Get(chatId)
	assistant.Chat.Id = chatId

	lastTask := assistant.Chat.GetLastMessage().Task
	if lastTask.IsComplete {
		var currentTaskType TaskType
		switch lastTask.TaskType {
		case TaskTypeGoalCreation:
			currentTaskType = TaskTypeMilestoneCreation
		case TaskTypeMilestoneCreation:
			currentTaskType = TaskTypeScheduleCreation
		case TaskTypeScheduleCreation:
			currentTaskType = TaskTypeChat
		}
		assistant.SetTask(nextTask)
	} else {
		assistant.SetTask(assistant.Chat.GetLastMessage().Task)
	}

	log.Print("Assistant defined")
	log.Printf("Loaded chat: %v", assistant.Chat)
	if err := assistant.SetTask(task); err != nil {
		log.Print("Here")
		log.Fatal(err)
	}
	return assistant
}

func (assistant *Assistant) SetTask(task *chat.Task) error {
	// TODO: Refactor to only load assistant description once.
	assistantDescriptionFilePath := "resources/assistant-description.txt"
	fileContents, err := os.ReadFile(assistantDescriptionFilePath)
	if err != nil {
		err := fmt.Errorf("Assistant description file could not be found at location: '%s'", assistantDescriptionFilePath)
		return err
	}
	assistantDescriptionTemplate := string(fileContents)
	taskDescription, err := task.GetDescription()
	if err != nil {
		return err
	}

	assistant.Description = strings.Replace(assistantDescriptionTemplate, "%task-description%", taskDescription, -1)
	return nil
}

func (assistant *Assistant) Respond(message *chat.Message) (*chat.Message, error) {
	assistant.Chat.Append(message)

	resp, err := assistant.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    openai.GPT4TurboPreview,
			Messages: toChatCompletionMessages(assistant.Description, assistant.Chat),
			ResponseFormat: &openai.ChatCompletionResponseFormat{
				Type: openai.ChatCompletionResponseFormatTypeJSONObject,
			},
		},
	)
	if err != nil {
		return nil, err
	}
	assistantMessage, err := parseAssistantResponse(resp.Choices[0].Message.Content)
	if err != nil {
		return nil, err
	}

	if assistantMessage.Task.IsComplete {
		switch assistantMessage.Task.Type {
		case TaskTypeGoalCreation:
			goal := assistantMessage.Task.Result.(*Goal)
			goalId := goalRepository.Save(goal)
		case TaskTypeScheduleCreation:
			milestones = assistantMessage.Task.Result.([]*Milestone)
			milestones = goalRepository.SaveMilestones(goalId, milestones)
		case TaskTypeScheduleCreation:
			// is his the first schedule?
			schedule = assistant.Task.Result(*Schedule)
			scheduleId = scheduleRepository.SaveMilestones(schedule)
		}
	}

	assistant.Chat.Append(assistantMessage)
	assistant.repository.SaveNewMessages(assistant.Chat.Id, message, assistantMessage)
	return assistantMessage, nil
}

var roleBySender = map[chat.Sender]string{
	chat.SenderUser:      "user",
	chat.SenderAssistant: "assistant",
	chat.SenderSystem:    "system",
}

func getRole(sender chat.Sender) string {
	role, ok := roleBySender[sender]
	if !ok {
		// Do something here.
	}
	return role
}

func toChatCompletionMessages(assistantDescription string, userChat *chat.Chat) []openai.ChatCompletionMessage {
	var messages []openai.ChatCompletionMessage

	messages = append(messages, openai.ChatCompletionMessage{
		Role:    getRole(chat.SenderSystem),
		Content: assistantDescription,
	})

	for _, message := range userChat.Messages {
		messageJson, err := json.Marshal(message)
		if err != nil {
			// Do something here.
		}
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    getRole(message.Sender),
			Content: string(messageJson),
		})
	}
	return messages
}

func parseAssistantResponse(jsonResponse string) (*chat.Message, error) {
	message := chat.NewEmptyAssistantMessage()
	fmt.Println(jsonResponse)
	jsonResponse = strings.TrimPrefix(jsonResponse, "```json\n")
	jsonResponse = strings.TrimSuffix(jsonResponse, "\n```")
	err := json.Unmarshal([]byte(jsonResponse), message)
	if err != nil {
		return nil, err
	}
	return message, nil
}
