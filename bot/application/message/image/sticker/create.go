package sticker

import (
	"emobot/bot/application"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/mongo"
)

type createStickerCommand struct {
	collection *mongo.Collection
}

const createStickerCommandName = "Add Sticker"

var createStickerCommandDefinition = &discordgo.ApplicationCommand{
	Name: createStickerCommandName,
	Type: discordgo.MessageApplicationCommand,
}

func NewCreateStickerCommand(collection *mongo.Collection) application.Command {
	return &createStickerCommand{collection: collection}
}

func (c *createStickerCommand) Name() string {
	return createStickerCommandName
}

func (c *createStickerCommand) Definition() *discordgo.ApplicationCommand {
	return createStickerCommandDefinition
}

func (c *createStickerCommand) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {

}
