package main

import (
	"aisu.ai/api/v2/internal/assistant"
	"aisu.ai/api/v2/internal/chat"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type NewChatRequest struct {
	Task chat.TaskType `json:"task"`
}

type NewChatResponse struct {
	Id   string `json:"id"`
	Text string `json:"text"`
}

func CreateChat(context *gin.Context) {
	request := NewChatRequest{}
	if err := context.BindJSON(&request); err != nil {
		log.Panicf("Could not de-serialize JSON request body to 'NewChatRequest' struct: '%v'", err)
	}

	userChat, userTask := chat.NewChat(), chat.NewTask(request.Task, nil)
	initialChatMessage, err := userTask.GetInitialMessage()
	if err != nil {
		log.Fatal(err)
	}
	userChat.Append(initialChatMessage)

	chatRepository := chat.NewChatRepository()
	chatId := chatRepository.Save(userChat)
	userChat.Id = chatId

	response := &NewChatResponse{Id: chatId, Text: initialChatMessage.Text}
	context.IndentedJSON(http.StatusCreated, response)
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
