package bot

import (
	"fmt"
	"log"
	"scanlation-discord-bot/config"
	"scanlation-discord-bot/database"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func HelpHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	LogCommand(i, "help")
	helpEmbed := BuildHelpEmbed()
	RespondEmbed(s, i, helpEmbed)
}

// Handler for add_any_reminder
func AddAnyReminderHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	LogCommand(i, "add_any_reminder")
	//Compiling values into Reminder struct
	var rem database.Reminder
	rem.Guild = i.GuildID
	rem.Channel = i.ChannelID
	options := OptionsToMap(i.ApplicationCommandData().Options)
	rem.User = options["user"].UserValue(s).ID
	rem.Days = options["days"].IntValue()
	rem.Message = options["message"].StringValue()
	rem.Repeat = false
	if _, ok := options["repeat"]; ok {
		rem.Repeat = options["repeat"].BoolValue()
	}
	var mod int64
	if _, ok := options["hours-mod"]; ok {
		mod = options["hours-mod"].IntValue()
	} else {
		mod = 0
	}
	log.Printf("User: %s Days: %d Message: %s Repeat: %t Mod: %d", rem.User, rem.Days, rem.Message, rem.Repeat, mod)

	//Reminder time is user specified number of days after current time, modified by user specified hour mod
	rem.Time = (time.Now().Add(time.Hour * time.Duration(rem.Days*24)).Add(time.Hour * time.Duration(mod))).Format("2006-01-02 15:04:05")

	//Add reminder to DB
	err := database.Repo.AddReminder(rem)
	response := ""
	if err != nil {
		response = "Error adding reminder to database: " + err.Error()
	} else {
		response = "Successfully added reminder to database"
	}
	Respond(s, i, response)
}

// Handler for add_reminder
func AddReminderHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	LogCommand(i, "add_reminder")
	//Compiling values into Reminder struct
	var rem database.Reminder
	rem.Guild = i.GuildID
	rem.Channel = i.ChannelID
	rem.User = i.Member.User.ID
	options := OptionsToMap(i.ApplicationCommandData().Options)
	rem.Days = options["days"].IntValue()
	rem.Message = options["message"].StringValue()
	rem.Repeat = false
	if _, ok := options["repeat"]; ok {
		rem.Repeat = options["repeat"].BoolValue()
	}
	var mod int64
	if _, ok := options["hours-mod"]; ok {
		mod = options["hours-mod"].IntValue()
	} else {
		mod = 0
	}
	log.Printf("Days: %d Message: %s Repeat: %t Mod: %d", rem.Days, rem.Message, rem.Repeat, mod)

	//Reminder time is user specified number of days after current time, modified by user specified hour mod
	rem.Time = (time.Now().Add(time.Hour * time.Duration(rem.Days*24)).Add(time.Hour * time.Duration(mod))).Format("2006-01-02 15:04:05")

	//Add reminder to DB
	err := database.Repo.AddReminder(rem)
	response := ""
	if err != nil {
		response = "Error adding reminder to database: " + err.Error()
	} else {
		response = "Successfully added reminder to database"
	}
	Respond(s, i, response)
}

// Handler for rem_any_reminder
func RemAnyReminderHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	LogCommand(i, "rem_any_reminder")
	//Uses reminder ID to identify for removal
	options := OptionsToMap(i.ApplicationCommandData().Options)
	remID := options["id"].IntValue()
	log.Printf("ID: %d", remID)
	//Send to DB for removal
	rows, err := database.Repo.RemoveReminder(remID, i.GuildID)
	response := ""
	if err != nil {
		response = "Error removing reminder from database: " + err.Error()
	} else if rows == 0 {
		response = "Was unable to locate reminder to be removed"
	} else {
		response = "Successfully removed reminder from database"
	}
	Respond(s, i, response)
}

// Handler for rem_reminder
func RemReminderHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	LogCommand(i, "rem_reminder")
	//Uses reminder ID to identify for removal, also sends user ID to ensure user isn't removing someone else's reminder
	options := OptionsToMap(i.ApplicationCommandData().Options)
	remID := options["id"].IntValue()
	log.Printf("ID: %d", remID)
	//Send to DB for removal
	rows, err := database.Repo.RemoveUserReminder(remID, i.Member.User.ID, i.GuildID)
	response := ""
	if err != nil {
		response = "Error removing reminder from database: " + err.Error()
	} else if rows == 0 {
		response = "Was unable to locate reminder to be removed"
	} else {
		response = "Successfully removed reminder from database"
	}
	Respond(s, i, response)
}

// Handler for my_reminders
func MyRemindersHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	LogCommand(i, "my_reminders")
	//Send user ID to DB to find corresponding reminders
	rems, err := database.Repo.GetUserReminders(i.Member.User.ID, i.GuildID)
	response := ""
	if err != nil {
		response = "Error getting reminders from database: " + err.Error()
	} else {
		//Build reminders table from results
		response = BuildRemindersTable(i.Member.User.Username, rems)
	}
	Respond(s, i, response)
}

