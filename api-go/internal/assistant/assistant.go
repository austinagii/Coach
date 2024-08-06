package assistant

import (
	"aisu.ai/api/v2/internal/chat"
	"aisu.ai/api/v2/internal/user"
	"errors"
	"fmt"
	openai "github.com/sashabaranov/go-openai"
	"go.mongodb.org/mongo-driver/bson"
	"log/slog"
	"os"
	"time"
)

var (
	// The generic description of an assistant excluding any details about it's current task.
	assistantDescription string

	// Maps objectives to tailored user prompts that an assistant will use to start a chat.
	// e.g. If a user requests a new assistant with the objective of goal creation, the message
	// 'What's goal do you want to set?' will be used to prompt the user to start defining their goal.
	chatPromptByObjective = map[Objective]string{}

	// ErrNotIntialized indicates that the creation of a new Assistant was attempted before loading the
	// static data required.
	ErrNotIntialized = errors.New("The static data required to create a new assistant has not been loaded")
)

// InitAssistants loads the static data required to create an assistant, returning an error if one occurs.
func InitAssistants() error {
	if err := loadAssistantDescription(); err != nil {
		return err
	}
	if err := LoadObjectiveDescriptions(); err != nil {
		return err
	}
	if err := loadObjectiveChatPrompts(); err != nil {
		return err
	}
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
	filePathByObjective := map[Objective]string{
		ObjectiveGoalCreation:      "../../resources/assistant/objectives/goal_creation/initial-message.txt",
		ObjectiveMilestoneCreation: "../../resources/assistant/objectives/milestone_creation/initial-message.txt",
		ObjectiveScheduleCreation:  "../../resources/assistant/objectives/schedule_creation/initial-message.txt",
		ObjectiveChat:              "../../resources/assistant/objectives/chat/initial-message.txt",
	}

	for objective, filePath := range filePathByObjective {
		fileContents, err := os.ReadFile(filePath)
		fileContents = fileContents[:len(fileContents)-1]
		if err != nil {
			errMsg := "An error occurred while reading the chat prompt file for objective"
			slog.Error(errMsg, "objective", objective, "error", err)
			return fmt.Errorf("%s '%s': %w", errMsg, objective, err)
		}
		chatPromptByObjective[objective] = string(fileContents)
	}
	return nil
}

// Assistant is an interactive agent responsible for completing a given task by exchanging text
// based messages with a user.
type Assistant struct {
	Id                      string                  `json:"id" bson:"_id,omitempty"`
	Task                    Task                    `json:"task" bson:"task"`
	User                    *user.User              `json:"-" bson:"user"`
	Chat                    *chat.Chat              `json:"chat" bson:"chat"`
	client                  *openai.Client          `json:"-" bson:"-"`
	modelExchangeRepository ModelExchangeRepository `json:"-" bson:"-"`
	CreatedAt               int64                   `json:"-" bson:"created_at"`
	UpdatedAt               int64                   `json:"-" bson:"updated_at"`
}

// NewAssistant creates an assistant to complete a given task.
//
// An ErrNotIntialized error will be returned if the `InitAssistants` function was not executed
// successfully before attempting to call this function.
func NewAssistant(
	user *user.User,
	task Task,
	openaiClient *openai.Client,
	// TODO(Medium): Remove model exchange repo as assistant dependency.
	modelExchangeRepository ModelExchangeRepository,
) (*Assistant, error) {
	assistant := &Assistant{
		User:                    user,
		Task:                    task,
		Chat:                    chat.NewChat(),
		client:                  openaiClient,
		modelExchangeRepository: modelExchangeRepository,
		CreatedAt:               time.Now().Local().UnixMilli(),
	}

	// Start the assistant's chat using the prompt defined for the specified objective.
	chatPromptText, ok := chatPromptByObjective[task.Objective()]
	if !ok {
		errorMsg := "No chat prompt found for objective"
		slog.Error(errorMsg, "objective", task.Objective().String())
		return nil, fmt.Errorf("%s '%s'", errorMsg, task.Objective().String())
	}
	chatPrompt := chat.NewAssistantMessage(chatPromptText)
	assistant.Chat.Append(chatPrompt)
	return assistant, nil
}

// Init configures an existing assistant to use the provided openai client and model exchange repository.
func (assistant *Assistant) Init(
	openaiClient *openai.Client,
	modelExchangeRepository *LanguageModelExchangeRepository,
) {
	assistant.client = openaiClient
	assistant.modelExchangeRepository = modelExchangeRepository
}

// Respond generates a response to the provided message and updates the assistant's state accordingly.
// It returns the generated response message or an error if one occurs.
func (assistant *Assistant) Respond(message *chat.Message) (*chat.Message, error) {
	// Add the incoming message to the chat history.
	assistant.Chat.Append(message)

	// Generate a response from the model based on the current chat context.
	modelResponse, err := promptModel(assistant)
	if err != nil {
		errMsg := "An error occurred while requesting a response from the model"
		slog.Error(errMsg, "error", err)
		return nil, fmt.Errorf("%s: %w", errMsg, err)
	}

	// Add the model's response message to the chat history.
	assistantMessage := chat.NewAssistantMessage(modelResponse.ResponseMessage)
	assistant.Chat.Append(assistantMessage)

	// Update the assistant's context based on the type of task completed.
	assistant.Task = modelResponse.Task
	assistant.User.Summary = modelResponse.UserSummary
	// TODO: (Low Priority) Move each case into it's own handler function that updates the assistant context.
	if modelResponse.IsComplete {
		switch modelResponse.Task.Objective() {
		case ObjectiveGoalCreation:
			// Avoid unnecessary error handling on type conversion since type setting is explicit in the
			// `promptModel` function.
			task := modelResponse.Task.(*GoalCreationTask)
			assistant.User.AddNewGoal(task.Goal)
			assistant.Task = NewMilestoneCreationTask(task.Goal.Id)
		case ObjectiveMilestoneCreation:
			task := modelResponse.Task.(*MilestoneCreationTask)
			goal, err := assistant.User.GetGoalById(task.GoalId)
			if err != nil {
				return nil, err
			}
			goal.Milestones = task.Milestones
			assistant.Task = NewScheduleCreationTask()
		case ObjectiveScheduleCreation:
			task := modelResponse.Task.(*ScheduleCreationTask)
			assistant.User.Schedule = task.Schedule
		}
	}

	return assistantMessage, nil
}

// TODO: (Medium Priority) Add comments for understandability.
func (assistant *Assistant) UnmarshalBSON(data []byte) error {
	var tempAssistant struct {
		Id        string     `bson:"_id"`
		Task      bson.Raw   `bson:"task"`
		User      *user.User `bson:"user"`
		Chat      *chat.Chat `bson:"chat"`
		CreatedAt int64      `bson:"created_at"`
		UpdatedAt int64      `bson:"updated_at"`
	}
	if err := bson.Unmarshal(data, &tempAssistant); err != nil {
		return err
	}

	tempTask := &BaseTask{}
	if err := bson.Unmarshal(tempAssistant.Task, tempTask); err != nil {
		return err
	}
	var t Task
	switch tempTask.Objective() {
	case ObjectiveGoalCreation:
		t = &GoalCreationTask{}
	case ObjectiveMilestoneCreation:
		t = &MilestoneCreationTask{}
	case ObjectiveScheduleCreation:
		t = &ScheduleCreationTask{}
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
	assistant.CreatedAt = tempAssistant.CreatedAt
	assistant.UpdatedAt = tempAssistant.UpdatedAt
	return nil
}
