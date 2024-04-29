package assistant

import (
	"aisu.ai/api/v2/internal/assistant/chat"
	"aisu.ai/api/v2/internal/assistant/task"
	"aisu.ai/api/v2/internal/user"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	openai "github.com/sashabaranov/go-openai"
	"log/slog"
	"os"
	"strings"
	"time"
)

// Maps objectives to tailored initial messages for starting conversations.
var chatPromptByObjective = map[task.Objective]string{}

// Assistant is an interactive agent responsible for completing a specified
// task by sending and receiving text based messages.
//
// An assistant manages the current task and chat states, updating both as
// new messages are received and processed. The assistant provides responses
// by wrapping an OpenAI model via the `openai.Client` struct and prompts
// that model to produce a response to a given message.
type Assistant struct {
	Id                      string                           `json:"id" bson:"_id"`
	Task                    *task.Task                       `json:"task"`
	description             string                           `json:"-"`
	User                    *user.User                       `json:"-"`
	Chat                    *chat.Chat                       `json:"chat"`
	client                  *openai.Client                   `json:"-"`
	modelExchangeRepository *LanguageModelExchangeRepository `json:"-"`
}

// NewAssistant creates an assistant to complete a given task.
//
// When a new assistant is created, it is initalized with a description
// and initial chat message based on the specified task. The description
// is a combination of the assistant's generic description plus the
// description of the task and the initial message is a pre-defined starter
// message relevant to the task's objective.
func NewAssistant(
	openaiClient *openai.Client,
	modelExchangeRepository *LanguageModelExchangeRepository,
	task *task.Task,
) (*Assistant, error) {
	// TODO: Add checks to ensure the client is available for use.
	assistant := &Assistant{
		client:                  openaiClient,
		modelExchangeRepository: modelExchangeRepository,
		Task:                    task,
	}

	chatPromptText, ok := chatPromptByObjective[task.Objective]
	if !ok {
		errorMsg := "No initial message found for objective"
		slog.Error(errorMsg, "objective", task.Objective.String())
		fmt.Errorf("%s '%s'", errorMsg, task.Objective.String())
	}
	chatPrompt := chat.NewAssistantMessage(chatPromptText)
	assistant.Chat.Append(chatPrompt)
	return assistant, nil
}

// Respond provides a reply message to the given message or an error if one occurrs.
//
// Respond takses into account information about the user to provide a context
// sensitive response to the given message.
func (assistant *Assistant) Respond(message *chat.Message) (*chat.Message, error) {
	assistant.Chat.Append(message)
	assistantMessage, err := assistant.promptModel()
	if err != nil {
		errMsg := "An error occurred while requesting a response from the model"
		slog.Error(errMsg, "error", err)
		return nil, fmt.Errorf("%s: %w", errMsg, err)
	}
	assistant.Chat.Append(assistantMessage)
	return assistantMessage, nil
}

func (assistant *Assistant) promptModel() (*chat.Message, error) {
	// Generate a unique ID to track this conversational exchange.
	exchangeId := uuid.NewString()

	systemMessageText := fmt.Sprintf("%s\n\n%s", assistant.description, assistant.Task.Description)
	systemMessage := openai.ChatCompletionMessage{
		Role:    "system",
		Content: systemMessageText,
	}

	// TODO: Marshall to yaml for better model understadability and size/cost reduction.
	userMessageBytes, err := json.Marshal(ChatPrompt{
		user: assistant.User,
		task: assistant.Task,
		chat: assistant.Chat,
	})
	if err != nil {
		errMsg := "An error occurred while marshalling the chat prompt to JSON"
		slog.Error(errMsg, "error", err)
		return nil, fmt.Errorf("%s: %w", errMsg, err)
	}
	userMessageText := string(userMessageBytes)
	userMessage := openai.ChatCompletionMessage{
		Role:    "user",
		Content: string(userMessageText),
	}

	initiationTime := time.Now().UnixMilli()
	resp, err := assistant.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    openai.GPT4TurboPreview,
			Messages: []openai.ChatCompletionMessage{systemMessage, userMessage},
			ResponseFormat: &openai.ChatCompletionResponseFormat{
				Type: openai.ChatCompletionResponseFormatTypeJSONObject,
			},
		},
	)
	completionTime := time.Now().UnixMilli()
	if err != nil {
		errMsg := "An error occurred while requesting a chat completion from the OpenAI API"
		slog.Error(errMsg, "error", err)
		return nil, fmt.Errorf("%s: %w", errMsg, err)
	}
	assistantMessageText := resp.Choices[0].Message.Content

	// Fire and forget the exchange audit. Success here is't critical
	go assistant.modelExchangeRepository.Save(
		exchangeId,
		systemMessageText,
		userMessageText,
		assistantMessageText,
		initiationTime,
		completionTime,
	)

	assistantMessage := chat.NewEmptyAssistantMessage()
	assistantMessageText = strings.TrimPrefix(assistantMessageText, "```json\n")
	assistantMessageText = strings.TrimSuffix(assistantMessageText, "\n```")
	err = json.Unmarshal([]byte(assistantMessageText), assistantMessage)
	if err != nil {
		errMsg := "An error occurred while unmarshalling the OpenAI API model response"
		slog.Error(errMsg, "error", err)
		return nil, fmt.Errorf("%s: %w", errMsg, err)
	}
	return assistantMessage, nil
}

// loadChatPrompts loads the chat prompts used by a newly created assistant to
// start a conversation from disk, returning an error if the file could not be
// read.
func loadChatPrompts() error {
	filePathByObjective := map[task.Objective]string{
		task.ObjectiveGoalCreation:      "resources/assistant/objectives/goal_creation/initial-message.txt",
		task.ObjectiveMilestoneCreation: "resources/assistant/objectives/milestone_creation/initial-message.txt",
		task.ObjectiveScheduleCreation:  "resources/assistant/objectives/schedule_creation/initial-message.txt",
		task.ObjectiveChat:              "resources/assistant/objectives/chat/initial-message.txt",
	}

	for objective, filePath := range filePathByObjective {
		fileContents, err := os.ReadFile(filePath)
		if err != nil {
			errMsg := fmt.Sprintf("An error occurred while reading the initial message file for objective '%s'", objective)
			slog.Error(errMsg, "error", err)
			return fmt.Errorf("%s: %w", errMsg, err)
		}
		chatPromptByObjective[objective] = string(fileContents)
	}
	return nil
}
