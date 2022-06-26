package quote

import (
	"emobot/bot/application"
	"emobot/bot/db"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type deleteQuoteCommand struct {
	collection *mongo.Collection
}

const deleteQuoteCommandName = "delete-quote"

var deleteQuoteCommandDefinition = &discordgo.ApplicationCommand{
	Name:        deleteQuoteCommandName,
	Description: "delete quote",
	Type:        discordgo.ChatApplicationCommand,
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:         "title",
			Description:  "title of the quote",
			Type:         discordgo.ApplicationCommandOptionString,
			Required:     true,
			Autocomplete: true,
		},
	},
}

func NewDeleteGetQuoteCommand(collection *mongo.Collection) application.Command {
	return &deleteQuoteCommand{collection}
}

func (c *deleteQuoteCommand) Name() string {
	return deleteQuoteCommandName
}

func (c *deleteQuoteCommand) Definition() *discordgo.ApplicationCommand {
	return deleteQuoteCommandDefinition
}

func (c *deleteQuoteCommand) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		data := i.ApplicationCommandData()
		quoteId := data.Options[0].StringValue()

		quote, err := db.DeleteQuote(c.collection, quoteId)

		var content string

		if err != nil {
			content = "server error, cannot delete quote"
			log.Println("cannot delete quote with reason:", err)
		} else {
			content = "quote `" + quote.Title + "` deleted"
			log.Printf("user %s deleted quote %s", i.Member.User.ID, quoteId)
		}

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

		var quoteChoices []*discordgo.ApplicationCommandOptionChoice
		var err error

		if findAttr != "" {
			quoteChoices, err = GetQuoteAutocompleteChoice(c.collection, findAttr)

			if err != nil {
				log.Println("autocomplete error with reason:", err)
			}
		}

		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionApplicationCommandAutocompleteResult,
			Data: &discordgo.InteractionResponseData{
				Choices: quoteChoices,
			},
		})

		if err != nil {
			log.Println("cannot send autocomplete command with reason:", err)
		}
	}
}
