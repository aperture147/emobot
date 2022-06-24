package sticker

import (
	"emobot/bot/cmd"
	"emobot/bot/db"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type DeleteStickerSlashCommand struct {
	Collection *mongo.Collection
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

func NewDeleteStickerSlashCommand(collection *mongo.Collection) cmd.SlashCommand {
	return &DeleteStickerSlashCommand{Collection: collection}
}

func (c *DeleteStickerSlashCommand) Name() string {
	return DeleteStickerCommandName
}

func (c *DeleteStickerSlashCommand) Definition() *discordgo.ApplicationCommand {
	return DeleteStickerCommandDefinition
}

func (c *DeleteStickerSlashCommand) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		data := i.ApplicationCommandData()
		stickerName := data.Options[0].StringValue()

		err := db.DeleteSticker(c.Collection, stickerName)

		content := "sticker `" + stickerName + "` deleted"
		if err != nil {
			content = "server error, cannot delete sticker"
			log.Println("cannot delete sticker with reason: ", err)
		}
		log.Printf("user %s deleted sticker %s", i.Member.User.ID, stickerName)

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
			stickerChoices, err = GetStickerAutocompleteChoice(c.Collection, findAttr)

			if err != nil {
				log.Println("autocomplete error,", err)
			}
		}

		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionApplicationCommandAutocompleteResult,
			Data: &discordgo.InteractionResponseData{
				Choices: stickerChoices,
			},
		})

		if err != nil {
			log.Println("cannot send autocomplete command with reason: ", err)
		}
	}
}
