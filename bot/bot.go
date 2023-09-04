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
	ActionsCh     chan bool
	ErrorsCh      chan func() (string, []any, string)

	DatabaseOps  int
	DatabaseErrs int
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
	//Handles all command additions and deletions in a single function. Why wasn't I using this from the start?
	_, err = goBot.ApplicationCommandBulkOverwrite(BotId, "", commands)
	if err != nil {
		log.Println("Failed to load application commands: " + err.Error())
		return
	}

	InitializeCache()

	//If every thing works fine we will be printing this.
	log.Println("Bot is running!")

	//Quit channel used to close all goroutines
	quit = make(chan struct{})
	log.Println("Starting reminder checking")
	go CheckReminders()
	log.Println("Starting backups")
	go BackupDB()

	//Create channels
	SeriesCh = make(chan func() (string, string))
	AssignmentsCh = make(chan string)
	ColorsCh = make(chan string)
	ActionsCh = make(chan bool)
	ErrorsCh = make(chan func() (string, []any, string))
	DatabaseOps = 0
	DatabaseErrs = 0
	database.RegisterChannels(SeriesCh, AssignmentsCh, ColorsCh, ActionsCh, ErrorsCh)
	go TrackDB()
	go HandlerErrors()
	go BillboardUpdates()

	goBot.UpdateListeningStatus("my Fans")
}

func Stop() {
	defer goBot.Close()
	//Close all goroutines
	close(quit)

	if config.RemoveCommands {
		log.Println("Removing commands...")
		_, err := goBot.ApplicationCommandBulkOverwrite(BotId, "", []*discordgo.ApplicationCommand{})
		if err != nil {
			log.Println("Failed to remove commands: " + err.Error())
		}
	}
}
