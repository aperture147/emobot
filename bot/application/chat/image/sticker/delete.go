package sticker

import (
	"emobot/bot/application"
	"emobot/bot/db"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type DeleteStickerCommand struct {
	collection *mongo.Collection
}

const DeleteStickerCommandName = "delete-sticker"

var DeleteStickerCommandDefinition = &discordgo.ApplicationCommand{
	Name:        DeleteStickerCommandName,
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

func NewDeleteStickerSlashCommand(collection *mongo.Collection) application.Command {
	return &DeleteStickerCommand{collection: collection}
}

func (c *DeleteStickerCommand) Name() string {
	return DeleteStickerCommandName
}

func (c *DeleteStickerCommand) Definition() *discordgo.ApplicationCommand {
	return DeleteStickerCommandDefinition
}

func (c *DeleteStickerCommand) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		data := i.ApplicationCommandData()
		stickerId := data.Options[0].StringValue()

		sticker, err := db.DeleteSticker(c.collection, stickerId)

		content := "sticker `" + sticker.Name + "` deleted"
		if err != nil {
			content = "server error, cannot delete sticker"
			log.Println("cannot delete sticker with reason:", err)
		}
		log.Printf("user %s deleted sticker %s", i.Member.User.ID, stickerId)

		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: content,
			},
		})

		if err != nil {
			log.Println("cannot send delete message with reason:", err)
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
