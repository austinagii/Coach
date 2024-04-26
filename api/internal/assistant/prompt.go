package assistant

import (
	"aisu.ai/api/v2/internal/assistant/task"
	"aisu.ai/api/v2/internal/chat"
	"aisu.ai/api/v2/internal/user"
)

type ChatPrompt struct {
	user *user.User
	task *task.Task
	chat *chat.Chat
}

type ModelResponse struct {
	UserSummary     string `json:"user_summary"`
	IsComplete      bool   `json:"is_complete"`
	ResponseMessage string `json:"response_message"`
}

type GoalCreationResponse struct {
	ModelResponse
	Goal *user.Goal `json:"goal"`
}

type MilestoneCreationResponse struct {
	ModelResponse
	Goal *user.Goal `json:"goal"`
}

type ScheduleCreationResponse struct {
	ModelResponse
	Goal *user.Schedule `json:"goal"`
}
