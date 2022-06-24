package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func getContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 5*time.Second)
}

func getObjectId(id string) *primitive.ObjectID {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil
	}
	return &objectId
}
