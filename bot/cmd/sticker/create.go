package sticker

import (
	"emobot/bot/cmd"
	"emobot/bot/db"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type CreateStickerSlashCommand struct {
	Collection *mongo.Collection
}

const CreateStickerSlashCommandName = "add-sticker"

var CreateStickerSlashDefinition = &discordgo.ApplicationCommand{
	Name:        CreateStickerSlashCommandName,
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
	return &CreateStickerSlashCommand{Collection: collection}
}

func (c *CreateStickerSlashCommand) Name() string {
	return CreateStickerSlashCommandName
}

func (c *CreateStickerSlashCommand) Definition() *discordgo.ApplicationCommand {
	return CreateStickerSlashDefinition
}

func (c *CreateStickerSlashCommand) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()

	stickerName := data.Options[0].StringValue()
	stickerUrl := data.Options[1].StringValue()

	stickerId, err := db.CreateSticker(c.Collection, stickerName, stickerUrl)

	var content string

	if err != nil {
		content = "server error, cannot add sticker"
		log.Println("cannot add sticker,", err)
	} else {
		log.Printf("user %s created sticker %s", i.Member.User.ID, stickerId)
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
