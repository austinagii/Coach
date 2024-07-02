package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"

	"aisu.ai/api/v2/cmd/server/shared/middleware"
	"aisu.ai/api/v2/internal/assistant"
	"aisu.ai/api/v2/internal/user"

	"github.com/gin-gonic/gin"
	openai "github.com/sashabaranov/go-openai"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var openaiClient *openai.Client
var database *mongo.Database

var modelExchangeRepository *assistant.LanguageModelExchangeRepository
var assistantRepository *assistant.AssistantRepository
var userRepository *user.UserRepository

func main() {
	router := gin.Default()

	router.Use(middleware.CorsMiddleware())

	router.POST("/users", CreateUser)
	router.GET("/users/:user_id", nil)

	router.POST("/chats", CreateChat)
	router.POST("/chats/:id/messages", HandleUserMessage)

	router.Run("0.0.0.0:8080")
}

func init() {
	slog.Info("Initializing the application...")
	var err error

	openaiClient, err = initializeOpenaiClient()
	if err != nil {
		slog.Error("An error occurred while initializing the openai client", "error", err)
		os.Exit(1)
	}

	database, err = initializeDatabaseConnection()
	if err != nil {
		slog.Error("An error occurred while initializing the mongo database connection", "error", err)
		os.Exit(1)
	}

	if err := assistant.InitAssistants(); err != nil {
		log.Fatalf("An error occurred while initializing static data: %v", err)
	}

	modelExchangeRepository = assistant.NewLanguageModelExchangeRepository(database)
	assistantRepository = assistant.NewAssistantRepository(database)
	userRepository = user.NewUserRepository(database)
}

func initializeDatabaseConnection() (*mongo.Database, error) {
	clientOptions := options.Client().ApplyURI("mongodb://host.docker.internal:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %v", err)
	}
	log.Println("Successfully connected to MongoDB")
	return client.Database("aisu"), nil
}

func initializeOpenaiClient() (*openai.Client, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("No environment variable defined for OpenAI API key")
	}
	return openai.NewClient(apiKey), nil
}
