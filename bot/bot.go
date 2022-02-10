package main

import (
	"context"
	"emobot/bot/cmd"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"os"
	"os/signal"
	"time"
)

var s *discordgo.Session
var db *mongo.Client

func init() {
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	mongoUri := os.Getenv("MONGO_URI")

	db, err = mongo.Connect(ctx, options.Client().ApplyURI(mongoUri))

	if err != nil {
		log.Fatalln(err)
	}

	err = db.Ping(ctx, readpref.Primary())

	if err != nil {
		log.Fatalln(err)
	}

	log.Println("connected to mongo")
}

// init discord session
func init() {
	var err error

	botToken := os.Getenv("BOT_TOKEN")

	s, err = discordgo.New("Bot " + botToken)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("new discord session created")

	err = s.Open()
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {

	defer s.Close()
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		err := db.Disconnect(ctx)
		if err != nil {
			log.Println(err)
		}
	}()

	stickerCommand := cmd.NewStickerSlashCommand(db)
	addStickerCommand := cmd.NewAddStickerSlashCommand(db)

	cmdHandlerMap, cmdDefinitionList := cmd.PrepareCommands(stickerCommand, addStickerCommand)

	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) { log.Println("Bot is up!") })
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if handler, ok := cmdHandlerMap[i.ApplicationCommandData().Name]; ok {
			handler(s, i)
		}
	})

	guildId := os.Getenv("GUILD_ID")

	createdCommands, err := s.ApplicationCommandBulkOverwrite(s.State.User.ID, guildId, cmdDefinitionList)

	if err != nil {
		log.Println(err)
	}

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt) //nolint: staticcheck
	<-stop
	log.Println("Gracefully shutting down")

	for _, c := range createdCommands {
		err = s.ApplicationCommandDelete(s.State.User.ID, guildId, c.ID)
		if err != nil {
			log.Fatalf("Cannot delete %q command: %v", c.Name, err)
		}
	}
}
