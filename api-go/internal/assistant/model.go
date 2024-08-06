package assistant

import (
	"aisu.ai/api/v2/internal/chat"
	"aisu.ai/api/v2/internal/user"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	openai "github.com/sashabaranov/go-openai"
	"go.mongodb.org/mongo-driver/mongo"
	"log/slog"
	"strings"
	"time"
)

type ModelPrompt struct {
	User *user.User `json:"user"`
	Task Task       `json:"task"`
	Chat *chat.Chat `json:"chat"`
}

type ModelResponse struct {
	Task            Task   `json:"task"`
	IsComplete      bool   `json:"is_complete"`
	UserSummary     string `json:"user_summary"`
	ResponseMessage string `json:"response"`
}

// languageModelExchange represents a single conversational exchange with a language model.
type languageModelExchange struct {
	Id            string `json:"_id,omitempty"`
	MessageId     string `json:"message_id"`
	SystemMessage string `json:"system_message"`
	Prompt        string `json:"prompt"`
	Response      string `json:"response"`
	PromptedAt    int64  `json:"prompted_at"`
	RespondedAt   int64  `json:"responded_at"`
}

type ModelExchangeRepository interface {
	Save(string, string, string, string, int64, int64) error
}

// LanguageModelExchangeRepository is responsible for auditing exchanges with a language model
// to a mongo database.
type LanguageModelExchangeRepository struct {
	collection *mongo.Collection
}

func NewLanguageModelExchangeRepository(database *mongo.Database) *LanguageModelExchangeRepository {
	return &LanguageModelExchangeRepository{collection: database.Collection("model_exchange")}
}

// Save persists a language model exchange to the database, returning any error that occurrs.
func (r *LanguageModelExchangeRepository) Save(
	messageId string,
	systemMessage string,
	prompt string,
	response string,
	promptedAt int64,
	respondedAt int64,
) error {
	exchange := &languageModelExchange{
		MessageId:     messageId,
		SystemMessage: systemMessage,
		Prompt:        prompt,
		Response:      response,
		PromptedAt:    promptedAt,
		RespondedAt:   respondedAt,
	}

	_, err := r.collection.InsertOne(context.TODO(), exchange)
	if err != nil {
		errMsg := "An error occurred while inserting a language model exchange into the database"
		slog.Error(errMsg, "error", err)
		return fmt.Errorf("%s: %w", errMsg, err)
	}
	return nil
}

// promptModel prompts the OpenAI model for a new chat message given the assistant's current context.
//
// This function encompasses three main steps:
// 1. Building the model prompt and marshalling it to the expected format.
// 2. Executing the chat completion HTTP request using the OpenAI client.
// 3. Unmarshalling the model's response to extract the response details.
//
// An error is returned if any of these steps fail.
// TODO: (Medium Priority) Return a model exchange instead of a response alone.
func promptModel(assistant *Assistant) (*ModelResponse, error) {
	// Generate a unique identifier for this exchange session
	exchangeId := uuid.NewString()

	// Create the system message using the assistant's description and task details
	systemMessageText := fmt.Sprintf("%s\n\n%s", assistantDescription, assistant.Task.Description())
	systemMessage := openai.ChatCompletionMessage{
		Role:    "system",
		Content: systemMessageText,
	}

	// Construct the user message using the user & task details along with the chat messages.
	// Though this may not seem like a traditional prompt, this user message along with the system message
	// will inform the model on what it should do.
	prompt := &ModelPrompt{
		User: assistant.User,
		Task: assistant.Task,
		Chat: assistant.Chat,
	}
	// TODO: (Low Priority) Marshall to yaml for better model understadability and size/cost reduction.
	promptBytes, err := json.Marshal(prompt)
	if err != nil {
		errMsg := "An error occurred while marshalling the model prompt to JSON"
		slog.Error(errMsg, "error", err)
		return nil, fmt.Errorf("%s: %w", errMsg, err)
	}
	userMessageText := string(promptBytes)
	userMessage := openai.ChatCompletionMessage{
		Role:    "user",
		Content: string(promptBytes),
	}

	// Record the time when the chat completion request is initiated
	initiationTime := time.Now().UnixMilli()
	// Send the chat completion request to the OpenAI API
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
	// Record the time when the chat completion request is completed
	completionTime := time.Now().UnixMilli()
	if err != nil {
		errMsg := "An error occurred while requesting a chat completion from the OpenAI API"
		slog.Error(errMsg, "error", err)
		return nil, fmt.Errorf("%s: %w", errMsg, err)
	}
	assistantMessageText := resp.Choices[0].Message.Content

	// Remove the json markdown tags from the model's response so that we just have the raw json string.
	assistantMessageText = strings.TrimPrefix(assistantMessageText, "```json\n")
	assistantMessageText = strings.TrimSuffix(assistantMessageText, "\n```")
	// Prepare the model response for unmarshalling, determining the correct task type based on the
	// objective of the request task.
	modelResponse := &ModelResponse{}
	switch prompt.Task.Objective() {
	case ObjectiveGoalCreation:
		modelResponse.Task = &GoalCreationTask{}
	case ObjectiveMilestoneCreation:
		modelResponse.Task = &MilestoneCreationTask{}
	case ObjectiveScheduleCreation:
		modelResponse.Task = &ScheduleCreationTask{}
	default:
		return nil, fmt.Errorf("Unsupported objective")
	}

	// Parse the model response.
	if err = json.Unmarshal([]byte(assistantMessageText), &modelResponse); err != nil {
		return nil, err
	}

	// Save the exchange audit asynchronously since success is not critical
	go assistant.modelExchangeRepository.Save(
		exchangeId,
		systemMessageText,
		userMessageText,
		assistantMessageText,
		initiationTime,
		completionTime,
	)

	return modelResponse, nil
}
