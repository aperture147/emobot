package eval

import (
	"emobot/bot/application"
	"github.com/bwmarrin/discordgo"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

type initEvalCommand struct {
	execEnv map[string]*interp.Interpreter
}

const initEvalCommandName = "init-eval"

var newEvalCommandDefinition = &discordgo.ApplicationCommand{
	Name:        initEvalCommandName,
	Description: "init a go evaluate environment",
	Type:        discordgo.ChatApplicationCommand,
}

func NewInitEvalCommand() application.Command {
	return &initEvalCommand{}
}

func (c *initEvalCommand) Name() string {
	return initEvalCommandName
}

func (c *initEvalCommand) Definition() *discordgo.ApplicationCommand {
	return newEvalCommandDefinition
}

func (c *initEvalCommand) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
