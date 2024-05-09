package assistant

import (
	"aisu.ai/api/v2/internal/assistant/chat"
	"aisu.ai/api/v2/internal/assistant/task"
	"aisu.ai/api/v2/internal/user"
)

type ModelPrompt struct {
	User *user.User `json:"user"`
	Task task.Task  `json:"task"`
	Chat *chat.Chat `json:"chat"`
}

func NewModelPrompt(user *user.User, task task.Task, chat *chat.Chat) *ModelPrompt {
	return &ModelPrompt{User: user, Task: task, Chat: chat}
}

type ModelResponse struct {
	Task            task.Task `json:"task"`
	IsComplete      bool      `json:"is_complete"`
	UserSummary     string    `json:"user_summary"`
	ResponseMessage string    `json:"response_message"`
}

func NewEmptyModelResponse() *ModelResponse {
	return &ModelResponse{}
}