// Handler for user_reminders
func UserRemindersHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	LogCommand(i, "user_reminders")
	options := OptionsToMap(i.ApplicationCommandData().Options)
	userID := options["user"].UserValue(s).ID
	log.Printf("User: %s", userID)
	//Send user ID to DB to find corresponding reminders
	rems, err := database.Repo.GetUserReminders(userID, i.GuildID)
	response := ""
	if err != nil {
		response = "Error getting reminders from database: " + err.Error()
	} else {
		//Build reminders table from results
		response = BuildRemindersTable(options["user"].UserValue(s).Username, rems)
	}
	Respond(s, i, response)
}

// Handler for all_reminders
func AllRemindersHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	LogCommand(i, "all_reminders")
	//Send query to DB
	rems, err := database.Repo.GetAllReminders(i.GuildID)
	response := ""
	if err != nil {
		response = "Error getting reminders from database: " + err.Error()
	} else {
		//Build reminders table from results
		resp, err := BuildVerboseRemindersTable(rems)
		if err != nil {
			response = "Error creating verbose reminders table: " + err.Error()
		} else {
			response = resp
		}
	}
	Respond(s, i, response)
}

// Handler for set_any_alarm
func SetAnyAlarmHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	LogCommand(i, "set_any_alarm")
	//Compiling values into Reminder struct
	var rem database.Reminder
	rem.Guild = i.GuildID
	rem.Channel = i.ChannelID
	options := OptionsToMap(i.ApplicationCommandData().Options)
	rem.User = options["user"].UserValue(s).ID
	rem.Time = options["date-time"].StringValue()
	//Checking if user input date-time is valid
	_, err := time.Parse("2006-01-02 15:04:05", rem.Time)
	if err != nil {
		Respond(s, i, "Couldn't understand the date-time you input. Please try again while ensuring in 'YYYY-MM-DD HH:MM:SS' format")
		return
	}
	rem.Message = options["message"].StringValue()
	rem.Days = 0
	rem.Repeat = false
	if _, ok := options["days"]; ok {
		rem.Repeat = true
		rem.Days = options["days"].IntValue()
	}
	log.Printf("User: %s Date-Time: %s Message: %s Days: %d", rem.User, rem.Time, rem.Message, rem.Days)

	//Adding reminder to DB
	err = database.Repo.AddReminder(rem)
	response := ""
	if err != nil {
		response = "Error adding reminder to database: " + err.Error()
	} else {
		response = "Successfully added reminder to database"
	}
	Respond(s, i, response)
}

// Handler for set_alarm
func SetAlarmHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	LogCommand(i, "set_alarm")
	//Compiling values into Reminder struct
	var rem database.Reminder
	rem.Guild = i.GuildID
	rem.Channel = i.ChannelID
	rem.User = i.Member.User.ID
	options := OptionsToMap(i.ApplicationCommandData().Options)
	rem.Time = options["date-time"].StringValue()
	//Checking if user input date-time is valid
	_, err := time.Parse("2006-01-02 15:04:05", rem.Time)
	if err != nil {
		Respond(s, i, "Couldn't understand the date-time you input. Please try again while ensuring in 'YYYY-MM-DD HH:MM:SS' format")
		return
	}
	rem.Message = options["message"].StringValue()
	rem.Days = 0
	rem.Repeat = false
	if _, ok := options["days"]; ok {
		rem.Repeat = true
		rem.Days = options["days"].IntValue()
	}
	log.Printf("Date-Time: %s Message: %s Days: %d", rem.Time, rem.Message, rem.Days)

	//Adding reminder to DB
	err = database.Repo.AddReminder(rem)
	response := ""
	if err != nil {
		response = "Error adding reminder to database: " + err.Error()
	} else {
		response = "Successfully added reminder to database"
	}
	Respond(s, i, response)
}

// Handler for add_series
func AddSeriesHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	LogCommand(i, "add_series")
	//Compiling values into Series struct
	var ser database.Series
	ser.Guild = i.GuildID
	options := OptionsToMap(i.ApplicationCommandData().Options)
	ser.NameSh = options["short-name"].StringValue()
	ser.NameFull = options["full-name"].StringValue()
	ser.PingRole = ""
	if _, ok := options["ping-role"]; ok {
		ser.PingRole = options["ping-role"].RoleValue(s, i.GuildID).ID
	}
	ser.RepoLink = ""
	if _, ok := options["repo-link"]; ok {
		ser.RepoLink = options["repo-link"].StringValue()
	}
	fullCreate := options["full-create"].BoolValue()
	log.Printf("Full-Name: %s Short-Name: %s Full-Create: %t Ping-Role: %s Repo-Link: %s", ser.NameFull, ser.NameSh, fullCreate, ser.PingRole, ser.RepoLink)

	//Adding series to DB
	err := database.Repo.AddSeries(ser)
	response := ""
	if err != nil {
		response = "Error adding series to database: " + err.Error()
	} else {
		response = "Successfully created series " + ser.NameFull
	}

	// Creating channels and roles for full create plus keeping track of results
	channelRes := ""
	roleRes := ""
	if fullCreate {
		err := CreateChannels(ser)
		if err != nil {
			channelRes = "Error creating channels: " + err.Error()
		} else {
			channelRes = "Successfully created channels"
		}
		if ser.PingRole == "" {
			ser.PingRole, err = CreatePingRole(ser)
			if err != nil {
				roleRes = "Error creating ping role: " + err.Error()
			} else {
				roleRes = "Successfully created ping role"
			}
		}
	}

	if channelRes != "" {
		response = response + "\n" + channelRes
	}
	if roleRes != "" {
		response = response + "\n" + roleRes
	}
	RespondNonEph(s, i, response)
}

