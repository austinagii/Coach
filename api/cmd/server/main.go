package main

import (
	// "aisu.ai/api/v2/cmd/server/controllers"
	"aisu.ai/api/v2/cmd/server/shared/middleware"
	// "aisu.ai/api/v2/internal/assistant"
	"aisu.ai/api/v2/internal/chat"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

func main() {
	if err := initializeStaticData(); err != nil {
		log.Panicf("The following error occurred while initializing the application's static data: %v", err)
	}

	router := gin.Default()
	router.Use(middleware.CorsMiddleware())
	router.POST("/chats", CreateChat)
	router.POST("/chats/:id/messages", HandleUserMessage)
	router.Run("0.0.0.0:8080")
}

func initializeStaticData() error {
	if err := chat.LoadTasks(); err != nil {
		return err
	}
	return nil
}

func initializeDatabaseConnection() (*mongo.Client, error) {
	// Do stuff here.
}
