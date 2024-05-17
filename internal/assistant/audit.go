package assistant

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"log/slog"
)

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
