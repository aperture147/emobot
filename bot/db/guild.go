package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type GuildInfo struct {
	GuildId string `json:"guild_id" bson:"guild_id"`
}

func GetGuildIdList(collection *mongo.Collection) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cursor, err := collection.Find(ctx, bson.M{})

	if err != nil {
		return nil, err
	}

	var guildInfoList []GuildInfo

	err = cursor.All(ctx, &guildInfoList)

	if err != nil {
		return nil, err
	}

	var guildIdList []string

	for _, guildInfo := range guildInfoList {
		guildIdList = append(guildIdList, guildInfo.GuildId)
	}

	return guildIdList, nil
}

func GetGuildDatabase(client *mongo.Client, guildId string) *mongo.Database {
	guildDbName := "guild-" + guildId
	return client.Database(guildDbName)
}