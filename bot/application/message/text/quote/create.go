package quote

import (
	"emobot/bot/application"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/mongo"
)

type createQuoteCommand struct {
	collection *mongo.Collection
}

const createQuoteCommandName = "Add Quote"

var createQuoteCommandDefinition = &discordgo.ApplicationCommand{
	Name: createQuoteCommandName,
	Type: discordgo.MessageApplicationCommand,
}

func NewCreateQuoteCommand(collection *mongo.Collection) application.Command {
	return &createQuoteCommand{collection}
}

func (c *createQuoteCommand) Name() string {
	return createQuoteCommandName
}

func (c *createQuoteCommand) Definition() *discordgo.ApplicationCommand {
	return createQuoteCommandDefinition
}

func (c *createQuoteCommand) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {

}
