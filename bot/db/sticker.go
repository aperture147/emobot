package db

import (
	"context"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Sticker struct {
	ObjectId string `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string `json:"name" bson:"name"`
	Url      string `json:"url" bson:"url"`
}

func getContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 5*time.Second)
}

func GetSticker(collection *mongo.Collection, stickerId string) (*Sticker, error) {
	ctx, cancel := getContext()
	defer cancel()
	objectId, err := primitive.ObjectIDFromHex(stickerId)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return nil, nil
		}
		return nil, err
	}

	findOpt := options.Find()
	findOpt.SetProjection(bson.M{
		"name": false,
	})
	result := collection.FindOne(ctx, bson.M{"_id": objectId})
	var sticker Sticker
	err = result.Decode(&sticker)
	return &sticker, err
}

func CreateSticker(collection *mongo.Collection, name string, url string) (string, error) {
	ctx, cancel := getContext()
	defer cancel()
	result, err := collection.InsertOne(ctx, Sticker{Name: name, Url: url})
	stickerId := result.InsertedID.(primitive.ObjectID)
	return stickerId.Hex(), err
}

func DeleteSticker(collection *mongo.Collection, stickerId string) (*Sticker, error) {
	ctx, cancel := getContext()
	defer cancel()
	objectId, err := primitive.ObjectIDFromHex(stickerId)
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			return nil, nil
		}
		return nil, err
	}

	var sticker Sticker

	deleteResult := collection.FindOneAndDelete(ctx, bson.M{
		"_id": objectId,
	})
	err = deleteResult.Decode(&sticker)
	return &sticker, err
}

func GetStickerAutocompleteList(collection *mongo.Collection, findAttr string) ([]Sticker, error) {
	findOpts := options.Find()
	findOpts.SetLimit(25)
	findOpts.SetProjection(bson.M{
		"url": false,
	})

	ctx, cancel := getContext()
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{
		"name": bson.M{"$regex": primitive.Regex{Pattern: findAttr + ".*", Options: "i"}},
	}, findOpts)

	if err != nil {
		return nil, err
	}

	var stickers []Sticker

	err = cursor.All(ctx, &stickers)

	if err != nil {
		return nil, err
	}

	return stickers, nil
}
