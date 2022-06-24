package sticker

import (
	"emobot/bot/cmd"
	"emobot/bot/db"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type CreateStickerSlashCommand struct {
	collection *mongo.Collection
}

const CreateStickerCommandName = "add-sticker"

var CreateStickerCommandDefinition = &discordgo.ApplicationCommand{
	Name:        CreateStickerCommandName,
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

func NewCreateStickerSlashCommand(collection *mongo.Collection) cmd.SlashCommand {
	return &CreateStickerSlashCommand{collection: collection}
}

func (c *CreateStickerSlashCommand) Name() string {
	return CreateStickerCommandName
}

func (c *CreateStickerSlashCommand) Definition() *discordgo.ApplicationCommand {
	return CreateStickerCommandDefinition
}

func (c *CreateStickerSlashCommand) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()

	stickerName := data.Options[0].StringValue()
	stickerUrl := data.Options[1].StringValue()

	stickerId, err := db.CreateSticker(c.collection, stickerName, stickerUrl)

	var content string

	if err != nil {
		content = "server error, cannot add sticker"
		log.Println("cannot add sticker with reason:", err)
	} else {
		log.Printf("user %s created sticker %q with ID %s\n", i.Member.User.ID, stickerName, stickerId)
		content = "sticker `" + stickerName + "` added"
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	})

	if err != nil {
		log.Println("cannot send create message with reason:", err)
	}
}
