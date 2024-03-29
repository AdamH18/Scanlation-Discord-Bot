package main

import (
	"log"
	"os"
	"os/signal"
	"scanlation-discord-bot/bot"
	"scanlation-discord-bot/config"
	"scanlation-discord-bot/database"
)

func init() {
	err := config.ReadConfig()

	if err != nil {
		log.Println(err.Error())
		return
	}
}

// https://golangexample.com/discord-bot-in-golang/
func main() {
	database.StartDatabase(config.DatabaseFile)
	bot.Start()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	bot.Stop()
	log.Println("Gracefully shutting down.")
}
