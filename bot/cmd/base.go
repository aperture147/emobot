package cmd

import "github.com/bwmarrin/discordgo"

type SlashCommand interface {
	Name() string
	Definition() *discordgo.ApplicationCommand
	Handler(s *discordgo.Session, i *discordgo.InteractionCreate)
}

func PrepareCommands(cmdList ...SlashCommand) (map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate), []*discordgo.ApplicationCommand) {
	cmdHandlerMap := make(map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate), len(cmdList))
	cmdDefinitionList := make([]*discordgo.ApplicationCommand, len(cmdList))

	for index, cmd := range cmdList {
		cmdHandlerMap[cmd.Name()] = cmd.Handler
		cmdDefinitionList[index] = cmd.Definition()
	}

	return cmdHandlerMap, cmdDefinitionList
}
