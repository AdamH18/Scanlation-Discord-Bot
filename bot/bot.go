package bot

import (
	//to print errors
	"log"
	"scanlation-discord-bot/config" //importing our config package which we have created above

	"github.com/bwmarrin/discordgo" //discordgo package from the repo of bwmarrin .
)

var BotId string
var goBot *discordgo.Session
var registeredCommands []*discordgo.ApplicationCommand

func Start() {

	//Creating new bot session
	var err error
	goBot, err = discordgo.New("Bot " + config.Token)
	if err != nil {
		log.Printf("Invalid bot parameters: %v\n", err)
	}

	//Register slash command handlers
	CreateHandlers()

	//Extract bot id
	u, err := goBot.User("@me")
	//Handling error
	if err != nil {
		log.Println(err.Error())
		return
	}
	// Storing our id from u to BotId .
	BotId = u.ID

	err = goBot.Open()
	//Error handling
	if err != nil {
		log.Printf("Failed to open socket: %v\n", err.Error())
		return
	}

	log.Println("Adding commands...")
	registeredCommands = make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := goBot.ApplicationCommandCreate(BotId, "", v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	//If every thing works fine we will be printing this.
	log.Println("Bot is running !")
}

func Stop() {
	defer goBot.Close()

	if config.RemoveCommands {
		log.Println("Removing commands...")
		for _, v := range registeredCommands {
			err := goBot.ApplicationCommandDelete(BotId, "", v.ID)
			if err != nil {
				log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
			}
		}
	}
}
