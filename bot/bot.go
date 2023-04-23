package bot

import (
	//to print errors
	"log"
	"scanlation-discord-bot/config" //importing our config package which we have created above

	"github.com/bwmarrin/discordgo" //discordgo package from the repo of bwmarrin .
)

var (
	BotId              string
	goBot              *discordgo.Session
	registeredCommands []*discordgo.ApplicationCommand
	quit               chan struct{}
)

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

	log.Println("Removing unused commands...")
	cmds, err := goBot.ApplicationCommands(BotId, "")
	if err != nil {
		log.Printf("Was unable to retrieve existing commands: %v\n", err)
	} else {
		for _, cmd := range cmds {
			if !IsCommand(cmd.Name) {
				log.Printf("Unused command %v found. Deleting...", cmd.Name)
				err := goBot.ApplicationCommandDelete(BotId, "", cmd.ID)
				if err != nil {
					log.Printf("Cannot delete '%v' command: %v", cmd.Name, err)
				}
			}
		}
	}

	//If every thing works fine we will be printing this.
	log.Println("Bot is running!")

	//Quit channel used to close all goroutines
	quit = make(chan struct{})
	log.Println("Starting reminder checking")
	go CheckReminders()
	log.Println("Starting backups")
	go BackupDB()
}

func Stop() {
	defer goBot.Close()
	//Close all goroutines
	close(quit)

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
