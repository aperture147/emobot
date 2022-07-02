package server

import (
	"context"
	"emobot/bot/application"
	"fmt"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

type pingCommand struct {
	client *mongo.Client
}

const pingCommandName = "ping"

var pingCommandDefinition = &discordgo.ApplicationCommand{
	Name:        pingCommandName,
	Description: "send ping signal to database and discord and get the latency between them",
	Type:        discordgo.ChatApplicationCommand,
}

func NewPingCommand(client *mongo.Client) application.Command {
	return &pingCommand{client}
}

func (c *pingCommand) Name() string {
	return pingCommandName
}

func (c *pingCommand) Definition() *discordgo.ApplicationCommand {
	return pingCommandDefinition
}

func getDatabaseLatency(client *mongo.Client) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	databaseStart := time.Now()
	err := client.Ping(ctx, readpref.Primary())
	return time.Since(databaseStart).Milliseconds(), err
}

const (
	databaseLatencyIcon = ":robot: → :leaves:"
	botLatencyIcon      = ":robot: → :office:"
	userLatencyIcon     = ":computer: → :robot:"
)

const warningIcon = ":warning:"

func (c *pingCommand) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	databaseLatency, err := getDatabaseLatency(c.client)
	var databaseLatencyStr string
	if err != nil {
		log.WithField("guild_id", i.GuildID).Println("cannot ping database with reason:", err)
		databaseLatencyStr = fmt.Sprintf("%s: %s\n", databaseLatencyIcon, warningIcon)
	} else {
		databaseLatencyStr = fmt.Sprintf("%s: %d ms\n", databaseLatencyIcon, databaseLatency)
	}

	botLatency := s.HeartbeatLatency().Milliseconds()
	botLatencyStr := fmt.Sprintf("%s: %d ms\n", botLatencyIcon, botLatency)

	msgTimestamp, _ := discordgo.SnowflakeTimestamp(i.ID)
	userLatency := time.Since(msgTimestamp).Milliseconds()
	userLatencyStr := fmt.Sprintf("%s: %d ms\n", userLatencyIcon, userLatency)

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: databaseLatencyStr + botLatencyStr + userLatencyStr,
		},
	})
	if err != nil {
		log.WithField("guild_id", i.GuildID).Println("cannot send pong with reason:", err)
	}
}
