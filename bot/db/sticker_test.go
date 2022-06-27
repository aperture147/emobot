package db

import (
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
	"testing"
)

func getTestStickerMongoArtifacts() (*mongo.Client, *mongo.Database, *mongo.Collection, error) {
	client, err := NewMongoClient()
	if err != nil {
		return nil, nil, nil, err
	}
	db := GetGuildDatabase("933206081570762763", client)
	collection := db.Collection("sticker")
	return client, db, collection, err
}

func TestGetStickerAutocompleteList(t *testing.T) {
	_, _, collection, err := getTestStickerMongoArtifacts()
	if err != nil {
		t.Fatal(err)
	}

	// test for existing autocomplete
	findAttr := "ro"

	stickers, err := GetStickerAutocompleteList(collection, findAttr)

	if err != nil {
		t.Fatal(err)
	}

	if len(stickers) == 0 {
		t.Fatal("empty sticker slice returned")
	}

	for _, sticker := range stickers {
		if !strings.HasPrefix(sticker.Name, findAttr) {
			t.Fatalf("sticker %s doesn't starts with %s\n", sticker.Name, findAttr)
		}
	}

	// test for not existing autocomplete
	findAttr = "abcxyz"

	stickers, err = GetStickerAutocompleteList(collection, findAttr)

	if err != nil {
		t.Fatal(err)
	}

	if len(stickers) != 0 {
		t.Fatal("sticker slice is supposed to be empty")
	}
}

func TestGetSticker(t *testing.T) {
	_, _, collection, err := getTestStickerMongoArtifacts()
	if err != nil {
		t.Fatal(err)
	}

	// test for existing sticker
	sticker, err := GetSticker(collection, "test")
	if err != nil {
		t.Fatal(err)
	}
	if sticker == nil {
		t.Fatal("sticker is supposed to be existed")
	}

	// test for not existing sticker
	sticker, err = GetSticker(collection, "test2")
	if err != nil {
		t.Fatal(err)
	}
	if sticker != nil {
		t.Fatal("sticker is supposed to be null")
	}

	t.Logf("ID: %s - URL: %s", sticker.ObjectId, sticker.Url)
}
