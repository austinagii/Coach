package user

import (
	// "go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepository struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewUserRepository(client *mongo.Client) *UserRepository {
  clientOptions := 
}
