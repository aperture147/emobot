package sticker

import (
	"emobot/bot/application"
	"emobot/bot/db"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type createStickerCommand struct {
	collection *mongo.Collection
}

const createStickerCommandName = "add-sticker"

var createStickerCommandDefinition = &discordgo.ApplicationCommand{
	Name:        createStickerCommandName,
	Description: "add a sticker",
	Type:        discordgo.ChatApplicationCommand,
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:         "name",
			Description:  "name of the sticker",
			Type:         discordgo.ApplicationCommandOptionString,
			Required:     true,
			Autocomplete: false,
		},
		{
			Name:         "url",
			Description:  "image URL of the sticker",
			Type:         discordgo.ApplicationCommandOptionString,
			Required:     true,
			Autocomplete: false,
		},
	},
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
	data := i.ApplicationCommandData()

	stickerName := data.Options[0].StringValue()
	stickerUrl := data.Options[1].StringValue()

	stickerId, err := db.CreateSticker(c.collection, stickerName, stickerUrl)

	var content string

	if err != nil {
		content = "server error, cannot add sticker"
		log.WithField("guild_id", i.GuildID).Println("cannot add sticker with reason:", err)
	} else if stickerId == "" {
		log.WithField("guild_id", i.GuildID).Printf("user %s failed to create a duplicated sticker %q", i.Member.User.ID, stickerName)
		content = "sticker `" + stickerName + "` already exists"
	} else {
		log.WithField("guild_id", i.GuildID).Printf("user %s created sticker %q with ID %s", i.Member.User.ID, stickerName, stickerId)
		content = "sticker `" + stickerName + "` added"
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	})

	if err != nil {
		log.WithField("guild_id", i.GuildID).Println("cannot send create message with reason:", err)
	}
}
