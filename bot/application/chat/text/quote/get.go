package quote

import (
	"emobot/bot/application"
	"emobot/bot/db"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type GetQuoteCommand struct {
	collection *mongo.Collection
}

const getQuoteCommandName = "quote"

var getQuoteCommandDefinition = &discordgo.ApplicationCommand{
	Name:        getQuoteCommandName,
	Description: "get quote",
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

func NewGetQuoteCommand(collection *mongo.Collection) application.Command {
	return &GetQuoteCommand{collection}
}

func (c *GetQuoteCommand) Name() string {
	return getQuoteCommandName
}

func (c *GetQuoteCommand) Definition() *discordgo.ApplicationCommand {
	return getQuoteCommandDefinition
}

func (c *GetQuoteCommand) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		data := i.ApplicationCommandData()
		quoteId := data.Options[0].StringValue()
		quote, err := db.GetQuote(c.collection, quoteId)

		var content string

		if err != nil {
			log.Println("cannot get quote with reason:", err)
			content = "server error, cannot get quote"
		} else if quote == nil {
			log.Printf("user %s cannot get quote %s\n", i.Member.User.ID, quoteId)
			content = "no quote found"
		} else {
			log.Printf("user %s used quote %s\n", i.Member.User.ID, quoteId)
			content = quote.Content
		}

		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: content,
			},
		})
		if err != nil {
			log.Println("cannot send quote with reason:", err)
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
