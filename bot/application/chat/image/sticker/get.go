package sticker

import (
	"emobot/bot/application"
	"emobot/bot/db"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type GetStickerCommand struct {
	collection *mongo.Collection
}

const GetStickerCommandName = "sticker"

var GetStickerCommandDefinition = &discordgo.ApplicationCommand{
	Name:        GetStickerCommandName,
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

func NewGetStickerSlashCommand(collection *mongo.Collection) application.Command {
	return &GetStickerCommand{collection: collection}
}

func (c *GetStickerCommand) Name() string {
	return GetStickerCommandName
}

func (c *GetStickerCommand) Definition() *discordgo.ApplicationCommand {
	return GetStickerCommandDefinition
}

func (c *GetStickerCommand) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		data := i.ApplicationCommandData()
		stickerId := data.Options[0].StringValue()
		sticker, err := db.GetSticker(c.collection, stickerId)

		var content string

		if err != nil {
			log.Println("cannot get sticker with reason:", err)
			content = "server error, cannot get sticker"
		} else if sticker == nil {
			log.Printf("user %s cannot get sticker %s\n", i.Member.User.ID, stickerId)
			content = "no sticker found"
		} else {
			log.Printf("user %s used sticker %s\n", i.Member.User.ID, stickerId)
			content = sticker.Url
		}

		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: content,
			},
		})
		if err != nil {
			log.Println("cannot send sticker with reason:", err)
		}

	case discordgo.InteractionApplicationCommandAutocomplete:
		data := i.ApplicationCommandData()
		findAttr := data.Options[0].StringValue()

		var stickerChoices []*discordgo.ApplicationCommandOptionChoice
		var err error

		if findAttr != "" {
			stickerChoices, err = GetStickerAutocompleteChoice(c.collection, findAttr)

			if err != nil {
				log.Println("autocomplete error with reason:", err)
			}
		}

		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionApplicationCommandAutocompleteResult,
			Data: &discordgo.InteractionResponseData{
				Choices: stickerChoices,
			},
		})

		if err != nil {
			log.Println("cannot send autocomplete command with reason:", err)
		}
	}
}
