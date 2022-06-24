package sticker

import (
	"emobot/bot/cmd"
	"emobot/bot/db"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewAllStickerCommands(stickerCollection *mongo.Collection) []cmd.SlashCommand {
	getStickerCommand := NewGetStickerSlashCommand(stickerCollection)
	addStickerCommand := NewCreateStickerSlashCommand(stickerCollection)
	deleteStickerCommand := NewDeleteStickerSlashCommand(stickerCollection)
	return []cmd.SlashCommand{getStickerCommand, addStickerCommand, deleteStickerCommand}
}

func GetStickerAutocompleteChoice(collection *mongo.Collection, findAttr string) ([]*discordgo.ApplicationCommandOptionChoice, error) {
	stickers, err := db.GetStickerAutocompleteList(collection, findAttr)
	if err != nil {
		return nil, err
	}

	var choices []*discordgo.ApplicationCommandOptionChoice
	for _, sticker := range stickers {
		choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
			Name:  sticker.Name,
			Value: sticker.ObjectId,
		})
	}

	return choices, nil
}
