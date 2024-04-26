package main

import (
	"aisu.ai/api/v2/internal/assistant"
	"aisu.ai/api/v2/internal/assistant/task"
	"fmt"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

type AssistantRequest struct {
	InitialTask    task.
	TargetEntityId string
}

func CreateNewAssistant(context *gin.Context) {
	request := &AssistantRequest{}
	if err := context.BindJSON(request); err != nil {
		slog.Error("Failed to parse request body to new assistant request")
		errResponse := NewApiError("invalid_request", "Failed to parse request body")
		context.IndentedJSON(http.StatusBadRequest, errResponse)
		return
	}

	task := assistant.NewTask(request.InitialTask, request.TargetEntityId)
	// assistant := assistant.NewAssistant(task)
	// assistantRepo.Save(assistant)
	context.IndentedJSON(http.StatusCreated, task)
}

func Chat(context *gin.Context) {
	context.Params.Get("assistant_id")
}
