package cmd

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

type StickerData struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type StickerSlashCommand struct {
	Database *mongo.Client
}

const StickerCommandName = "sticker"

var StickerCommandDefinition = &discordgo.ApplicationCommand{
	Name:        StickerCommandName,
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

func NewStickerSlashCommand(db *mongo.Client) SlashCommand {
	return &StickerSlashCommand{Database: db}
}

func (c *StickerSlashCommand) Name() string {
	return StickerCommandName
}

func (c *StickerSlashCommand) Definition() *discordgo.ApplicationCommand {
	return StickerCommandDefinition
}

func (c *StickerSlashCommand) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		//data := i.ApplicationCommandData()
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "test interaction",
			},
		})
		if err != nil {
			log.Println("sticker command handler error,", err)
		}
	case discordgo.InteractionApplicationCommandAutocomplete:
		data := i.ApplicationCommandData()
		findAttr := data.Options[0].StringValue()

		var choices []*discordgo.ApplicationCommandOptionChoice

		if findAttr != "" {
			collection := c.Database.Database("emobot").Collection("data")
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			res, err := collection.Find(ctx, bson.D{
				{"name", findAttr + ".*/"},
			})

			if err == nil {
				var stickers []StickerData
				err = res.All(ctx, stickers)

				for _, sticker := range stickers {
					choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
						Name:  sticker.Name,
						Value: sticker.Url,
					})
				}
			}

		}

		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionApplicationCommandAutocompleteResult,
			Data: &discordgo.InteractionResponseData{
				Choices: choices,
			},
		})

		if err != nil {
			log.Println("cannot autocomplete,", err)
		}
	}
}

type AddStickerSlashCommand struct {
	Database *mongo.Client
}

const AddStickerSlashCommandName = "add-sticker"

var AddStickerSlashDefinition = &discordgo.ApplicationCommand{
	Name:        AddStickerSlashCommandName,
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
			Description:  "image url of the sticker",
			Type:         discordgo.ApplicationCommandOptionString,
			Required:     true,
			Autocomplete: false,
		},
	},
}

func NewAddStickerSlashCommand(db *mongo.Client) SlashCommand {
	return &AddStickerSlashCommand{Database: db}
}

func (c *AddStickerSlashCommand) Name() string {
	return AddStickerSlashCommandName
}

func (c *AddStickerSlashCommand) Definition() *discordgo.ApplicationCommand {
	return AddStickerSlashDefinition
}

func (c *AddStickerSlashCommand) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()

	stickerName := data.Options[0].StringValue()
	stickerUrl := data.Options[1].StringValue()

	collection := c.Database.Database("emobot").Collection("data")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, StickerData{stickerName, stickerUrl})

	if err != nil {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "server error, cannot add sticker",
			},
		})
		log.Println("cannot add sticker,", err)
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "sticker added",
		},
	})
	if err != nil {
		log.Println("add sticker command handler error,", err)
	}
}
