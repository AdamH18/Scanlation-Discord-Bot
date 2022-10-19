package bot

import "github.com/bwmarrin/discordgo"

var (
	adminPerms int64 = discordgo.PermissionAdministrator
	dmPerms          = false

	//Definitions for all slash commands and their expected parameters
	commands = []*discordgo.ApplicationCommand{
		{
			Name:         "test",
			Description:  "Test command",
			DMPermission: &dmPerms,
		},
		{
			Name:                     "test_restricted",
			Description:              "Should only be available if you're an admin",
			DefaultMemberPermissions: &adminPerms,
			DMPermission:             &dmPerms,
		},
	}

	//Map to link slash commands to their handler
	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"test": TestHandler,
	}
)
