package user

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"log/slog"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(database *mongo.Database) *UserRepository {
	return &UserRepository{
		collection: database.Collection("users"),
	}
}

func (r *UserRepository) Save(user *User) (*User, error) {
	result, err := r.collection.InsertOne(context.TODO(), user)
	if err != nil {
		return nil, err
	}

	userID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		errMsg := "An error occurred while converting the user's ID to a mongo object ID"
		slog.Error(errMsg, "error", err)
		return nil, fmt.Errorf("%s: %w", errMsg, err)
	}
	user.Id = userID.Hex()

	return user, nil
}

func (r *UserRepository) Get(id string) (*User, error) {
	userId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Failed to convert hex user ID to mongo object ID")
		return nil, err
	}

	var user = &User{}
	err = r.collection.FindOne(
		context.TODO(),
		bson.M{"_id": userId},
	).Decode(user)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Printf("No user with id '%s' could be found\n", id)
		}
		log.Printf("An error occurred while retrieving user with id '%s'\n", id)
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) Update(user *User) error {
	userId, err := primitive.ObjectIDFromHex(user.Id)
	if err != nil {
		log.Printf("Failed to convert hex user ID '%s' to a mongo object ID", user.Id)
		return err
	}

	result, err := r.collection.UpdateByID(context.TODO(), userId, bson.M{
		"$set": bson.M{"summary": user.Summary, "goals": user.Goals, "schedule": user.Schedule},
	})
	if err != nil {
		slog.Error("An error occurred while updating the user", "error", err)
		return err
	}

	// Assert that only one user has the specified ID. Throw an error if more than one user
	// has the same ID to avoid cross contamination
	if result.MatchedCount != 1 {
		slog.Error("Failed to update user")
		return fmt.Errorf("Failed to update user with id '%s'", user.Id)
	}

	// if result.UpsertedCount != result.MatchedCount {
	// 	log.Printf("Failed to identify and/or update user with unique ID: '%s'", user.Id)
	// 	return errors.New(fmt.Sprintf("Unique user with ID '%s' could not be updated", user.Id))
	// }
	return nil
}
