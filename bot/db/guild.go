package db

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type GuildInfo struct {
	GuildId string `json:"guild_id" bson:"guild_id"`
}

func GetGuildIdList(collection *mongo.Collection) ([]string, error) {
	ctx, cancel := getContext()
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

func GetGuildDatabase(guildId string, client *mongo.Client) *mongo.Database {
	guildDbName := "guild-" + guildId
	return client.Database(guildDbName)
}
