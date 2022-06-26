package utils

import (
	"emobot/bot/application"
	"emobot/bot/application/chat"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

func AddGuildCommands(session *discordgo.Session, client *mongo.Client, guildIdList []string) map[string][]*discordgo.ApplicationCommand {
	guildCreatedCommands := make(map[string][]*discordgo.ApplicationCommand, len(guildIdList))
	for _, guildId := range guildIdList {
		go func(guildId string) {
			chatCollection := chat.NewCommandCollection(guildId, client)
			//messageCollection := message.NewCommandCollection(guildId, client)

			var commands []application.Command
			commands = append(commands, chatCollection.GetAllCommands()...)
			//commands = append(commands, messageCollection.GetAllCommands()...)

			masterCmdHandler, cmdDefinitionList := application.PrepareHandler(guildId, commands...)
			createdCommands, err := session.ApplicationCommandBulkOverwrite(session.State.User.ID, guildId, cmdDefinitionList)
			session.AddHandler(masterCmdHandler)

			if err != nil {
				log.Println("failed to add command to guild", guildId, "with reason:", err)
			} else {
				log.Println("application command added for guild", guildId)
			}
			guildCreatedCommands[guildId] = createdCommands
		}(guildId)
	}
	return guildCreatedCommands
}

func DeleteGuildCommands(session *discordgo.Session, guildCreatedCommands map[string][]*discordgo.ApplicationCommand) {
	wg, ctx, cancel := WaitGroupTimeOut(60 * time.Second)
	defer cancel()
	for guildId, createdCommands := range guildCreatedCommands {
		wg.Add(1)
		go func(guildId string, createdCommands []*discordgo.ApplicationCommand) {
			defer wg.Done()
			application.DeleteGuildCreatedCommands(session, guildId, createdCommands)
		}(guildId, createdCommands)
	}
	wg.Wait()
	<-ctx.Done()
}
