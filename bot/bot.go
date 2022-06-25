package main

import (
	"emobot/bot/application"
	"emobot/bot/db"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var session *discordgo.Session
var client *mongo.Client

// init mongo client
func init() {
	var err error
	client, err = db.NewMongoClient()
	if err != nil {
		log.Fatalln("cannot connect to mongo, ", err)
	}
}

// init discord session
func init() {
	var err error
	session, err = application.NewDiscordSession()
	if err != nil {
		log.Fatalln("cannot create discord token, ", err)
	}
}

func main() {
	defer application.CloseDiscordSession(session)
	defer db.CloseMongoClient(client)

	guildCollection := client.Database("global").Collection("guild")
	guildIdList, err := db.GetGuildIdList(guildCollection)
	if err != nil {
		log.Fatalln(err)
	}

	guildCreatedCommands := AddGuildCommands(guildIdList)

	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	log.Println("gracefully shutting down")

	DeleteGuildCommands(guildCreatedCommands)
}