// Handler for remove_series
func RemoveSeriesHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	LogCommand(i, "remove_series")
	options := OptionsToMap(i.ApplicationCommandData().Options)
	nameSh := options["short-name"].StringValue()
	nameFull := options["full-name"].StringValue()
	log.Printf("Full-Name: %s Short-Name: %s", nameFull, nameSh)

	//Removing series from DB
	done, err := database.Repo.RemoveSeries(nameSh, nameFull, i.GuildID)
	response := ""
	if err != nil {
		response = "Error removing series from database: " + err.Error()
	} else if !done {
		response = "Could not locate series for removal"
	} else {
		response = "Successfully removed series and all references from databases"
	}
	Respond(s, i, response)
}

// Handler for server_series
func ServerSeriesHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	LogCommand(i, "server_series")

	embed, err := BuildServerSeriesEmbed(i.GuildID)
	if err != nil {
		Respond(s, i, "Failed to build embed: "+err.Error())
		return
	}
	RespondEmbed(s, i, embed)
}

// Handler for change_series_title
func ChangeSeriesTitleHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	LogCommand(i, "change_series_title")
	options := OptionsToMap(i.ApplicationCommandData().Options)
	nameSh := options["short-name"].StringValue()
	nameFull := options["new-full-name"].StringValue()
	log.Printf("Short-Name: %s New-Full-Name: %s", nameSh, nameFull)

	//Updating title
	done, err := database.Repo.UpdateSeriesName(nameSh, nameFull, i.GuildID)
	response := ""
	if err != nil {
		response = "Error changing series name: " + err.Error()
	} else if !done {
		response = "Could not locate series for name change"
	} else {
		response = "Successfully changed series name"
	}
	Respond(s, i, response)
}

// Handler for change_series_repo
func ChangeSeriesRepoHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	LogCommand(i, "change_series_repo")
	options := OptionsToMap(i.ApplicationCommandData().Options)
	nameSh := options["short-name"].StringValue()
	repoLink := options["new-repo-link"].StringValue()
	log.Printf("Short-Name: %s New-Repo-Link: %s", nameSh, repoLink)

	//Updating link
	done, err := database.Repo.UpdateSeriesRepoLink(nameSh, repoLink, i.GuildID)
	response := ""
	if err != nil {
		response = "Error changing series link: " + err.Error()
	} else if !done {
		response = "Could not locate series for link change"
	} else {
		response = "Successfully changed series link"
	}
	Respond(s, i, response)
}

// Handler for add_series_channel
func AddSeriesChannelHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	LogCommand(i, "add_series_channel")
	//Compiling values into Channel struct
	var cha database.Channel
	cha.Guild = i.GuildID
	cha.Channel = i.ChannelID
	options := OptionsToMap(i.ApplicationCommandData().Options)
	cha.Series = options["series"].StringValue()
	log.Printf("Series: %s", cha.Series)

	//Check if series to be added to exists
	if !database.Repo.RegisteredSeries(cha.Series, cha.Guild) {
		Respond(s, i, "Could not find series in database. Did not register channel")
		return
	}

	//Add channel to DB
	err := database.Repo.AddChannel(cha)
	response := ""
	if err != nil {
		response = "Error adding channel to database: " + err.Error()
	} else {
		response = "Successfully added channel to database"
	}
	Respond(s, i, response)
}

// Handler for remove_series_channel
func RemoveSeriesChannelHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	LogCommand(i, "remove_series_channel")

	//Removing channel from DB
	done, err := database.Repo.RemoveChannel(i.ChannelID)
	response := ""
	if err != nil {
		response = "Error removing channel from database: " + err.Error()
	} else if !done {
		response = "This channel was not registered in the first place"
	} else {
		response = "Successfully removed channel from database"
	}
	Respond(s, i, response)
}

// Handler for add_user
func AddUserHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	LogCommand(i, "add_user")
	//Compiling values into Channel struct
	var usr database.User
	usr.Guild = i.GuildID
	usr.Color = ""
	usr.VanityRole = ""
	options := OptionsToMap(i.ApplicationCommandData().Options)
	usr.User = options["user"].UserValue(s).ID
	log.Printf("User: %s", usr.User)

	//Add user to DB
	err := database.Repo.AddUser(usr)
	response := ""
	if err != nil {
		response = "Error adding user to database: " + err.Error()
	} else {
		response = "Successfully added user to database"

		//If member role is set, give role to new user
		mem := database.Repo.GetMemberRole(i.GuildID)
		if mem != "" {
			err = s.GuildMemberRoleAdd(i.GuildID, usr.User, mem)
			if err != nil {
				response += "\nError giving member role: " + err.Error()
			} else {
				response += "\nMember role successfully given"
			}
		}
	}
	Respond(s, i, response)
}

