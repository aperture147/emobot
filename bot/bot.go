package main

import (
	"context"
	"emobot/bot/cmd"
	"emobot/bot/cmd/sticker"
	"emobot/bot/db"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"os"
	"os/signal"
	"time"
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
	session, err = cmd.NewDiscordSession()
	if err != nil {
		log.Fatalln("cannot create discord token, ", err)
	}
}

func main() {
	defer func() {
		err := session.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		err := client.Disconnect(ctx)
		if err != nil {
			log.Println(err)
		}
	}()

	guildCollection := client.Database("global").Collection("guild")
	guildIdList, err := db.GetGuildIdList(guildCollection)
	if err != nil {
		log.Fatalln(err)
	}

	userId := session.State.User.ID

	for _, guildId := range guildIdList {
		guildDatabase := db.GetGuildDatabase(client, guildId)
		stickerCollection := guildDatabase.Collection("sticker")
		stickerCommands := sticker.NewAllStickerCommands(stickerCollection)
		masterCmdHandler, cmdDefinitionList := cmd.PrepareCommands(guildId, stickerCommands...)

		createdCommands, err := session.ApplicationCommandBulkOverwrite(userId, guildId, cmdDefinitionList)
		session.AddHandler(masterCmdHandler)
		if err != nil {
			log.Println("failed to add command to guild "+guildId+" with reason: ", err)
		}
		log.Println("slash command added for guild " + guildId)
		defer cmd.DeleteCreatedCommand(session, guildId, createdCommands)
	}

	if err != nil {
		log.Fatalln(err)
	}

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("gracefully shutting down")
}
