package server

import (
	"emobot/bot/application"
	"go.mongodb.org/mongo-driver/mongo"
)

type commandCollection struct {
	client *mongo.Client
}

func NewCommandCollection(client *mongo.Client) application.CommandCollection {
	return commandCollection{client: client}
}

func (c commandCollection) GetAllCommands() (commands []application.Command) {
	pingCommand := NewPingCommand(c.client)

	commands = append(commands, pingCommand)
	return commands
}
