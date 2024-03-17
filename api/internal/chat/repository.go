package chat

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type ChatRepository struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewChatRepository() *ChatRepository {
	clientOptions := options.Client().ApplyURI("mongodb://host.docker.internal:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	collection := client.Database("aisu").Collection("chats")
	return &ChatRepository{client: client, collection: collection}
}

func (r *ChatRepository) Save(chat *Chat) string {
	result, err := r.collection.InsertOne(context.TODO(), chat)
	if err != nil {
		log.Fatal(err)
	}

	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		log.Fatal("InsertedID is not an ObjectID")
	}

	return insertedID.Hex()
}

func (r *ChatRepository) Get(id string) *Chat {
	chatObjectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Created hex")
	userChat := &Chat{}
	filter := bson.M{"_id": chatObjectId}
	err = r.collection.FindOne(context.TODO(), filter).Decode(userChat)
	if err != nil {
		log.Fatal(err)
	}

	return userChat
}

func (r *ChatRepository) SaveNewMessages(id string, messages ...*Message) {
	chatObjectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Fatal(err)
	}
	filter := bson.M{"_id": chatObjectId}

	update := bson.M{
		"$push": bson.M{
			"messages": bson.M{
				"$each": messages,
			},
		},
	}
	_, err = r.collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	return
}
