package sticker

import (
	"emobot/bot/application"
	"emobot/bot/db"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type deleteStickerCommand struct {
	collection *mongo.Collection
}

const deleteStickerCommandName = "delete-sticker"

var deleteStickerCommandDefinition = &discordgo.ApplicationCommand{
	Name:        deleteStickerCommandName,
	Description: "delete sticker",
	Type:        discordgo.ChatApplicationCommand,
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:         "name",
			Description:  "name of the sticker",
			Type:         discordgo.ApplicationCommandOptionString,
			Required:     true,
			Autocomplete: true,
		},
	},
}

func NewDeleteStickerCommand(collection *mongo.Collection) application.Command {
	return &deleteStickerCommand{collection: collection}
}

func (c *deleteStickerCommand) Name() string {
	return deleteStickerCommandName
}

func (c *deleteStickerCommand) Definition() *discordgo.ApplicationCommand {
	return deleteStickerCommandDefinition
}

func (c *deleteStickerCommand) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		data := i.ApplicationCommandData()
		stickerId := data.Options[0].StringValue()

		sticker, err := db.DeleteSticker(c.collection, stickerId)

		var content string

		if err != nil {
			content = "server error, cannot delete sticker"
			log.WithField("guild_id", i.GuildID).Println("cannot delete sticker with reason:", err)
		} else {
			content = "sticker `" + sticker.Name + "` deleted"
			log.WithField("guild_id", i.GuildID).Printf("user %s deleted sticker %s", i.Member.User.ID, stickerId)
		}

		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: content,
			},
		})

		if err != nil {
			log.WithField("guild_id", i.GuildID).Println("cannot send delete message with reason:", err)
		}

	case discordgo.InteractionApplicationCommandAutocomplete:
		data := i.ApplicationCommandData()
		findAttr := data.Options[0].StringValue()

		var stickerChoices []*discordgo.ApplicationCommandOptionChoice
		var err error

		if findAttr != "" {
			stickerChoices, err = GetStickerAutocompleteChoice(c.collection, findAttr)

			if err != nil {
				log.WithField("guild_id", i.GuildID).Println("autocomplete error with reason:", err)
			}
		}

		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionApplicationCommandAutocompleteResult,
			Data: &discordgo.InteractionResponseData{
				Choices: stickerChoices,
			},
		})

		if err != nil {
			log.WithField("guild_id", i.GuildID).Println("cannot send autocomplete command with reason:", err)
		}
	}
}
