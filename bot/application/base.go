package application

import (
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
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
		if handler, ok := cmdHandlerMap[applicationCommandName]; ok {
			defer func() {
				if r := recover(); r != nil {
					log.WithField("guild_id", i.GuildID).Warningf("Recovered in handler %s with reason %s", applicationCommandName, r)
				}
			}()
			handler(s, i)
		} else {
			log.WithField("guild_id", i.GuildID).Printf("received unknown application command %q", applicationCommandName)
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
			log.WithField("guild_id", guildId).Printf("cannot delete %q command on guild %s, %v", c.Name, guildId, err)
		}
	}
	log.WithField("guild_id", guildId).Printf("deleted all commands on guild %s", guildId)
}
