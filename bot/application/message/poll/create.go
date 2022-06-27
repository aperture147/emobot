package poll

import (
	"emobot/bot/application"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/mongo"
)

type createPollCommand struct {
	collection *mongo.Collection
}

const createQuoteCommandName = "Make Poll"

var createPollCommandDefinition = &discordgo.ApplicationCommand{
	Name: createQuoteCommandName,
	Type: discordgo.MessageApplicationCommand,
}

func NewCreatePollCommand(collection *mongo.Collection) application.Command {
	return &createPollCommand{collection}
}

func (c *createPollCommand) Name() string {
	return createQuoteCommandName
}

func (c *createPollCommand) Definition() *discordgo.ApplicationCommand {
	return createPollCommandDefinition
}

func (c *createPollCommand) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {

}
