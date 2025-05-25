package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserNotification struct {
	UserID    int       `bson:"user_id"`
	Name      string    `bson:"name"`
	Email     string    `bson:"email"`
	UpdatedAt time.Time `bson:"updated_at"`
}

type MongoDB struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoDB(uri, dbName, collectionName string) (*MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	collection := client.Database(dbName).Collection(collectionName)
	return &MongoDB{
		client:     client,
		collection: collection,
	}, nil
}

func (m *MongoDB) UpsertUser(notification *UserNotification) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"user_id": notification.UserID}
	update := bson.M{"$set": notification}
	opts := options.Update().SetUpsert(true)

	_, err := m.collection.UpdateOne(ctx, filter, update, opts)
	return err
}

func (m *MongoDB) GetUser(userID int) (*UserNotification, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	notification := &UserNotification{}
	err := m.collection.FindOne(ctx, bson.M{"user_id": userID}).Decode(notification)
	if err != nil {
		return nil, err
	}

	return notification, nil
}

func (m *MongoDB) GetAllNotifications() ([]*UserNotification, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := m.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var notifications []*UserNotification
	if err = cursor.All(ctx, &notifications); err != nil {
		return nil, err
	}

	return notifications, nil
}
