package db

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetOne(collection *mongo.Collection, id string, container interface{}, options ...*options.FindOneOptions) (bool, error) {
	ctx, cancel := getContext()
	defer cancel()
	objectId := getObjectId(id)
	if objectId == nil {
		return false, nil
	}

	result := collection.FindOne(ctx, bson.M{"_id": objectId}, options...)
	err := result.Decode(container)
	if errors.Is(err, bson.ErrNilRegistry) {
		return false, nil
	}
	return true, err
}

func CreateOne(collection *mongo.Collection, object interface{}, options ...*options.InsertOneOptions) (string, error) {
	ctx, cancel := getContext()
	defer cancel()
	result, err := collection.InsertOne(ctx, object, options...)
	if isDuplicated(err) {
		return "", nil
	}
	stickerId := result.InsertedID.(primitive.ObjectID)
	return stickerId.Hex(), err
}

func UpdateOne(collection *mongo.Collection, stickerId string, update interface{}, options ...*options.UpdateOptions) (bool, error) {
	ctx, cancel := getContext()
	defer cancel()
	objectId := getObjectId(stickerId)
	if objectId == nil {
		return false, nil
	}

	_, err := collection.UpdateOne(ctx, bson.M{"_id": objectId}, update, options...)
	if err != nil {
		return false, err
	}
	return true, nil
}

func DeleteOne(collection *mongo.Collection, stickerId string, container interface{}, options ...*options.FindOneAndDeleteOptions) (bool, error) {
	ctx, cancel := getContext()
	defer cancel()
	objectId := getObjectId(stickerId)
	if objectId == nil {
		return false, nil
	}

	deleteResult := collection.FindOneAndDelete(ctx, bson.M{"_id": objectId}, options...)
	err := deleteResult.Decode(container)
	if errors.Is(err, bson.ErrNilRegistry) {
		return false, nil
	}
	return true, err
}
