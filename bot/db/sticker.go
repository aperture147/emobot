package db

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Sticker struct {
	ObjectId string `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string `json:"name" bson:"name"`
	Url      string `json:"url" bson:"url"`
}

func GetSticker(collection *mongo.Collection, id string) (*Sticker, error) {
	opt := options.FindOne()
	opt.SetProjection(bson.M{
		"name": false,
	})

	var sticker Sticker
	exists, err := GetOne(collection, id, &sticker, opt)
	if !exists {
		return nil, nil
	}
	return &sticker, err
}

func CreateSticker(collection *mongo.Collection, name string, url string) (string, error) {
	return CreateOne(collection, Sticker{Name: name, Url: url})
}

func DeleteSticker(collection *mongo.Collection, id string) (*Sticker, error) {
	var sticker Sticker
	opt := options.FindOneAndDelete()
	opt.SetProjection(bson.M{
		"url": false,
	})

	exists, err := DeleteOne(collection, id, &sticker, opt)
	if !exists {
		return nil, nil
	}
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
