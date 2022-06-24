package server

import (
	"context"
	"emobot/bot/cmd"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

type PingSlashCommand struct {
	client *mongo.Client
}

const PingCommandName = "ping"

var PingCommandDefinition = &discordgo.ApplicationCommand{
	Name:        PingCommandName,
	Description: "send ping signal to database and discord",
	Type:        discordgo.ChatApplicationCommand,
}

func NewPingCommand(client *mongo.Client) cmd.SlashCommand {
	return &PingSlashCommand{client}
}

func (c *PingSlashCommand) Name() string {
	return PingCommandName
}

func (c *PingSlashCommand) Definition() *discordgo.ApplicationCommand {
	return PingCommandDefinition
}

func (c *PingSlashCommand) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	databaseStart := time.Now()
	err := c.client.Ping(ctx, readpref.Primary())
	databaseLatency := time.Since(databaseStart).Milliseconds()
	if err != nil {
		log.Println("cannot ping database with reason:", err)
	}
	botLatency := s.HeartbeatLatency().Milliseconds()
	msgTimestamp, _ := discordgo.SnowflakeTimestamp(i.ID)
	userLatency := time.Since(msgTimestamp).Milliseconds()

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf(":leaves: %d ms\n:robot: %d ms\n:ping_pong: %d ms", databaseLatency, botLatency, userLatency),
		},
	})
	if err != nil {
		log.Println("cannot send pong with reason:", err)
	}
}
