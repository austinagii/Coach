package main

import (
	"aisu.ai/api/v2/cmd/server/shared/api"
	"aisu.ai/api/v2/internal/assistant"
	"aisu.ai/api/v2/internal/assistant/task"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

type NewChatRequest struct {
	Task task.Task `json:"task"`
}

type NewChatResponse struct {
	Id   string `json:"id"`
	Text string `json:"text"`
}

func CreateChat(context *gin.Context) {
	request := NewChatRequest{}
	if err := context.BindJSON(&request); err != nil {
		// TODO: Add user friendly error messages for request validation.
		slog.Error("Failed to deserialize request body to new chat request", "err", err)
		context.IndentedJSON(
			http.StatusBadRequest,
			api.NewApiError(api.InvalidRequest, err.Error()),
		)
		return
	}

	// Create a new assistant with it's own chat and save it for future use.
	assistant := assistant.NewAssistant(openaiClient, request.Task)
	repository := assistant.NewAssistantRepository(dbClient.Database("assistant"))
	if err := repository.Save(assistant); err != nil {

	}

	response := &NewChatResponse{
		Id:   assistant.Id,
		Text: assistant.Chat.GetLastMessage().Text,
	}
	context.IndentedJSON(http.StatusCreated, response)
	context.BindHeader(map[string]string{"Location": "https://api.superu.ai/v1/chat/%s"})
	return
}

func HandleUserMessage(context *gin.Context) {
	log.Print("Here 1")
	userMessage := chat.NewEmptyUserMessage()
	if err := context.BindJSON(userMessage); err != nil {
		log.Panicf("Could not de-serialize JSON request body to 'ChatRequest' struct: '%v'", err)
	}
	log.Print("Here 2")

	chatId := context.Param("id")
	log.Printf("task: %v, chatId: %s", userMessage.Task, chatId)
	chatAssistant := assistant.LoadAssistant(userMessage.Task, chatId)
	log.Print("Here 3")
	assistantMessage, err := chatAssistant.Respond(userMessage)
	if err != nil {
		log.Panicf("An error occurred while the assistant was responding: %v", err)
	}

	context.IndentedJSON(http.StatusOK, assistantMessage)
	return
}
