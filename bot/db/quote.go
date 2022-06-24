package db

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Quote struct {
	ObjectId string `json:"_id,omitempty" bson:"_id,omitempty"`
	Title    string `json:"title" bson:"title"`
	Content  string `json:"content" bson:"content"`
}

func GetQuote(collection *mongo.Collection, stickerId string) (*Quote, error) {
	ctx, cancel := getContext()
	defer cancel()
	objectId := getObjectId(stickerId)

	if objectId == nil {
		return nil, nil
	}

	findOpt := options.Find()
	findOpt.SetProjection(bson.M{
		"name": false,
	})
	result := collection.FindOne(ctx, bson.M{"_id": objectId})
	var quote Quote
	err := result.Decode(&quote)
	return &quote, err
}

func GetQuoteAutocompleteList(collection *mongo.Collection, findAttr string) ([]Quote, error) {
	findOpts := options.Find()
	findOpts.SetLimit(25)
	findOpts.SetProjection(bson.M{
		"content": false,
	})

	ctx, cancel := getContext()
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{
		"title": bson.M{"$regex": primitive.Regex{Pattern: findAttr + ".*", Options: "i"}},
	}, findOpts)

	if err != nil {
		return nil, err
	}

	var quotes []Quote

	err = cursor.All(ctx, &quotes)

	if err != nil {
		return nil, err
	}

	return quotes, nil
}
