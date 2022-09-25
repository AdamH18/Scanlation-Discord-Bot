package main

import (
	"fmt"
	"scanlation-discord-bot/bot"
	"scanlation-discord-bot/config"
)

// https://golangexample.com/discord-bot-in-golang/
func main() {
	err := config.ReadConfig()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	bot.Start()

	<-make(chan struct{})
	return
}
