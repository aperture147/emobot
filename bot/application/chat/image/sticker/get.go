package sticker

import (
	"emobot/bot/application"
	"emobot/bot/db"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type getStickerCommand struct {
	collection *mongo.Collection
}

const getStickerCommandName = "sticker"

var getStickerCommandDefinition = &discordgo.ApplicationCommand{
	Name:        getStickerCommandName,
	Description: "get a sticker",
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

func NewGetStickerCommand(collection *mongo.Collection) application.Command {
	return &getStickerCommand{collection: collection}
}

func (c *getStickerCommand) Name() string {
	return getStickerCommandName
}

func (c *getStickerCommand) Definition() *discordgo.ApplicationCommand {
	return getStickerCommandDefinition
}

func (c *getStickerCommand) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		stickerId := data.Options[0].StringValue()
		sticker, err := db.GetSticker(c.collection, stickerId)

		var content string

		if err != nil {
			log.WithField("guild_id", i.GuildID).Println("cannot get sticker with reason:", err)
			content = "server error, cannot get sticker"
		} else if sticker == nil {
			log.WithField("guild_id", i.GuildID).Printf("user %s cannot get sticker %s", i.Member.User.ID, stickerId)
			content = "no sticker found"
		} else {
			log.WithField("guild_id", i.GuildID).Printf("user %s used sticker %s", i.Member.User.ID, stickerId)
			content = sticker.Url
		}

		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: content,
			},
		})
		if err != nil {
			log.WithField("guild_id", i.GuildID).Println("cannot send sticker with reason:", err)
		}

	case discordgo.InteractionApplicationCommandAutocomplete:
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
