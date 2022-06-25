package main

import (
	"emobot/bot/application"
	"emobot/bot/application/chat"
	"github.com/bwmarrin/discordgo"
	"log"
	"sync"
)

func AddGuildCommands(guildIdList []string) map[string][]*discordgo.ApplicationCommand {
	guildCreatedCommands := make(map[string][]*discordgo.ApplicationCommand, len(guildIdList))
	for _, guildId := range guildIdList {
		go func(guildId string) {
			chatCollection := chat.NewCommandCollection(guildId, client)
			var commands []application.Command
			commands = append(commands, chatCollection.GetAllCommands()...)

			masterCmdHandler, cmdDefinitionList := application.PrepareHandler(guildId, commands...)
			createdCommands, err := session.ApplicationCommandBulkOverwrite(session.State.User.ID, guildId, cmdDefinitionList)
			session.AddHandler(masterCmdHandler)

			if err != nil {
				log.Println("failed to add command to guild "+guildId+" with reason: ", err)
			}
			log.Println("slash command added for guild " + guildId)
			guildCreatedCommands[guildId] = createdCommands
		}(guildId)
	}
	return guildCreatedCommands
}

func DeleteGuildCommands(guildCreatedCommands map[string][]*discordgo.ApplicationCommand) {
	var wg sync.WaitGroup
	for guildId, createdCommands := range guildCreatedCommands {
		wg.Add(1)
		go func(guildId string, createdCommands []*discordgo.ApplicationCommand) {
			defer wg.Done()
			application.DeleteGuildCreatedCommands(session, guildId, createdCommands)
		}(guildId, createdCommands)
	}
	wg.Wait()
}
