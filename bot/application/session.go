package application

import (
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

func NewDiscordSession() (*discordgo.Session, error) {
	botToken := os.Getenv("BOT_TOKEN")

	s, err := discordgo.New("Bot " + botToken)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Println("new discord session created, latency", time.Since(start).Milliseconds(), "ms")
	})

	err = s.Open()
	if err != nil {
		return nil, err
	}

	return s, nil
}

func CloseDiscordSession(session *discordgo.Session) {
	err := session.Close()
	if err != nil {
		log.Println(err)
	}
}
