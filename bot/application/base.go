package application

import (
	"github.com/bwmarrin/discordgo"
	"log"
)

type CommandCollection interface {
	GetAllCommands() []Command
}

type Command interface {
	Name() string
	Definition() *discordgo.ApplicationCommand
	Handler(s *discordgo.Session, i *discordgo.InteractionCreate)
}

func PrepareHandler(guildId string, cmdList ...Command) (func(s *discordgo.Session, i *discordgo.InteractionCreate), []*discordgo.ApplicationCommand) {
	cmdHandlerMap := make(map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate), len(cmdList))
	cmdDefinitionList := make([]*discordgo.ApplicationCommand, len(cmdList))

	for index, cmd := range cmdList {
		cmdHandlerMap[cmd.Name()] = cmd.Handler
		cmdDefinitionList[index] = cmd.Definition()
	}

	masterCmdHandler := func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if guildId != i.GuildID {
			return
		}
		applicationCommandName := i.ApplicationCommandData().Name
		if handler, ok := cmdHandlerMap[i.ApplicationCommandData().Name]; ok {
			handler(s, i)
		} else {
			log.Printf("received unknown application command %q\n", applicationCommandName)
		}
	}

	return masterCmdHandler, cmdDefinitionList
}

func DeleteGuildCreatedCommands(s *discordgo.Session, guildId string, createdCommands []*discordgo.ApplicationCommand) {
	var err error
	userId := s.State.User.ID
	for _, c := range createdCommands {
		err = s.ApplicationCommandDelete(userId, guildId, c.ID)
		if err != nil {
			log.Printf("cannot delete %q command on guild %s, %v\n", c.Name, guildId, err)
		}
	}
	log.Printf("deleted all commands on guild %s\n", guildId)
}