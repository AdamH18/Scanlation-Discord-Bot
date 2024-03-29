package bot

import "github.com/bwmarrin/discordgo"

type AutoComplete struct {
	Loc     int
	Choices func(string, string) []*discordgo.ApplicationCommandOptionChoice
}

var (
	adminPerms int64 = discordgo.PermissionAdministrator
	dmPerms          = false
	daysMin          = 0.0
	hourModMin       = -12.0
	emptyStr         = ""

	//Definitions for all slash commands and their expected parameters
	commands = []*discordgo.ApplicationCommand{
		{
			Name:         "help",
			Description:  "Get info on how to use this bot",
			DMPermission: &dmPerms,
		},

		// REMINDER COMMANDS
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

		// REGISTRATION COMMANDS
		{
			Name:                     "add_series",
			Description:              "Register new series for group (admin only)",
			DMPermission:             &dmPerms,
			DefaultMemberPermissions: &adminPerms,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "full-name",
					Description: "Full name for the series",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "short-name",
					Description: "Shorthand name for the series (used for role and category names)",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionBoolean,
					Name:        "full-create",
					Description: "Should channels and roles be created for this series",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionRole,
					Name:        "ping-role",
					Description: "If full-create was not selected, include role for pinging on release here",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "repo-link",
					Description: "Link to repo where series files can be found",
					Required:    false,
				},
			},
		},
		{
			Name:                     "remove_series",
			Description:              "Removes series for group, including all related settings. Channels are not deleted (admin only)",
			DMPermission:             &dmPerms,
			DefaultMemberPermissions: &adminPerms,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:         discordgo.ApplicationCommandOptionString,
					Name:         "full-name",
					Description:  "Full name for the series",
					Required:     true,
					Autocomplete: true,
				},
				{
					Type:         discordgo.ApplicationCommandOptionString,
					Name:         "short-name",
					Description:  "Shorthand name for the series (just to make sure)",
					Required:     true,
					Autocomplete: true,
				},
			},
		},
		{
			Name:                     "server_series",
			Description:              "See all existing series on the server (admin only)",
			DMPermission:             &dmPerms,
			DefaultMemberPermissions: &adminPerms,
		},
		{
			Name:                     "change_series_title",
			Description:              "Changes the full name of the series. Shorthand name is unchanged (admin only)",
			DMPermission:             &dmPerms,
			DefaultMemberPermissions: &adminPerms,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:         discordgo.ApplicationCommandOptionString,
					Name:         "short-name",
					Description:  "Shorthand name for the series",
					Required:     true,
					Autocomplete: true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "new-full-name",
					Description: "New full name for the series",
					Required:    true,
				},
			},
		},
		{
			Name:                     "change_series_repo",
			Description:              "Changes the repo link of the series (admin only)",
			DMPermission:             &dmPerms,
			DefaultMemberPermissions: &adminPerms,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:         discordgo.ApplicationCommandOptionString,
					Name:         "short-name",
					Description:  "Shorthand name for the series",
					Required:     true,
					Autocomplete: true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "new-repo-link",
					Description: "New repo link for the series",
					Required:    true,
				},
			},
		},
		{
			Name:                     "add_series_channel",
			Description:              "Register a channel with a given series (contextual command, admin only)",
			DMPermission:             &dmPerms,
			DefaultMemberPermissions: &adminPerms,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:         discordgo.ApplicationCommandOptionString,
					Name:         "series",
					Description:  "Shorthand for series to add to",
					Required:     true,
					Autocomplete: true,
				},
			},
		},
		{
			Name:                     "remove_series_channel",
			Description:              "Deregister a channel with a given series, channel is not deleted (contextual command, admin only)",
			DMPermission:             &dmPerms,
			DefaultMemberPermissions: &adminPerms,
		},
		{
			Name:                     "add_user",
			Description:              "Register a user as a member of the group (admin only)",
			DMPermission:             &dmPerms,
			DefaultMemberPermissions: &adminPerms,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user",
					Description: "The user to be added",
					Required:    true,
				},
			},
		},
		{
			Name:                     "remove_user",
			Description:              "Remove a user from the group, deletes all related settings. User is not kicked (admin only)",
			DMPermission:             &dmPerms,
			DefaultMemberPermissions: &adminPerms,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user",
					Description: "The user to be removed",
					Required:    true,
				},
			},
		},
		{
			Name:                     "remove_user_by_id",
			Description:              "Remove a user by ID (if left server), deletes all related settings. User is not kicked (admin only)",
			DMPermission:             &dmPerms,
			DefaultMemberPermissions: &adminPerms,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "user",
					Description: "The id of user to be removed",
					Required:    true,
				},
			},
		},
		{
			Name:                     "server_users",
			Description:              "See all registered users on the server (admin only)",
			DMPermission:             &dmPerms,
			DefaultMemberPermissions: &adminPerms,
		},
		{
			Name:                     "add_job",
			Description:              "Register a new job type for the group (admin only)",
			DMPermission:             &dmPerms,
			DefaultMemberPermissions: &adminPerms,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "name-full",
					Description: "Full name of the new job",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "name-short",
					Description: "Name shorthand (PR, TS, etc.)",
					Required:    true,
				},
			},
		},
		{
			Name:                     "add_global_job",
			Description:              "Register a new job type for all users (owner only)",
			DMPermission:             &dmPerms,
			DefaultMemberPermissions: &adminPerms,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "name-full",
					Description: "Full name of the new job",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "name-short",
					Description: "Name shorthand (PR, TS, etc.)",
					Required:    true,
				},
			},
		},
		{
			Name:                     "remove_job",
			Description:              "Remove a job type for the group, including all assignments to that job (admin only)",
			DMPermission:             &dmPerms,
			DefaultMemberPermissions: &adminPerms,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:         discordgo.ApplicationCommandOptionString,
					Name:         "name-short",
					Description:  "Shorthand for the job",
					Required:     true,
					Autocomplete: true,
				},
			},
		},
		{
			Name:                     "server_jobs",
			Description:              "See all existing jobs on the server (admin only)",
			DMPermission:             &dmPerms,
			DefaultMemberPermissions: &adminPerms,
		},
		{
			Name:                     "add_member_role",
			Description:              "Registers the role used to determine group members (admin only)",
			DMPermission:             &dmPerms,
			DefaultMemberPermissions: &adminPerms,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionRole,
					Name:        "role",
					Description: "Group member role",
					Required:    true,
				},
			},
		},
		{
			Name:                     "remove_member_role",
			Description:              "Deregisters the role used to determine group members. Role is not deleted (admin only)",
			DMPermission:             &dmPerms,
			DefaultMemberPermissions: &adminPerms,
		},
		{
			Name:                     "reg_series_channels",
			Description:              "Registers bounds for series channels. Should be IDs of first and last categories (admin only)",
			DMPermission:             &dmPerms,
			DefaultMemberPermissions: &adminPerms,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "top",
					Description: "ID for first series category",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "bottom",
					Description: "ID for last series category",
					Required:    true,
				},
			},
		},

		// JOB COMMANDS
		{
			Name:                     "add_series_assignment",
			Description:              "Register an assignment to a series for a group member (admin only)",
			DMPermission:             &dmPerms,
			DefaultMemberPermissions: &adminPerms,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user",
					Description: "User to be assigned",
					Required:    true,
				},
				{
					Type:         discordgo.ApplicationCommandOptionString,
					Name:         "job",
					Description:  "Shorthand name for the job",
					Required:     true,
					Autocomplete: true,
				},
				{
					Type:         discordgo.ApplicationCommandOptionString,
					Name:         "series",
					Description:  "Shorthand name for the series",
					Required:     false,
					Autocomplete: true,
				},
			},
		},
		{
			Name:                     "remove_series_assignment",
			Description:              "Remove an assignment to a series for a group member (admin only)",
			DMPermission:             &dmPerms,
			DefaultMemberPermissions: &adminPerms,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user",
					Description: "User to have assignment removed",
					Required:    true,
				},
				{
					Type:         discordgo.ApplicationCommandOptionString,
					Name:         "job",
					Description:  "Shorthand name for the job",
					Required:     true,
					Autocomplete: true,
				},
				{
					Type:         discordgo.ApplicationCommandOptionString,
					Name:         "series",
					Description:  "Shorthand name for the series",
					Required:     false,
					Autocomplete: true,
				},
			},
		},
		{
			Name:                     "remove_all_assignments",
			Description:              "Remove all assignments for a group member. Does not kick member from group (admin only)",
			DMPermission:             &dmPerms,
			DefaultMemberPermissions: &adminPerms,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user",
					Description: "User to have all assignments removed",
					Required:    true,
				},
			},
		},
		{
			Name:         "series_assignments",
			Description:  "See the assignments and user colors for a given series (contextual command)",
			DMPermission: &dmPerms,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:         discordgo.ApplicationCommandOptionString,
					Name:         "series",
					Description:  "Series shorthand if non-contextual",
					Required:     false,
					Autocomplete: true,
				},
			},
		},
		{
			Name:         "my_assignments",
			Description:  "See your personal assignments",
			DMPermission: &dmPerms,
		},
		{
			Name:         "user_assignments",
			Description:  "See the assignments of a given user",
			DMPermission: &dmPerms,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user",
					Description: "User to check assignments for",
					Required:    true,
				},
			},
		},
		{
			Name:         "job_assignments",
			Description:  "See everyone assigned to a given job",
			DMPermission: &dmPerms,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:         discordgo.ApplicationCommandOptionString,
					Name:         "job",
					Description:  "Job to check assignments for (shorthand)",
					Required:     true,
					Autocomplete: true,
				},
			},
		},
		{
			Name:         "tl",
			Description:  "Ping the translator(s) (contextual command)",
			DMPermission: &dmPerms,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "message",
					Description: "Message to include with ping",
					Required:    false,
				},
			},
		},
		{
			Name:         "rd",
			Description:  "Ping the redrawer(s) (contextual command)",
			DMPermission: &dmPerms,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "message",
					Description: "Message to include with ping",
					Required:    false,
				},
			},
		},
		{
			Name:         "ts",
			Description:  "Ping the typesetter(s) (contextual command)",
			DMPermission: &dmPerms,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "message",
					Description: "Message to include with ping",
					Required:    false,
				},
			},
		},
		{
			Name:         "pr",
			Description:  "Ping the proofreader(s) (contextual command)",
			DMPermission: &dmPerms,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "message",
					Description: "Message to include with ping",
					Required:    false,
				},
			},
		},

		// USER CUSTOMIZATION
		{
			Name:         "my_settings",
			Description:  "See your server settings",
			DMPermission: &dmPerms,
		},
		{
			Name:                     "user_settings",
			Description:              "See user's server settings",
			DMPermission:             &dmPerms,
			DefaultMemberPermissions: &adminPerms,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user",
					Description: "User to see settings for",
					Required:    true,
				},
			},
		},
		{
			Name:         "set_color",
			Description:  "Set your color for credits pages",
			DMPermission: &dmPerms,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "color",
					Description: "Set your color in RGB hex",
					Required:    true,
				},
			},
		},
		{
			Name:                     "set_user_color",
			Description:              "Set user color for credits pages (admin only)",
			DMPermission:             &dmPerms,
			DefaultMemberPermissions: &adminPerms,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user",
					Description: "User to set color for",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "color",
					Description: "Set your color in RGB hex",
					Required:    true,
				},
			},
		},
		{
			Name:         "vanity_role",
			Description:  "Give yourself a vanity role (or edit existing one)",
			DMPermission: &dmPerms,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "name",
					Description: "Name of the role",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "color",
					Description: "Set role color in RGB hex. Must be exactly 6 characters and parseable",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "copy-user",
					Description: "Take existing vanity role instead. Ignores first two parameters. Works even if at max roles",
					Required:    false,
				},
			},
		},
		{
			Name:         "rem_vanity_role",
			Description:  "Removes your vanity role",
			DMPermission: &dmPerms,
		},

		// BILLBOARD COMMANDS
		{
			Name:                     "create_series_billboard",
			Description:              "Create a billboard showcasing series information in this channel (admin only)",
			DMPermission:             &dmPerms,
			DefaultMemberPermissions: &adminPerms,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:         discordgo.ApplicationCommandOptionString,
					Name:         "series",
					Description:  "Shorthand name for the series",
					Required:     true,
					Autocomplete: true,
				},
			},
		},
		{
			Name:                     "delete_series_billboard",
			Description:              "Deregister the billboard showcasing series information. Does not delete message (admin only)",
			DMPermission:             &dmPerms,
			DefaultMemberPermissions: &adminPerms,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:         discordgo.ApplicationCommandOptionString,
					Name:         "series",
					Description:  "Shorthand name for the series",
					Required:     true,
					Autocomplete: true,
				},
			},
		},
		{
			Name:                     "create_assignments_billboard",
			Description:              "Create a billboard showcasing all assignments in this channel (admin only)",
			DMPermission:             &dmPerms,
			DefaultMemberPermissions: &adminPerms,
		},
		{
			Name:                     "delete_assignments_billboard",
			Description:              "Deregister the billboard showcasing all assignments. Does not delete message (admin only)",
			DMPermission:             &dmPerms,
			DefaultMemberPermissions: &adminPerms,
		},
		{
			Name:                     "create_colors_billboard",
			Description:              "Create a billboard showcasing all color prefs in this channel (admin only)",
			DMPermission:             &dmPerms,
			DefaultMemberPermissions: &adminPerms,
		},
		{
			Name:                     "delete_colors_billboard",
			Description:              "Deregister the billboard showcasing all color prefs. Does not delete message (admin only)",
			DMPermission:             &dmPerms,
			DefaultMemberPermissions: &adminPerms,
		},
		{
			Name:                     "refresh_all_billboards",
			Description:              "Refreshes all billboards on the server (admin only)",
			DMPermission:             &dmPerms,
			DefaultMemberPermissions: &adminPerms,
		},

		// GENERAL COMMANDS
		{
			Name:                     "add_notification_channel",
			Description:              "Sets channel to receive messages containing updates from bot owner (admin only)",
			DMPermission:             &dmPerms,
			DefaultMemberPermissions: &adminPerms,
		},
		{
			Name:                     "send_notification",
			Description:              "Send message to all registered notification channels (owner only)",
			DMPermission:             &dmPerms,
			DefaultMemberPermissions: &adminPerms,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "message",
					Description: "Message to send",
					Required:    true,
				},
			},
		},
	}

	//Map to link slash commands to their handler
	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"help": HelpHandler,
		//"admin_help": AdminHelpHandler,

		//TODO: Organize by subcommands

		"add_any_reminder": AddAnyReminderHandler,
		"add_reminder":     AddReminderHandler,
		"rem_any_reminder": RemAnyReminderHandler,
		"rem_reminder":     RemReminderHandler,
		"my_reminders":     MyRemindersHandler,
		"user_reminders":   UserRemindersHandler,
		"all_reminders":    AllRemindersHandler,
		"set_any_alarm":    SetAnyAlarmHandler,
		"set_alarm":        SetAlarmHandler,

		"add_series":            AddSeriesHandler,
		"remove_series":         RemoveSeriesHandler,
		"server_series":         ServerSeriesHandler,
		"change_series_title":   ChangeSeriesTitleHandler,
		"change_series_repo":    ChangeSeriesRepoHandler,
		"add_series_channel":    AddSeriesChannelHandler,
		"remove_series_channel": RemoveSeriesChannelHandler,
		"add_user":              AddUserHandler,
		"remove_user":           RemoveUserHandler,
		"remove_user_by_id":     RemoveUserByIDHandler,
		"server_users":          ServerUsersHandler,
		"add_job":               AddJobHandler,
		"add_global_job":        AddGlobalJobHandler,
		"remove_job":            RemoveJobHandler,
		"server_jobs":           ServerJobsHandler,
		"add_member_role":       AddMemberRoleHandler,
		"remove_member_role":    RemoveMemberRoleHandler,
		"reg_series_channels":   RegSeriesChannelsHandler,
		//"server_settings": ServerSettingsHandler,

		"add_series_assignment":    AddSeriesAssignmentHandler,
		"remove_series_assignment": RemoveSeriesAssignmentHandler,
		"remove_all_assignments":   RemoveAllAssignmentsHandler,
		"series_assignments":       SeriesAssignmentsHandler,
		"my_assignments":           MyAssignmentsHandler,
		"user_assignments":         UserAssignmentsHandler,
		"job_assignments":          JobAssignmentsHandler,
		"tl":                       TLPingHandler,
		"rd":                       RDPingHandler,
		"ts":                       TSPingHandler,
		"pr":                       PRPingHandler,

		"my_settings":     MySettingsHandler,
		"user_settings":   UserSettingsHandler,
		"set_color":       SetColorHandler,
		"set_user_color":  SetUserColorHandler,
		"vanity_role":     VanityRoleHandler,
		"rem_vanity_role": RemVanityRoleHandler,

		"create_series_billboard":      CreateSeriesBillboardHandler,
		"delete_series_billboard":      DeleteSeriesBillboardHandler,
		"create_assignments_billboard": CreateAssignmentsBillboardHandler,
		"delete_assignments_billboard": DeleteAssignmentsBillboardHandler,
		"create_colors_billboard":      CreateColorsBillboardHandler,
		"delete_colors_billboard":      DeleteColorsBillboardHandler,
		"refresh_all_billboards":       RefreshAllBillboardsHandler,

		"add_notification_channel": AddNotificationChannelHandler,
		"send_notification":        SendNotificationHandler,
	}

	completeHandlers = map[string][]AutoComplete{
		"remove_series":            {AutoComplete{Loc: 0, Choices: SeriesFullNameAutocomplete}, AutoComplete{Loc: 1, Choices: SeriesShortNameAutocomplete}},
		"change_series_title":      {AutoComplete{Loc: 0, Choices: SeriesShortNameAutocomplete}},
		"change_series_repo":       {AutoComplete{Loc: 0, Choices: SeriesShortNameAutocomplete}},
		"add_series_channel":       {AutoComplete{Loc: 0, Choices: SeriesShortNameAutocomplete}},
		"remove_job":               {AutoComplete{Loc: 0, Choices: JobShortNameAutocompleteNonG}},
		"add_series_assignment":    {AutoComplete{Loc: 1, Choices: JobShortNameAutocomplete}, AutoComplete{Loc: 2, Choices: SeriesShortNameAutocomplete}},
		"remove_series_assignment": {AutoComplete{Loc: 1, Choices: JobShortNameAutocomplete}, AutoComplete{Loc: 2, Choices: SeriesShortNameAutocomplete}},
		"series_assignments":       {AutoComplete{Loc: 0, Choices: SeriesShortNameAutocomplete}},
		"job_assignments":          {AutoComplete{Loc: 0, Choices: JobShortNameAutocomplete}},
		"create_series_billboard":  {AutoComplete{Loc: 0, Choices: SeriesShortNameAutocomplete}},
		"delete_series_billboard":  {AutoComplete{Loc: 0, Choices: SeriesShortNameAutocomplete}},
	}

	componentHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate, identifier string){
		"note_add_button":    NoteAddButtonHandler,
		"note_add_modal":     NoteAddModalHandler,
		"note_remove_button": NoteRemoveButtonHandler,
		"note_remove_modal":  NoteRemoveModalHandler,
	}
)