// Handler for remove_user
func RemoveUserHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	LogCommand(i, "remove_user")
	options := OptionsToMap(i.ApplicationCommandData().Options)
	user := options["user"].UserValue(s).ID
	log.Printf("User: %s", user)

	//Removing user from DB
	done, err := database.Repo.RemoveUser(user, i.GuildID)
	response := ""
	if err != nil {
		response = "Error removing user from database: " + err.Error()
	} else if !done {
		response = "This user was not registered in the first place"
	} else {
		response = "Successfully removed user from database"

		//If member role is set, remove role from user
		mem := database.Repo.GetMemberRole(i.GuildID)
		if mem != "" {
			err = s.GuildMemberRoleRemove(i.GuildID, user, mem)
			if err != nil {
				response += "\nError removing member role: " + err.Error()
			} else {
				response += "\nMember role successfully removed"
			}
		}
	}
	Respond(s, i, response)
}

// Handler for server_users
func ServerUsersHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	LogCommand(i, "server_users")

	embed, err := BuildServerUsersEmbed(i.GuildID)
	if err != nil {
		Respond(s, i, "Failed to build embed: "+err.Error())
		return
	}
	RespondEmbed(s, i, embed)
}

// Handler for add_job
func AddJobHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	LogCommand(i, "add_job")
	//Compiling values into Job struct
	var job database.Job
	job.Guild = i.GuildID
	options := OptionsToMap(i.ApplicationCommandData().Options)
	job.JobFull = options["name-full"].StringValue()
	job.JobSh = options["name-short"].StringValue()
	log.Printf("Name-Full: %s Name-Short: %s", job.JobFull, job.JobSh)

	//Add job to DB
	err := database.Repo.AddJob(job)
	response := ""
	if err != nil {
		response = "Error adding job to database: " + err.Error()
	} else {
		response = "Successfully added job to database"
	}
	Respond(s, i, response)
}

// Handler for add_global_job
func AddGlobalJobHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	LogCommand(i, "add_global_job")
	if i.Member.User.ID != config.Owner {
		Respond(s, i, "Owner only command. Please use add_job instead")
		return
	}
	//Compiling values into Job struct
	var job database.Job
	job.Guild = "GLOBAL"
	options := OptionsToMap(i.ApplicationCommandData().Options)
	job.JobFull = options["name-full"].StringValue()
	job.JobSh = options["name-short"].StringValue()
	log.Printf("Name-Full: %s Name-Short: %s", job.JobFull, job.JobSh)

	//Add job to DB
	err := database.Repo.AddJob(job)
	response := ""
	if err != nil {
		response = "Error adding job to database: " + err.Error()
	} else {
		response = "Successfully added job to database"
	}
	Respond(s, i, response)
}

// Handler for remove_job
func RemoveJobHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	LogCommand(i, "remove_job")
	options := OptionsToMap(i.ApplicationCommandData().Options)
	nameSh := options["name-short"].StringValue()
	log.Printf("Name-Short: %s", nameSh)

	//Removing job from DB
	done, err := database.Repo.RemoveJob(nameSh, i.GuildID)
	response := ""
	if err != nil {
		response = "Error removing job from database: " + err.Error()
	} else if !done {
		response = "Could not locate job for removal"
	} else {
		response = "Successfully removed job and all references from databases"
	}
	Respond(s, i, response)
}

// Handler for server_jobs
func ServerJobsHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	LogCommand(i, "server_jobs")

	embed, err := BuildServerJobsEmbed(i.GuildID)
	if err != nil {
		Respond(s, i, "Failed to build embed: "+err.Error())
		return
	}
	RespondEmbed(s, i, embed)
}

// Handler for add_member_role
func AddMemberRoleHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	LogCommand(i, "add_member_role")
	//Compiling values into MemberRole struct
	var mem database.MemberRole
	mem.Guild = i.GuildID
	options := OptionsToMap(i.ApplicationCommandData().Options)
	mem.Role = options["role"].RoleValue(s, mem.Guild).ID
	log.Printf("Role: %s", mem.Role)

	//Add role to DB
	err := database.Repo.AddMemberRole(mem)
	response := ""
	if err != nil {
		response = "Error adding member role to database: " + err.Error()
	} else {
		response = "Successfully added member role to database"
	}
	Respond(s, i, response)
}

// Handler for remove_member_role
func RemoveMemberRoleHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	LogCommand(i, "remove_member_role")

	//Removing member role from DB
	done, err := database.Repo.RemoveMemberRole(i.GuildID)
	response := ""
	if err != nil {
		response = "Error removing member role from database: " + err.Error()
	} else if !done {
		response = "No role was registered in the first place"
	} else {
		response = "Successfully removed member role from database"
	}
	Respond(s, i, response)
}

// Handler for reg_series_channels
func RegSeriesChannelsHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	LogCommand(i, "reg_series_channels")
	//Compiling values into SeriesChannels struct
	var sec database.SeriesChannels
	sec.Guild = i.GuildID
	options := OptionsToMap(i.ApplicationCommandData().Options)
	sec.Top = options["top"].StringValue()
	sec.Bottom = options["bottom"].StringValue()
	log.Printf("Top: %s Bottom: %s", sec.Top, sec.Bottom)

	//Add role to DB
	err := database.Repo.AddSeriesChannels(sec)
	response := ""
	if err != nil {
		response = "Error adding series channels to database: " + err.Error()
	} else {
		response = "Successfully added series channels to database"
	}
	Respond(s, i, response)
}

