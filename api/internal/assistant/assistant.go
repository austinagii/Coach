package assistant

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"
	"time"
	"uuid"

	"aisu.ai/api/v2/internal/assistant/chat"
	"aisu.ai/api/v2/internal/assistant/task"
	"aisu.ai/api/v2/internal/user"
	openai "github.com/sashabaranov/go-openai"
)

// The initial message to be
var initialMessageByOjective = map[task.Objective]string{}
var assistantDescription string

// TODO: Properly initialize model exchange repoistory
var modelExchangeRepository = NewLanguageModelExchangeRepository()

// Assistant is an interactive agent responsible for completing a specified
// task by sending and receiving text based messages.
//
// An assistant manages the current task and chat states, updating both as
// new messages are received and processed. The assistant provides responses
// by wrapping an OpenAI model via the `openai.Client` struct and prompts
// that model to produce a response to a given message.
type Assistant struct {
	Id          string     `json:"id" bson:"_id"`
	Task        *task.Task `json:"task"`
	description string     `json:"-"`
	// Want to make the user struct readonly so that the assistant doesn't
	// update the users goals accidentally. But would still like the User
	// description to be updated as the assistant interacts with the user.
	User *user.User `json:"-"`
	// Each time a message is sent to openai that full prompt and response
	// should be saved to the database for auditing. Each time a prompt is sent
	// a unique id should be generated and the corresponding chat message
	// produced should be linked back to that prompt's ID.
	Chat   *chat.Chat     `json:"chat"`
	client *openai.Client `json:"-"`
}

// NewAssistant creates an assistant to complete the specified task.
//
// When a new assistant is created, it is initalized with a description
// and initial chat message based on the specified task. The description
// is a combination of the assistant's generic description plus the
// description of the task and the initial message is a pre-defined starter
// message relevant to the task's objective.
func NewAssistant(
	openaiClient *openai.Client,
	task *task.Task,
) *Assistant {
	assistant := &Assistant{
		client: openaiClient,
		Task:   task,
	}
	return assistant
}

func (assistant *Assistant) getInitialMessage(objective task.Objective) (*chat.Message, error) {
	initialMessageText, ok := initialMessageByOjective[objective]
	if !ok {
		err := fmt.Errorf("No initial message found for task '%s'", assistant.task.Objective)
		return chat.NewEmptyAssistantMessage(), err
	}
	return chat.NewAssistantMessage(initialMessageText), nil
}

func (assistant *Assistant) Respond(message *chat.Message) (*chat.Message, error) {
	assistant.Chat.Append(message)
	assistantMessage, err := assistant.requestMessageResponseFromModel()
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
	go modelExchangeRepository.Save(
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

// loadDescription loads the agent's description from the description file.
func loadDescription() error {
	filePath := "resources/assistant/description.txt"
	fileContents, err := os.ReadFile(filePath)
	if err != nil {
		errMsg := "An error occurred while attempting to read the assistant description file"
		slog.Error(errMsg, "error", err)
		return fmt.Errorf("%s: %w", filePath, err)
	}

	assistantDescription = string(fileContents)
	return nil
}

// loadTaskInitialMessages loads the first message to be used by a newly created assistant for each objective.
//
// Each message is loaded from it's config file and an error is returned if the message file could
// not be read.
func loadTaskInitialMessages() error {
	filePathByObjective := map[task.Objective]string{
		task.ObjectiveGoalCreation:      "resources/objectives/goal_creation/initial-message.txt",
		task.ObjectiveMilestoneCreation: "resources/objectives/milestone_creation/initial-message.txt",
		task.ObjectiveScheduleCreation:  "resources/objectives/schedule_creation/initial-message.txt",
		task.ObjectiveChat:              "resources/objectives/chat/initial-message.txt",
	}

	for objective, filePath := range filePathByObjective {
		fileContents, err := os.ReadFile(filePath)
		if err != nil {
			errMsg := fmt.Sprintf("An error occurred while reading the initial message file for objective '%s'", objective)
			slog.Error(errMsg, "error", err)
			return fmt.Errorf("%s: %w", errMsg, err)
		}
		initialMessageByOjective[objective] = string(fileContents)
	}
	return nil
}
