package utils

import (
	"emobot/bot/application"
	"emobot/bot/application/chat"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func AddGuildCommands(session *discordgo.Session, client *mongo.Client, guildIdList []string) map[string][]*discordgo.ApplicationCommand {
	guildCreatedCommands := make(map[string][]*discordgo.ApplicationCommand, len(guildIdList))
	for _, guildId := range guildIdList {
		go func(guildId string) {
			var commands []application.Command

			chatCollection := chat.NewCommandCollection(guildId, client)
			commands = append(commands, chatCollection.GetAllCommands()...)

			//messageCollection := message.NewCommandCollection(guildId, client)
			//commands = append(commands, messageCollection.GetAllCommands()...)

			masterCmdHandler, cmdDefinitionList := application.PrepareHandler(guildId, commands...)
			createdCommands, err := session.ApplicationCommandBulkOverwrite(session.State.User.ID, guildId, cmdDefinitionList)
			session.AddHandler(masterCmdHandler)

			if err != nil {
				log.WithField("guild_id", guildId).Println("failed to add command to guild", guildId, "with reason:", err)
			} else {
				log.WithField("guild_id", guildId).Println("application command added for guild", guildId)
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