// Handler for add_series_assignment
func AddSeriesAssignmentHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	LogCommand(i, "add_series_assignment")
	//Compiling values into SeriesAssignment struct
	var sea database.SeriesAssignment
	sea.Guild = i.GuildID
	options := OptionsToMap(i.ApplicationCommandData().Options)
	sea.User = options["user"].UserValue(s).ID
	sea.Job = options["job"].StringValue()
	sea.Series = ""
	if _, ok := options["series"]; ok {
		sea.Series = options["series"].StringValue()
	}
	log.Printf("User: %s Job: %s Series: %s", sea.User, sea.Job, sea.Series)
	var err error
	if sea.Series == "" {
		sea.Series, err = database.Repo.GetLocalSeries(i.ChannelID)
		if err != nil {
			Respond(s, i, "Unable to locate series for this command: "+err.Error())
			return
		}
	}

	//Check if user, series, and job exist
	if !database.Repo.RegisteredUser(sea.User, sea.Guild) {
		Respond(s, i, "Could not find user in database. Did not register assignment")
		return
	}
	if !database.Repo.RegisteredSeries(sea.Series, sea.Guild) {
		Respond(s, i, "Could not find series in database. Did not register assignment")
		return
	}
	if !database.Repo.RegisteredJob(sea.Job, sea.Guild) {
		Respond(s, i, "Could not find job in database. Did not register assignment")
		return
	}

	//Add assignment to DB
	err = database.Repo.AddAssignment(sea)
	response := ""
	if err != nil {
		response = "Error adding assignment to database: " + err.Error()
	} else {
		response = "Successfully added assignment to database"
	}
	Respond(s, i, response)
}

// Handler for remove_series_assignment
func RemoveSeriesAssignmentHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	LogCommand(i, "remove_series_assignment")
	options := OptionsToMap(i.ApplicationCommandData().Options)
	user := options["user"].UserValue(s).ID
	job := options["job"].StringValue()
	series := ""
	if _, ok := options["series"]; ok {
		series = options["series"].StringValue()
	}
	log.Printf("User: %s Series: %s Job: %s", user, series, job)
	var err error
	if series == "" {
		series, err = database.Repo.GetLocalSeries(i.ChannelID)
		if err != nil {
			Respond(s, i, "Unable to locate series for this command: "+err.Error())
			return
		}
	}

	//Removing assignemnt from DB
	done, err := database.Repo.RemoveSeriesAssignment(user, series, job, i.GuildID)
	response := ""
	if err != nil {
		response = "Error removing assignment from database: " + err.Error()
	} else if !done {
		response = "Could not locate assignment for removal"
	} else {
		response = "Successfully removed assignment from database"
	}
	Respond(s, i, response)
}

// Handler for remove_all_assignments
func RemoveAllAssignmentsHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	LogCommand(i, "remove_all_assignments")
	options := OptionsToMap(i.ApplicationCommandData().Options)
	user := options["user"].UserValue(s).ID
	log.Printf("User: %s", user)

	//Removing assignemnt from DB
	done, err := database.Repo.RemoveAllAssignments(user, i.GuildID)
	response := ""
	if err != nil {
		response = "Error removing assignments from database: " + err.Error()
	} else if !done {
		response = "User had no assignments to remove"
	} else {
		response = "Successfully removed assignments from database"
	}
	Respond(s, i, response)
}

// Handler for series_assignments
func SeriesAssignmentsHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	LogCommand(i, "series_assignments")
	options := OptionsToMap(i.ApplicationCommandData().Options)
	series := ""
	var err error
	if _, ok := options["series"]; ok {
		series = options["series"].StringValue()
	}
	log.Printf("Series: %s", series)
	if series == "" {
		series, err = database.Repo.GetLocalSeries(i.ChannelID)
		if err != nil {
			Respond(s, i, "Unable to locate series for this command: "+err.Error())
			return
		}
	}

	//Get all series assignments
	assMap, err := database.Repo.GetSeriesAssignments(series, i.GuildID)
	if err != nil {
		Respond(s, i, "Unable to get series assignments: "+err.Error())
		return
	}
	embed, err := BuildSeriesAssignmentsEmbed(assMap, series, i.GuildID)
	if err != nil {
		Respond(s, i, "Failed to build embed: "+err.Error())
		return
	}
	RespondEmbed(s, i, embed)
}

// Handler for my_assignments
func MyAssignmentsHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	LogCommand(i, "my_assignments")

	//Get all user assignments
	assMap, err := database.Repo.GetUserAssignments(i.Member.User.ID, i.GuildID)
	if err != nil {
		Respond(s, i, "Unable to get user assignments: "+err.Error())
		return
	}
	embed, err := BuildUserAssignmentsEmbed(assMap, i.Member.User.ID, i.GuildID)
	if err != nil {
		Respond(s, i, "Failed to build embed: "+err.Error())
		return
	}
	RespondEmbed(s, i, embed)
}

// Handler for user_assignments
func UserAssignmentsHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	LogCommand(i, "user_assignments")
	options := OptionsToMap(i.ApplicationCommandData().Options)
	user := options["user"].UserValue(s).ID
	log.Printf("User: %s", user)

	//Get all user assignments
	assMap, err := database.Repo.GetUserAssignments(user, i.GuildID)
	if err != nil {
		Respond(s, i, "Unable to get user assignments: "+err.Error())
		return
	}
	embed, err := BuildUserAssignmentsEmbed(assMap, user, i.GuildID)
	if err != nil {
		Respond(s, i, "Failed to build embed: "+err.Error())
		return
	}
	RespondEmbed(s, i, embed)
}

