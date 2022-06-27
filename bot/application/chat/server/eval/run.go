package eval

import (
	"emobot/bot/application"
	"fmt"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/traefik/yaegi/interp"
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

func NewEvalCommand(execEnv map[string]*interp.Interpreter) application.Command {
	return &evalCommand{execEnv}
}

func (c *evalCommand) Name() string {
	return evalCommandName
}

func (c *evalCommand) Definition() *discordgo.ApplicationCommand {
	return evalCommandDefinition
}

func (c *evalCommand) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var exec *interp.Interpreter
	var ok bool
	var content string
	userId := i.Member.User.ID

	if exec, ok = c.execEnv[userId]; !ok {
		content = fmt.Sprintf("no existing environment found, create one by using `%s` command", initEvalCommandName)
		log.Infof("user %s evaluated code on nil environment", userId)
	} else {
		data := i.ApplicationCommandData()
		code := data.Options[0].StringValue()
		v, err := exec.Eval(code)
		if err != nil {
			log.Println("cannot evaluate code with reason:", err)
			content = "cannot evaluate sticker with reason: " + err.Error()
		} else {
			log.Printf("user %s evaluated code %s", i.Member.User.ID, code)
			if !v.IsValid() {
				content = "no return"
			} else {
				content = v.String()
			}
		}
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	})

	if err != nil {
		log.Println("cannot send message with reason:", err)
	}

}
