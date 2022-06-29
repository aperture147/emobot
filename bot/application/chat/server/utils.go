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
	pingCmd := NewPingCommand(c.client)
	commands = append(commands, pingCmd)
	return commands
}