// Handler for job_assignments
func JobAssignmentsHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	LogCommand(i, "job_assignments")
	options := OptionsToMap(i.ApplicationCommandData().Options)
	job := options["job"].StringValue()
	log.Printf("Job: %s", job)

	//Get all job assignments
	assMap, err := database.Repo.GetJobAssignments(job, i.GuildID)
	if err != nil {
		Respond(s, i, "Unable to get user assignments: "+err.Error())
		return
	} else if len(assMap) == 0 {
		Respond(s, i, "No users found assigned to that job")
		return
	}
	embed, err := BuildJobAssignmentsEmbed(assMap, job, i.GuildID)
	if err != nil {
		Respond(s, i, "Failed to build embed: "+err.Error())
		return
	}
	RespondEmbed(s, i, embed)
}

// Pings the given job with a provided message
func JobPinger(job string, s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := OptionsToMap(i.ApplicationCommandData().Options)
	message := ""
	if _, ok := options["message"]; ok {
		message = options["message"].StringValue()
	}
	log.Printf("Message: %s", message)
	series, err := database.Repo.GetLocalSeries(i.ChannelID)
	if err != nil {
		Respond(s, i, "Unable to locate series for this command: "+err.Error())
		return
	}

	//Get all user assignments
	users, err := database.Repo.GetSeriesJobAssignments(series, job, i.GuildID)
	if err != nil {
		Respond(s, i, "Unable to get users: "+err.Error())
		return
	}
	if len(users) == 0 {
		Respond(s, i, "No users assigned to this job for this series")
		return
	}

	//Build message
	response := ""
	for _, user := range users {
		ping, err := GetUserPing(i.GuildID, user)
		if err != nil {
			Respond(s, i, "Error building ping: "+err.Error())
			return
		}
		response += ping
	}
	if message != "" {
		response += " " + message
	}
	RespondNonEph(s, i, response)
}

// Handler for tl
func TLPingHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	LogCommand(i, "tl")
	JobPinger("tl", s, i)
}

// Handler for rd
func RDPingHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	LogCommand(i, "rd")
	JobPinger("rd", s, i)
}

// Handler for ts
func TSPingHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	LogCommand(i, "ts")
	JobPinger("ts", s, i)
}

// Handler for pr
func PRPingHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	LogCommand(i, "pr")
	JobPinger("pr", s, i)
}

// Handler for my_settings
func MySettingsHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	LogCommand(i, "my_settings")

	embed, err := BuildUserSettingsEmbed(i.Member.User.ID, i.GuildID)
	if err != nil {
		Respond(s, i, "Failed to build embed: "+err.Error())
		return
	}
	RespondEmbed(s, i, embed)
}

// Handler for user_settings
func UserSettingsHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	LogCommand(i, "user_settings")
	options := OptionsToMap(i.ApplicationCommandData().Options)
	user := options["user"].UserValue(s).ID
	log.Printf("User: %s", user)

	embed, err := BuildUserSettingsEmbed(user, i.GuildID)
	if err != nil {
		Respond(s, i, "Failed to build embed: "+err.Error())
		return
	}
	RespondEmbed(s, i, embed)
}

// Handler for set_color
func SetColorHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	LogCommand(i, "set_color")
	options := OptionsToMap(i.ApplicationCommandData().Options)
	color := options["color"].StringValue()
	user := i.Member.User.ID
	log.Printf("Color: %s", color)

	//Update color in DB
	err := database.Repo.UpdateColor(user, color, i.GuildID)
	response := ""
	if err != nil {
		response = "Error updating color: " + err.Error()
	} else {
		response = "Successfully updated your credits color"
	}
	Respond(s, i, response)
}

// Handler for set_user_color
func SetUserColorHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	LogCommand(i, "set_user_color")
	options := OptionsToMap(i.ApplicationCommandData().Options)
	color := options["color"].StringValue()
	user := options["user"].UserValue(s).ID
	log.Printf("User: %s Color: %s", user, color)

	//Update color in DB
	err := database.Repo.UpdateColor(user, color, i.GuildID)
	response := ""
	if err != nil {
		response = "Error updating color: " + err.Error()
	} else {
		response = "Successfully updated user's credits color"
	}
	Respond(s, i, response)
}

