package bot

import "github.com/bwmarrin/discordgo"

var (
	//Definitions for all slash commands and their expected parameters
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "test",
			Description: "Test command",
		},
	}

	//Map to link slash commands to their handler
	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"test": TestHandler,
	}
)
