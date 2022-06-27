package main

import (
	"emobot/bot/application"
	"emobot/bot/db"
	"emobot/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/elastic/go-elasticsearch/v7"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/go-extras/elogrus.v7"
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

// init logrus elasticsearch hook
func init() {
	elasticClient, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{os.Getenv("BONSAI_URL")},
	})
	if err != nil {
		log.Fatalln("cannot connect to elastic with reason:", err)
	} else {
		log.Infoln("elasticsearch client created")
	}
	hook, err := elogrus.NewAsyncElasticHook(elasticClient, "localhost", log.InfoLevel, "emobot")
	if err != nil {
		log.Fatalln("cannot using elasticsearch hook with reason:", err)
	} else {
		log.Infoln("elasticsearch hook added")
	}
	log.AddHook(hook)
}

func main() {
	defer application.CloseDiscordSession(session)
	defer db.CloseMongoClient(client)

	guildCollection := client.Database("global").Collection("guild")
	guildIdList, err := db.GetGuildIdList(guildCollection)
	if err != nil {
		log.Fatalln(err)
	}

	guildCreatedCommands := utils.AddGuildCommands(session, client, guildIdList)

	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	log.Println("gracefully shutting down")

	utils.DeleteGuildCommands(session, guildCreatedCommands)
}
