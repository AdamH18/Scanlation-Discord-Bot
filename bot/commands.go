package bot

import "github.com/bwmarrin/discordgo"

var (
	adminPerms int64 = discordgo.PermissionAdministrator
	dmPerms          = false
	daysMin          = 1.0
	hourModMin       = -12.0

	//Definitions for all slash commands and their expected parameters
	commands = []*discordgo.ApplicationCommand{
		{
			Name:                     "add_any_reminder",
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
					Description: "Time until reminder in days",
					MinValue:    &daysMin,
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "message",
					Description: "Message to user (max 100 char)",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionBoolean,
					Name:        "repeat",
					Description: "Should this reminder repeat?",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "hours-mod",
					Description: "Modify reminder time by plus or minus 12 hours",
					MinValue:    &hourModMin,
					MaxValue:    12.0,
					Required:    false,
				},
			},
		},
		{
			Name:         "add_reminder",
			Description:  "Add reminder for yourself",
			DMPermission: &dmPerms,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "days",
					Description: "Time until reminder in days",
					MinValue:    &daysMin,
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "message",
					Description: "Message to self (max 100 char)",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionBoolean,
					Name:        "repeat",
					Description: "Should this reminder repeat?",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "hours-mod",
					Description: "Modify reminder time by plus or minus 12 hours",
					MinValue:    &hourModMin,
					MaxValue:    12.0,
					Required:    false,
				},
			},
		},
		{
			Name:                     "rem_any_reminder",
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
			Name:         "rem_reminder",
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
			Name:                     "user_reminders",
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
		{
			Name:                     "set_any_alarm",
			Description:              "Set alarm for any user (admin only)",
			DMPermission:             &dmPerms,
			DefaultMemberPermissions: &adminPerms,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user",
					Description: "User to add alarm for",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "date-time",
					Description: "Date and time of alarm. Must follow format 'YYYY-MM-DD HH:MM:SS'",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "message",
					Description: "Message to user (max 100 char)",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "days",
					Description: "If you want this alarm to repeat every X days, add this",
					MinValue:    &daysMin,
					Required:    false,
				},
			},
		},
		{
			Name:         "set_alarm",
			Description:  "Set an alarm for yourself",
			DMPermission: &dmPerms,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "date-time",
					Description: "Date and time of alarm. Must follow format 'YYYY-MM-DD HH:MM:SS' (GMT)",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "message",
					Description: "Message to self (max 100 char)",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "days",
					Description: "If you want this alarm to repeat every X days, add this",
					MinValue:    &daysMin,
					Required:    false,
				},
			},
		},
	}

	//Map to link slash commands to their handler
	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"add_any_reminder": AddAnyReminderHandler,
		"add_reminder":     AddReminderHandler,
		"rem_any_reminder": RemAnyReminderHandler,
		"rem_reminder":     RemReminderHandler,
		"my_reminders":     MyRemindersHandler,
		"user_reminders":   UserRemindersHandler,
		"all_reminders":    AllRemindersHandler,
		"set_any_alarm":    SetAnyAlarmHandler,
		"set_alarm":        SetAlarmHandler,
	}
)
