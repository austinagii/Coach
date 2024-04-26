package main

import (
	"aisu.ai/api/v2/cmd/server/shared/middleware"
	"aisu.ai/api/v2/internal/assistant"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	openai "github.com/sashabaranov/go-openai"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

var openaiClient *openai.Client = openai.NewClient(os.Getenv("OPENAI_API_KEY"))
var dbClient *mongo.Client

func main() {
	router := gin.Default()

	router.Use(middleware.CorsMiddleware())

	router.POST("/users", CreateUser)
	router.POST("/users/:user_id/goals", nil)
	router.POST("/users/:user_id/goals/:goal_id/milestones", nil)
	router.POST("/chats", CreateChat)
	router.POST("/chats/:id/messages", HandleUserMessage)

	router.Run("0.0.0.0:8080")
}

func init() {
	if err := initializeStaticData(); err != nil {
		log.Fatalf("An error occurred while initializing static data: %v", err)
	}

	client, err := initializeDatabaseConnection()
	if err != nil {
		log.Fatalf("An error occurred while initializing the database connection: %v", err)
	}
	dbClient = client
}

func initializeStaticData() error {
	if err := LoadObjectiveDescriptions(); err != nil {
		return err
	}
	if err := LoadTaskInitialMessages(); err != nil {
		return err
	}
	if err := assistant.LoadTasks(); err != nil {
		return err
	}
	return nil
}

func initializeDatabaseConnection() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI("mongodb://host.docker.internal:27017")
	log.Println("Connecting to MongoDB")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %v", err)
	}
	log.Println("Successfully connected to MongoDB")
	return client, nil
}
