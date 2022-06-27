package server

import (
	"emobot/bot/application"
	"emobot/bot/application/chat/server/eval"
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
	evalCommand := eval.NewEvalCommand()
	commands = append(commands, pingCommand, evalCommand)
	return commands
}
