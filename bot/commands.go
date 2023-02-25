package bot

import "github.com/bwmarrin/discordgo"

var (
	adminPerms int64 = discordgo.PermissionAdministrator
	dmPerms          = false

	//Definitions for all slash commands and their expected parameters
	commands = []*discordgo.ApplicationCommand{
		{
			Name:                     "add_reminder",
			Description:              "Add reminder for a user (admin only)",
			DMPermission:             &dmPerms,
			DefaultMemberPermissions: &adminPerms,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user",
					Description: "User to be reminded",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "days",
					Description: "Frequency of reminder in days",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "message",
					Description: "Message to user",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionBoolean,
					Name:        "repeat",
					Description: "Should this reminder repeat?",
					Required:    false,
				},
			},
		},
		{
			Name:         "add_personal_reminder",
			Description:  "Add reminder for yourself",
			DMPermission: &dmPerms,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "days",
					Description: "Frequency of reminder in days",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "message",
					Description: "Message to user",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionBoolean,
					Name:        "repeat",
					Description: "Should this reminder repeat?",
					Required:    false,
				},
			},
		},
		{
			Name:                     "rem_reminder",
			Description:              "Remove reminder for any user (admin only)",
			DMPermission:             &dmPerms,
			DefaultMemberPermissions: &adminPerms,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "id",
					Description: "ID of reminder to be removed",
					Required:    true,
				},
			},
		},
		{
			Name:         "rem_personal_reminder",
			Description:  "Remove reminder for yourself",
			DMPermission: &dmPerms,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "id",
					Description: "ID of reminder to be removed (can be found through my_reminders command)",
					Required:    true,
				},
			},
		},
		{
			Name:         "my_reminders",
			Description:  "Show all personal reminders",
			DMPermission: &dmPerms,
		},
		{
			Name:                     "user_reminder",
			Description:              "Show all reminders for a user (admin only)",
			DMPermission:             &dmPerms,
			DefaultMemberPermissions: &adminPerms,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user",
					Description: "User to check reminders for",
					Required:    true,
				},
			},
		},
		{
			Name:                     "all_reminders",
			Description:              "Show all reminders (admin only)",
			DMPermission:             &dmPerms,
			DefaultMemberPermissions: &adminPerms,
		},
	}

	//Map to link slash commands to their handler
	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"add_reminder":          AddReminderHandler,
		"add_personal_reminder": AddPersonalReminderHandler,
		"rem_reminder":          RemReminderHandler,
		"rem_personal_reminder": RemPersonalReminderHandler,
		"my_reminders":          MyRemindersHandler,
		"user_reminders":        UserRemindersHandler,
		"all_reminders":         AllRemindersHandler,
	}
)
