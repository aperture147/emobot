package eval

import (
	"emobot/bot/application"
	"fmt"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/traefik/yaegi/interp"
)

type initEvalCommand struct {
	execEnv map[string]*interp.Interpreter
}

const initEvalCommandName = "init-eval"

var newEvalCommandDefinition = &discordgo.ApplicationCommand{
	Name:        initEvalCommandName,
	Description: "init a go evaluate environment",
	Type:        discordgo.ChatApplicationCommand,
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:         "destroy",
			Description:  "automatically destroy old environment",
			Type:         discordgo.ApplicationCommandOptionBoolean,
			Required:     false,
			Autocomplete: false,
		},
	},
}

func NewInitEvalCommand(execEnv map[string]*interp.Interpreter) application.Command {
	return &initEvalCommand{execEnv}
}

func (c *initEvalCommand) Name() string {
	return initEvalCommandName
}

func (c *initEvalCommand) Definition() *discordgo.ApplicationCommand {
	return newEvalCommandDefinition
}

func (c *initEvalCommand) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	allowDestroy := false
	if len(options) != 0 {
		allowDestroy = options[0].BoolValue()
	}

	var content string
	userId := i.Member.User.ID
	if _, ok := c.execEnv[i.Member.User.ID]; ok {
		if allowDestroy {
			content = fmt.Sprintf("environment already exists, use `/destroy-eval` to destroy it or set `destroy = true` to overwrite")
			log.Infof("user %s's init request canceled due to existing environment", userId)
		} else {
			c.execEnv[i.Member.User.ID] = newEvalEnvironment()
			content = "old environment was destroyed and replaced by a new one"
			log.Infof("user %s replaced old environment by a new one", userId)
		}
	} else {
		c.execEnv[i.Member.User.ID] = newEvalEnvironment()
		content = "a new environment has been initialized"
		log.Infof("user %s initialized a new environment", userId)
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
