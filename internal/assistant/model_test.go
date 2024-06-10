package assistant

import (
	"aisu.ai/api/v2/internal/chat"
	"aisu.ai/api/v2/internal/user"
	openai "github.com/sashabaranov/go-openai"
	"os"
	"testing"
)

type MockModelExchangeRepository struct{}

func (r *MockModelExchangeRepository) Save(
	messageId string,
	systemMessage string,
	prompt string,
	response string,
	promptedAt int64,
	respondedAt int64,
) error {
	return nil
}

var assistant *Assistant

func setupTestCase(t *testing.T) {
	user := user.NewUser("Kai")
	task := NewGoalCreationTask()
	openaiClient := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	modelExchangeRepo := &MockModelExchangeRepository{}

	err := InitAssistants()
	if err != nil {
		t.Fatalf("Failed to initialize assistants: '%v'", err)
	}

	assistant, err = NewAssistant(user, task, openaiClient, modelExchangeRepo)
	if err != nil {
		t.Fatalf("Failed to create new assistant")
	}

	message := chat.NewMessage(chat.SenderUser, "I want to learn to swim")
	assistant.Chat.Append(message)
}

func TestPromptModelExecutesSuccessfully(t *testing.T) {
	setupTestCase(t)

	response, err := promptModel(assistant)
	if err != nil {
		t.Errorf("An error occurred while prompting the openai model: '%s'", err.Error())
	}

	if len(response.UserSummary) == 0 {
		t.Errorf("No user summary returned by openai model")
	}

	if len(response.ResponseMessage) == 0 {
		t.Errorf("No response message returned by openai model")
	}
}
