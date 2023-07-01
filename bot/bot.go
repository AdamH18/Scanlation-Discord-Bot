package bot

import (
	//to print errors
	"log"
	"scanlation-discord-bot/config" //importing our config package which we have created above
	"scanlation-discord-bot/database"

	"github.com/bwmarrin/discordgo" //discordgo package from the repo of bwmarrin .
)

var (
	BotId string
	goBot *discordgo.Session
	quit  chan struct{}

	SeriesCh      chan func() (string, string)
	AssignmentsCh chan string
	ColorsCh      chan string
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

	//Get commands registed with Discord
	cmds, err := goBot.ApplicationCommands(BotId, "")
	if err != nil {
		log.Panicf("Was unable to retrieve existing commands: %v\n", err)
	}

	log.Println("Adding commands...")
	for _, v := range commands {
		if !DiscordCommand(cmds, v.Name) || config.RefreshCommands {
			_, err := goBot.ApplicationCommandCreate(BotId, "", v)
			if err != nil {
				log.Panicf("Cannot create '%v' command: %v", v.Name, err)
			}
			log.Printf("Command %s added successfully\n", v.Name)
		}
	}

	log.Println("Removing unused commands...")
	for _, cmd := range cmds {
		if !IsCommand(cmd.Name) {
			log.Printf("Unused command %v found. Deleting...", cmd.Name)
			err := goBot.ApplicationCommandDelete(BotId, "", cmd.ID)
			if err != nil {
				log.Printf("Cannot delete '%v' command: %v", cmd.Name, err)
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

	//Create channels for billboard updates
	SeriesCh = make(chan func() (string, string))
	AssignmentsCh = make(chan string)
	ColorsCh = make(chan string)
	database.RegisterChannels(SeriesCh, AssignmentsCh, ColorsCh)
	go BillboardUpdates()
}

func Stop() {
	defer goBot.Close()
	//Close all goroutines
	close(quit)

	if config.RemoveCommands {
		log.Println("Removing commands...")
		cmds, err := goBot.ApplicationCommands(BotId, "")
		if err != nil {
			log.Printf("Was unable to retrieve existing commands: %v\n", err)
		}
		for _, v := range cmds {
			err := goBot.ApplicationCommandDelete(BotId, "", v.ID)
			if err != nil {
				log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
			}
		}
	}
}
