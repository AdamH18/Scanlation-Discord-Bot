package bot

import (
	"log"
	"scanlation-discord-bot/database"
	"time"

	"github.com/bwmarrin/discordgo"
)

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
	rows, err := database.Repo.RemoveReminder(remID)
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
	rows, err := database.Repo.RemoveUserReminder(remID, i.Member.User.ID)
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
	rems, err := database.Repo.GetUserReminders(i.Member.User.ID)
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
	rems, err := database.Repo.GetUserReminders(userID)
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
	rems, err := database.Repo.GetAllReminders()
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

// Creates handlers for all slash commands based on relationship defined in commandHandlers
func CreateHandlers() {
	goBot.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}
