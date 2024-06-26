package assistant

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"aisu.ai/api/v2/internal/chat"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AssistantRepository struct {
	collection *mongo.Collection
}

func NewAssistantRepository(database *mongo.Database) *AssistantRepository {
	return &AssistantRepository{
		collection: database.Collection("assistants"),
	}
}

func (r *AssistantRepository) Save(assistant *Assistant) (*Assistant, error) {
	if r.collection == nil {
		slog.Warn("Collection not initialized correctly")
	}
	assistant.UpdatedAt = time.Now().UnixMilli()
	result, err := r.collection.InsertOne(context.TODO(), assistant)
	if err != nil {
		errMsg := "An error occurred while inserting the assistant into the database"
		slog.Error(errMsg, "error", err)
		return nil, fmt.Errorf("%s: %w", errMsg, err)
	}

	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		errMsg := "An error occurred while converting the assistant's ID to a mongo object ID"
		slog.Error(errMsg, "error", err)
		return nil, fmt.Errorf("%s: %w", errMsg, err)
	}

	assistant.Id = insertedID.Hex()
	return assistant, nil
}

func (r *AssistantRepository) Get(id string) (*Assistant, error) {
	assistantId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		errMsg := "Failed to convert hex value to mongo objectID"
		slog.Error(errMsg, "value", id, "error", err)
		return nil, fmt.Errorf("%s: %w", errMsg, err)
	}

	// Find the assistant with the specified id and return it and the
	// most recent chat messages up to the message limit.
	filter := bson.M{"_id": assistantId}
	options := options.FindOne().SetProjection(bson.M{
		"messages": bson.M{"$slice": -chat.DefaultMessageLimit},
	})

	assistant := &Assistant{}
	err = r.collection.FindOne(context.TODO(), filter, options).Decode(assistant)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			slog.Error("No assistant with the specified ID could be found", "id", id, "err", err)
			return nil, fmt.Errorf("No assistant with ID '%s' could be found: %w", id, err)
		}
		slog.Error("Failed to load the assistant with the specified ID", "id", id, "err", err)
		return nil, fmt.Errorf("Failed to load the assistant with ID '%s': %w", id, err)
	}

	return assistant, nil
}

func (r *AssistantRepository) Update(assistant *Assistant, numNewMessages int) (*Assistant, error) {
	id, err := primitive.ObjectIDFromHex(assistant.Id)
	if err != nil {
		errMsg := "Failed to convert hex value to mongo objectID"
		slog.Error(errMsg, "value", assistant.Id, "error", err)
		return nil, fmt.Errorf("%s: %w", errMsg, err)
	}

	documentUpdate := bson.M{}
	switch assistant.Task.Objective() {
	case ObjectiveGoalCreation:
		t, ok := assistant.Task.(*GoalCreationTask)
		if !ok {
			return nil, errors.New("Failed to convert task with objective 'goal_creation' to expected struct 'GoalCreationTask'")
		}
		documentUpdate["task"] = t
	case ObjectiveMilestoneCreation:
		t, ok := assistant.Task.(*MilestoneCreationTask)
		if !ok {
			return nil, errors.New("Failed to convert task with objective 'milestone_creation' to expected struct 'MilestoneCreationTask'")
		}
		documentUpdate["task"] = t
	case ObjectiveScheduleCreation:
		t, ok := assistant.Task.(*ScheduleCreationTask)
		if !ok {
			return nil, errors.New("Failed to convert task with objective 'schedule_creation' to expected struct 'ScheduleCreationTask'")
		}
		documentUpdate["task"] = t
	}

	// Identity the new messages to be saved.
	newMessages := []*chat.Message{}
	for _, message := range assistant.Chat.Messages {
		if message.CreatedAt > assistant.UpdatedAt {
			newMessages = append(newMessages, message)
		}
	}
	assistant.UpdatedAt = time.Now().UnixMilli()
	documentUpdate["updated_at"] = assistant.UpdatedAt

	// Save the assistant's current task and any new messages.
	update := bson.M{
		// Save the current the
		"$set": documentUpdate,
		"$push": bson.M{
			"chat.messages": bson.M{
				"$each": newMessages,
			},
		},
	}
	_, err = r.collection.UpdateByID(context.TODO(), id, update)
	if err != nil {
		// TODO: Create a custom 'Not Found' error for assistants.
		if errors.Is(err, mongo.ErrNoDocuments) {
			slog.Error("No assistant with the specified ID could be found", "id", assistant.Id, "err", err)
			return nil, fmt.Errorf("No assistant with ID '%s' could be found: %w", id, err)
		}
		slog.Error("Failed to load the assistant with the specified ID", "id", id, "err", err)
		return nil, fmt.Errorf("Failed to load the assistant with ID '%s': %w", id, err)
	}
	return assistant, nil
}
