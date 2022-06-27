package eval

import (
	"emobot/bot/application"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

type evalCommand struct {
	execEnv map[string]*interp.Interpreter
}

const evalCommandName = "eval"

var evalCommandDefinition = &discordgo.ApplicationCommand{
	Name:        evalCommandName,
	Description: "evaluate a go command",
	Type:        discordgo.ChatApplicationCommand,
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:         "code",
			Description:  "code you want to evaluate",
			Type:         discordgo.ApplicationCommandOptionString,
			Required:     true,
			Autocomplete: false,
		},
	},
}

func NewEvalCommand() application.Command {
	return &evalCommand{}
}

func (c *evalCommand) Name() string {
	return evalCommandName
}

func (c *evalCommand) Definition() *discordgo.ApplicationCommand {
	return evalCommandDefinition
}

func (c *evalCommand) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	exec := interp.New(interp.Options{})
	_ = exec.Use(stdlib.Symbols)

	data := i.ApplicationCommandData()
	code := data.Options[0].StringValue()
	v, err := exec.Eval(code)

	var content string

	if err != nil {
		log.Println("cannot evaluate code with reason:", err)
		content = "cannot evaluate sticker"
	} else {
		log.Printf("user %s evaluated code %s", i.Member.User.ID, code)
		content = v.String()
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	})

	if err != nil {
		log.Println("cannot send message with reason:", err)
	}

}