// Handler for vanity_role
func VanityRoleHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	LogCommand(i, "vanity_role")
	options := OptionsToMap(i.ApplicationCommandData().Options)
	name := options["name"].StringValue()
	color := options["color"].StringValue()
	user := ""
	if _, ok := options["copy-user"]; ok {
		user = options["copy-user"].UserValue(s).ID
	}
	log.Printf("Name: %s Color: %s Copy-User: %s", name, color, user)

	//If user was set, just give same role (if specified user has one)
	var err error
	if user != "" {
		roleId, err := database.Repo.GetUserVanityRole(user, i.GuildID)
		response := ""
		if err != nil {
			response = "Error getting role from user: " + err.Error()
		} else if roleId == "" {
			response = "User does not have a vanity role"
		} else {
			err = s.GuildMemberRoleAdd(i.GuildID, i.Member.User.ID, roleId)
			if err != nil {
				response = "Error adding vanity role:" + err.Error()
			} else {
				response = "Successfully added vanity role"
				err = database.Repo.UpdateVanityRole(i.Member.User.ID, roleId, i.GuildID)
				if err != nil {
					log.Println("Error recording role assignment in database: " + err.Error())
				}
			}
		}
		Respond(s, i, response)
		return
	}

	//Now that it's determined the field will be used, parse color
	color = strings.TrimSpace(color)
	if len(color) != 6 {
		Respond(s, i, "Color is not right length")
		return
	}
	colorInt64, err := strconv.ParseInt(color, 16, 64)
	if err != nil {
		Respond(s, i, "Error parsing color: "+err.Error())
		return
	}
	colorInt := int(colorInt64)

	//Check if existing vanity role can just be updated
	usrRole, err := database.Repo.GetUserVanityRole(i.Member.User.ID, i.GuildID)
	if err != nil {
		Respond(s, i, "Error getting roles from database: "+err.Error())
		return
	}
	if usrRole != "" {
		num, err := database.Repo.NumUsersWithVanity(usrRole, i.GuildID)
		if err != nil {
			Respond(s, i, "Error getting roles count from database: "+err.Error())
			return
		}
		//If user running command is only one with role, edit it
		if num == 1 {
			falseVar := false
			role := discordgo.RoleParams{
				Name:        name,
				Color:       &colorInt,
				Hoist:       &falseVar,
				Mentionable: &falseVar,
			}
			_, err = s.GuildRoleEdit(i.GuildID, usrRole, &role)
			response := ""
			if err != nil {
				log.Println("Error editing role: " + err.Error() + " - Creating new role instead")
			} else {
				response = "Successfully edited role to new specs"
				Respond(s, i, response)
				return
			}
		}
	}

	//If new role needs to be created, make sure server has space for it
	roles, err := s.GuildRoles(i.GuildID)
	if err != nil {
		Respond(s, i, "Unable to check server roles: "+err.Error())
		return
	} else if len(roles) > 240 {
		Respond(s, i, "Too few role slots left in server to make new role")
		return
	}

	//If making new role and user has a current role, remove it. Ignore if fails
	if usrRole != "" {
		err = s.GuildMemberRoleRemove(i.GuildID, i.Member.User.ID, usrRole)
		if err != nil {
			log.Println("Error removing old vanity role: " + err.Error())
		}
	}

	//Make and give new role
	falseVar := false
	roleP := discordgo.RoleParams{
		Name:        name,
		Color:       &colorInt,
		Hoist:       &falseVar,
		Mentionable: &falseVar,
	}
	role, err := s.GuildRoleCreate(i.GuildID, &roleP)
	if err != nil {
		Respond(s, i, "Error creating new role: "+err.Error())
		return
	}
	err = s.GuildMemberRoleAdd(i.GuildID, i.Member.User.ID, role.ID)
	if err != nil {
		Respond(s, i, "Error giving new role: "+err.Error())
		return
	}
	//Register new role in database
	err = database.Repo.UpdateVanityRole(i.Member.User.ID, role.ID, i.GuildID)
	if err != nil {
		Respond(s, i, "Error adding new role to database: "+err.Error())
		return
	}
	Respond(s, i, "Successfully created and gave new vanity role")
}

// Handler for rem_vanity_role
func RemVanityRoleHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	LogCommand(i, "rem_vanity_role")

	usrRole, err := database.Repo.GetUserVanityRole(i.Member.User.ID, i.GuildID)
	if err != nil {
		Respond(s, i, "Error getting roles from database: "+err.Error())
		return
	} else if usrRole == "" {
		Respond(s, i, "No vanity role to remove")
		return
	}

	//Remove role in Discord
	err = s.GuildMemberRoleRemove(i.GuildID, i.Member.User.ID, usrRole)
	if err != nil {
		Respond(s, i, "Error removing role: "+err.Error())
		return
	}
	//Update database
	err = database.Repo.UpdateVanityRole(i.Member.User.ID, "", i.GuildID)
	response := ""
	if err != nil {
		response = "Error updating database: " + err.Error()
	} else {
		response = "Successfully removed role"
	}
	Respond(s, i, response)
}

// Handler for create_assignments_billboard
func CreateAssignmentsBillboardHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	LogCommand(i, "create_assignments_billboard")

	bill, _, err := database.Repo.GetRolesBillboard(i.GuildID)
	if err != nil {
		Respond(s, i, "Error getting existing billboard info: "+err.Error())
		return
	} else if bill != "" {
		Respond(s, i, "Server already has an assignments billboard. Please remove the existing one first")
		return
	}

	//Billboard should be created, so gather data
	assMap, err := database.Repo.GetAllAssignments(i.GuildID)
	if err != nil {
		Respond(s, i, "Error getting server assignments info: "+err.Error())
		return
	}

	embed, err := BuildFullAssignmentsEmbed(assMap, i.GuildID)
	if err != nil {
		Respond(s, i, "Error building embed: "+err.Error())
		return
	}

	message := discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{embed},
	}
	msg, err := s.ChannelMessageSendComplex(i.ChannelID, &message)
	if err != nil {
		Respond(s, i, "Error sending message: "+err.Error())
		return
	}

	bb := database.JobBB{
		Guild:   i.GuildID,
		Channel: i.ChannelID,
		Message: msg.ID,
	}
	err = database.Repo.AddRolesBillboard(bb)
	response := ""
	if err != nil {
		response = "Error updating database: " + err.Error()
	} else {
		response = "Successfully created assignments billboard"
	}
	Respond(s, i, response)
}

