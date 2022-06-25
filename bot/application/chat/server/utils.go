package server

import (
	"emobot/bot/application"
	"go.mongodb.org/mongo-driver/mongo"
)

type CommandCollection struct {
	client *mongo.Client
}

func NewCommandCollection(client *mongo.Client) CommandCollection {
	return CommandCollection{client: client}
}

func (c CommandCollection) GetAllCommands() (commands []application.Command) {
	pingCommand := NewPingCommand(c.client)

	commands = append(commands, pingCommand)
	return commands
}
