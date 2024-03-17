package assistant

import (
	"encoding/json"
	"fmt"
	"os"
)

type TaskType int

const (
	TaskTypeGoalCreation TaskType = iota
	TaskTypeMilestoneCreation
	TaskTypeScheduleCreation
	TaskTypeChat
)

func (t TaskType) MarshalJSON() ([]byte, error) {
	var s string

	switch t {
	case TaskTypeGoalCreation:
		s = "goal_creation"
	case TaskTypeMilestoneCreation:
		s = "milestone_creation"
	case TaskTypeScheduleCreation:
		s = "schedule_creation"
	case TaskTypeChat:
		s = "chat"
	default:
		return nil, fmt.Errorf("invalid task type: %s", s)
	}
	return json.Marshal(s)
}

func (t *TaskType) UnmarshalJSON(data []byte) error {
	// Unmarshal the JSON data into a string
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	// Convert the string to a Role type, checking if it's a valid value
	switch s {
	case "goal_creation":
		*t = TaskTypeGoalCreation
	case "milestone_creation":
		*t = TaskTypeMilestoneCreation
	case "schedule_creation":
		*t = TaskTypeScheduleCreation
	case "chat":
		*t = TaskTypeChat
	default:
		return fmt.Errorf("invalid task type: %s", s)
	}
	return nil
}

func (taskType TaskType) String() (str string) {
	switch taskType {
	case TaskTypeGoalCreation:
		str = "Goal Creation"
	case TaskTypeMilestoneCreation:
		str = "Milestone Creation"
	case TaskTypeScheduleCreation:
		str = "Schedule Creation"
	case TaskTypeChat:
		str = "Chat"
	}
	return
}

type Task struct {
	Type           TaskType `json:"type"`
	IsComplete     bool     `json:"is_complete"`
	TargetEntityId *int     `json:"target,omitempty"`
}

func NewTask(taskType TaskType, targetEntityId *int) *Task {
	return &Task{
		Type:           taskType,
		IsComplete:     false,
		TargetEntityId: targetEntityId,
	}
}

var descriptionByTask = map[TaskType]string{}
var initialMessageByTask = map[TaskType]string{}

func LoadTasks() error {
	if err := LoadTaskDescriptions(); err != nil {
		return err
	}
	if err := LoadTaskInitialMessages(); err != nil {
		return err
	}
	return nil
}

func LoadTaskDescriptions() error {
	descriptionFileByTaskType := map[TaskType]string{
		TaskTypeGoalCreation:      "resources/tasks/goal_creation/description.txt",
		TaskTypeMilestoneCreation: "resources/tasks/milestone_creation/description.txt",
		TaskTypeScheduleCreation:  "resources/tasks/schedule_creation/description.txt",
		TaskTypeChat:              "resources/tasks/chat/description.txt",
	}

	for taskType, taskDescriptionFile := range descriptionFileByTaskType {
		fileContents, err := os.ReadFile(taskDescriptionFile)
		if err != nil {
			err = fmt.Errorf("Description file for task: '%s' could not be found at location: '%s'", taskType, taskDescriptionFile, err)
			return err
		}
		descriptionByTask[taskType] = string(fileContents)
	}
	return nil
}

func (task Task) GetDescription() (string, error) {
	description, ok := descriptionByTask[task.Type]
	if !ok {
		err := fmt.Errorf("No description found for task '%s'", task.Type)
		return "", err
	}
	return description, nil
}

func LoadTaskInitialMessages() error {
	initialMessageFileByTaskType := map[TaskType]string{
		TaskTypeGoalCreation:      "resources/tasks/goal_creation/initial-message.txt",
		TaskTypeMilestoneCreation: "resources/tasks/milestone_creation/initial-message.txt",
		TaskTypeScheduleCreation:  "resources/tasks/schedule_creation/initial-message.txt",
		TaskTypeChat:              "resources/tasks/chat/initial-message.txt",
	}

	for taskType, taskInitialMessageFile := range initialMessageFileByTaskType {
		fileContents, err := os.ReadFile(taskInitialMessageFile)
		if err != nil {
			err = fmt.Errorf("Initial message file for task: '%s' could not be found at location: '%s'", taskType, taskInitialMessageFile)
			return err
		}
		initialMessageByTask[taskType] = string(fileContents)
	}
	return nil
}

func (task Task) GetInitialMessage() (*Message, error) {
	initialMessageText, ok := initialMessageByTask[task.Type]
	if !ok {
		err := fmt.Errorf("No initial message found for task '%s'", task.Type)
		return NewEmptyAssistantMessage(), err
	}
	initialMessage := NewAssistantMessage(initialMessageText, nil, nil)
	return initialMessage, nil
}
