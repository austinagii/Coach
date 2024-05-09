package assistant

import (
	"aisu.ai/api/v2/internal/assistant/chat"
	"aisu.ai/api/v2/internal/assistant/task"
	"aisu.ai/api/v2/internal/user"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	openai "github.com/sashabaranov/go-openai"
	"go.mongodb.org/mongo-driver/bson"
	"log/slog"
	"os"
	"strings"
	"time"
)

// The shared description of all assistants excluding task info.
var assistantDescription string

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
	Id                      string                           `json:"id" bson:"_id,omitempty"`
	Task                    task.Task                        `json:"task"`
	User                    *user.User                       `json:"-"`
	Chat                    *chat.Chat                       `json:"chat"`
	client                  *openai.Client                   `json:"-"`
	modelExchangeRepository *LanguageModelExchangeRepository `json:"-"`
}

// initiAssistant loads the chat prompt defined for each objective from disk,
// returning an error if the file could not be
func InitAssistants() error {
	if err := loadAssistantDescription(); err != nil {
		return err
	}
	if err := task.LoadObjectiveDescriptions(); err != nil {
		return err
	}
	if err := loadObjectiveChatPrompts(); err != nil {
		return err
	}
	return nil
}

// NewAssistant creates an assistant to complete a given task.
//
// When a new assistant is created, it is initalized with a description
// and initial chat message based on the specified task. The description
// is a combination of the assistant's generic description plus the
// description of the task and the initial message is a pre-defined starter
// message relevant to the task's objective.
func NewAssistant(
	user *user.User,
	task task.Task,
	openaiClient *openai.Client,
	modelExchangeRepository *LanguageModelExchangeRepository,
) (*Assistant, error) {
	// TODO: Add checks to ensure the client is available for use.
	assistant := &Assistant{
		User:                    user,
		Task:                    task,
		Chat:                    chat.NewChat(),
		client:                  openaiClient,
		modelExchangeRepository: modelExchangeRepository,
	}

	chatPromptText, ok := chatPromptByObjective[task.Objective()]
	for key := range chatPromptByObjective {
		slog.Error("Found key", "key", key)
	}
	if !ok {
		errorMsg := "No initial message found for objective"
		slog.Error(errorMsg, "objective", task.Objective().String())
		return nil, fmt.Errorf("%s '%s'", errorMsg, task.Objective().String())
	}
	chatPrompt := chat.NewAssistantMessage(chatPromptText, "")
	assistant.Chat.Append(chatPrompt)
	return assistant, nil
}

func (assistant *Assistant) Init(
	openaiClient *openai.Client,
	modelExchangeRepository *LanguageModelExchangeRepository,
) {
	assistant.client = openaiClient
	assistant.modelExchangeRepository = modelExchangeRepository
}

// Respond provides a reply message to the given message or an error if one occurrs.
//
// Respond takses into account information about the user to provide a context
// sensitive response to the given message.
func (assistant *Assistant) Respond(message *chat.Message) (*chat.Message, error) {
	assistant.Chat.Append(message)
	exchangeId, modelPrompt, err := assistant.promptModel()
	if err != nil {
		errMsg := "An error occurred while requesting a response from the model"
		slog.Error(errMsg, "error", err)
		return nil, fmt.Errorf("%s: %w", errMsg, err)
	}

	assistant.User.Summary = modelPrompt.UserSummary
	if modelPrompt.IsComplete {
		if modelPrompt.Task.Objective() == task.ObjectiveGoalCreation {
			t, ok := modelPrompt.Task.(*task.GoalCreationTask)
			if !ok {
				return nil, errors.New("Failed to convert task with objective 'goal_creation' to expected struct 'GoalCreationTask'")
			}
			assistant.User.AddNewGoal(t.Goal)
			assistant.Task = task.NewMilestoneCreationTask(t.Goal.Id)
		} else if modelPrompt.Task.Objective() == task.ObjectiveMilestoneCreation {
			task, ok := modelPrompt.Task.(*task.MilestoneCreationTask)
			if !ok {
				return nil, errors.New("Failed to convert task with objective 'milestone_creation' to expected struct 'MilestoneCreationTask'")
			}
			goal, err := assistant.User.GetGoalById(task.GoalId)
			if err != nil {
				return nil, err
			}
			goal.Milestones = task.Milestones
		} else {
			// TODO: Use more descriptive error message.
			return nil, errors.New("Unsupported objective")
		}
	}

	assistantMessage := chat.NewAssistantMessage(modelPrompt.ResponseMessage, exchangeId)
	assistant.Chat.Append(assistantMessage)
	return assistantMessage, nil
}

