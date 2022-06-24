package cmd

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
)

func NewDiscordSession() (*discordgo.Session, error) {
	botToken := os.Getenv("BOT_TOKEN")

	s, err := discordgo.New("Bot " + botToken)
	if err != nil {
		return nil, err
	}

	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) { log.Println("new discord session created") })

	err = s.Open()
	if err != nil {
		return nil, err
	}

	return s, nil
}