// Handler for delete_assignments_billboard
func DeleteAssignmentsBillboardHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	LogCommand(i, "delete_assignments_billboard")

	//Removing roles billboard from DB
	done, err := database.Repo.RemoveRolesBillboard(i.GuildID)
	response := ""
	if err != nil {
		response = "Error removing billboard from database: " + err.Error()
	} else if !done {
		response = "Could not locate billboard for removal"
	} else {
		response = "Successfully removed billboard from database"
	}
	Respond(s, i, response)
}

// Handler for create_colors_billboard
func CreateColorsBillboardHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	LogCommand(i, "create_colors_billboard")

	bill, _, err := database.Repo.GetColorsBillboard(i.GuildID)
	if err != nil {
		Respond(s, i, "Error getting existing billboard info: "+err.Error())
		return
	} else if bill != "" {
		Respond(s, i, "Server already has a colors billboard. Please remove the existing one first")
		return
	}

	//Billboard should be created, so gather data
	assMap, err := database.Repo.GetAllColors(i.GuildID)
	if err != nil {
		Respond(s, i, "Error getting user color info: "+err.Error())
		return
	}

	embed, err := BuildColorsEmbed(assMap, i.GuildID)
	if err != nil {
		Respond(s, i, "Error building embed: "+err.Error())
		return
	}

	message := discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{embed},
	}
	msg, err := s.ChannelMessageSendComplex(i.ChannelID, &message)
	if err != nil {
		Respond(s, i, "Error sending message: "+err.Error())
		return
	}

	bb := database.ColorBB{
		Guild:   i.GuildID,
		Channel: i.ChannelID,
		Message: msg.ID,
	}
	err = database.Repo.AddColorsBillboard(bb)
	response := ""
	if err != nil {
		response = "Error updating database: " + err.Error()
	} else {
		response = "Successfully created colors billboard"
	}
	Respond(s, i, response)
}

// Handler for delete_colors_billboard
func DeleteColorsBillboardHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	LogCommand(i, "delete_colors_billboard")

	//Removing roles billboard from DB
	done, err := database.Repo.RemoveColorsBillboard(i.GuildID)
	response := ""
	if err != nil {
		response = "Error removing billboard from database: " + err.Error()
	} else if !done {
		response = "Could not locate billboard for removal"
	} else {
		response = "Successfully removed billboard from database"
	}
	Respond(s, i, response)
}

// Handler for refresh_all_billboards
func RefreshAllBillboardsHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	LogCommand(i, "refresh_all_billboards")

	go UpdateAllSeriesBillboards(i.GuildID)
	go UpdateAssignmentsBillboard(i.GuildID)
	go UpdateColorsBillboard(i.GuildID)

	Respond(s, i, "Update started for all billboards")
}

// Handler for add_notification_channel
func AddNotificationChannelHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	LogCommand(i, "add_notification_channel")
	//Compiling values into Channel struct
	var cha database.NotificationChannel
	cha.Guild = i.GuildID
	cha.Channel = i.ChannelID

	//Add channel to DB
	err := database.Repo.AddNotificationChannel(cha)
	response := ""
	if err != nil {
		response = "Error adding notification channel to database: " + err.Error()
	} else {
		response = "Successfully added notification channel to database"
	}
	Respond(s, i, response)
}

// Handler for send_notification
func SendNotificationHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	LogCommand(i, "send_notification")
	if i.Member.User.ID != config.Owner {
		Respond(s, i, "Owner only command")
		return
	}
	options := OptionsToMap(i.ApplicationCommandData().Options)
	message := options["message"].StringValue()
	log.Printf("Message: %s", message)

	//Get all channels to send message to
	channels, err := database.Repo.GetAllNotificationChannels()
	if err != nil {
		Respond(s, i, "Error getting channels to send to: "+err.Error())
		return
	}

	for i := 0; i < len(message); i++ {
		if message[i] != '\\' {
			continue
		}
		if i == len(message)-1 || (message[i+1] != 'n' && message[i+1] != '\\') {
			continue
		}
		if message[i+1] == 'n' {
			after := ""
			if len(message) > i+1 {
				after = message[i+2:]
			}
			message = message[:i] + "\n" + after
		} else if message[i+1] == '\\' {
			after := ""
			if len(message) > i+1 {
				after = message[i+2:]
			}
			message = message[:i] + "\\" + after
		}
	}

	good := 0
	bad := 0
	for _, channel := range channels {
		_, err = s.ChannelMessageSend(channel.Channel, message)
		if err != nil {
			bad++
			log.Printf("Message failed to send to server %s in channel %s: %s", channel.Guild, channel.Channel, err.Error())
		} else {
			good++
		}
	}
	Respond(s, i, fmt.Sprintf("%d messages sent successfully\n%d messages failed to send", good, bad))
}

// Creates handlers for all slash commands based on relationship defined in commandHandlers
func CreateHandlers() {
	goBot.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}