func (assistant *Assistant) promptModel() (string, *ModelResponse, error) {
	// Generate a unique ID to track this conversational exchange.
	exchangeId := uuid.NewString()

	systemMessageText := fmt.Sprintf("%s\n\n%s", assistantDescription, assistant.Task.Description())
	systemMessage := openai.ChatCompletionMessage{
		Role:    "system",
		Content: systemMessageText,
	}

	// TODO: Marshall to yaml for better model understadability and size/cost reduction.
	modelPrompt := NewModelPrompt(assistant.User, assistant.Task, assistant.Chat)
	userMessageBytes, err := json.Marshal(modelPrompt)
	if err != nil {
		errMsg := "An error occurred while marshalling the model prompt to JSON"
		slog.Error(errMsg, "error", err)
		return "", nil, fmt.Errorf("%s: %w", errMsg, err)
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
		return "", nil, fmt.Errorf("%s: %w", errMsg, err)
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

	assistantMessageText = strings.TrimPrefix(assistantMessageText, "```json\n")
	assistantMessageText = strings.TrimSuffix(assistantMessageText, "\n```")
	modelResponse := NewEmptyModelResponse()
	var tempModelResponse struct {
		Task        json.RawMessage `json:"task"`
		IsComplete  bool            `json:"is_complete"`
		UserSummary string          `json:"user_summary"`
		Response    string          `json:"response"`
	}
	err = json.Unmarshal([]byte(assistantMessageText), &tempModelResponse)
	if err != nil {
		return "", nil, err
	}
	modelResponse.IsComplete = tempModelResponse.IsComplete
	modelResponse.UserSummary = tempModelResponse.UserSummary
	modelResponse.ResponseMessage = tempModelResponse.Response
	switch modelPrompt.Task.Objective() {
	case task.ObjectiveGoalCreation:
		task := &task.GoalCreationTask{}
		if err := json.Unmarshal(tempModelResponse.Task, task); err != nil {
			return "", nil, err
		}
		modelResponse.Task = task
	case task.ObjectiveMilestoneCreation:
		task := &task.MilestoneCreationTask{}
		if err := json.Unmarshal(tempModelResponse.Task, task); err != nil {
			return "", nil, err
		}
		modelResponse.Task = task
	default:
		return "", nil, fmt.Errorf("Unsupported objective")
	}
	if err != nil {
		errMsg := "An error occurred while unmarshalling the OpenAI API model response"
		slog.Error(errMsg, "error", err)
		return "", nil, fmt.Errorf("%s: %w", errMsg, err)
	}
	return exchangeId, modelResponse, nil
}

func (assistant *Assistant) UnmarshalBSON(data []byte) error {
	var tempAssistant struct {
		Id   string     `bson:"_id"`
		Task bson.Raw   `bson:"task"`
		User *user.User `bson:"user"`
		Chat *chat.Chat `bson:"chat"`
	}
	if err := bson.Unmarshal(data, &tempAssistant); err != nil {
		return err
	}

	tempTask := &task.BaseTask{}
	if err := bson.Unmarshal(tempAssistant.Task, tempTask); err != nil {
		return err
	}
	slog.Info("Data: ", "temp task", tempAssistant.Task)
	var t task.Task
	switch tempTask.Objective() {
	case task.ObjectiveGoalCreation:
		t = &task.GoalCreationTask{}
	case task.ObjectiveMilestoneCreation:
		t = &task.MilestoneCreationTask{}
	default:
		return fmt.Errorf("No task found for objective %s", tempTask.Objective().String())
	}
	if err := bson.Unmarshal(tempAssistant.Task, t); err != nil {
		return err
	}

	assistant.Id = tempAssistant.Id
	assistant.Task = t
	assistant.User = tempAssistant.User
	assistant.Chat = tempAssistant.Chat
	return nil
}

func loadAssistantDescription() error {
	fileContents, err := os.ReadFile("../../resources/assistant/description.txt")
	if err != nil {
		return fmt.Errorf("An error occurred while loading the assistant description: %w", err)
	}
	assistantDescription = string(fileContents)
	return nil
}

func loadObjectiveChatPrompts() error {
	filePathByObjective := map[task.Objective]string{
		task.ObjectiveGoalCreation:      "../../resources/assistant/objectives/goal_creation/initial-message.txt",
		task.ObjectiveMilestoneCreation: "../../resources/assistant/objectives/milestone_creation/initial-message.txt",
		task.ObjectiveScheduleCreation:  "../../resources/assistant/objectives/schedule_creation/initial-message.txt",
		task.ObjectiveChat:              "../../resources/assistant/objectives/chat/initial-message.txt",
	}

	for objective, filePath := range filePathByObjective {
		slog.Error("Loading prompt for objective", "objective", objective)
		fileContents, err := os.ReadFile(filePath)
		if err != nil {
			errMsg := "An error occurred while reading the chat prompt file for objective"
			slog.Error(errMsg, "objective", objective, "error", err)
			return fmt.Errorf("%s '%s': %w", errMsg, objective, err)
		}
		chatPromptByObjective[objective] = string(fileContents)
	}
	return nil
}
