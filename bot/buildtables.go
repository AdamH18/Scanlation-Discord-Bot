package bot

import (
	"fmt"
	"log"
	"scanlation-discord-bot/database"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/bwmarrin/discordgo"
)

// Returns table with all values in reminders DB included
func BuildVerboseRemindersTable(rems []database.Reminder) (string, error) {
	var buf strings.Builder
	w := tabwriter.NewWriter(&buf, 1, 1, 1, ' ', 0)
	fmt.Fprintln(w, "Reminders for all users:\n\nID\tGuild\tChannel\tUser\tDays\tMessage\tRepeat\tTime\t")
	for _, rem := range rems {
		usr, err := GetUserName(rem.Guild, rem.User)
		if err != nil {
			return "", err
		}
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%d\t%s\t%t\t%s\t\n", rem.ID, rem.Guild, rem.Channel, usr, rem.Days, rem.Message, rem.Repeat, rem.Time)
	}
	w.Flush()
	return "```\n" + buf.String() + "```", nil
}

// Returns table with useful values in reminders DB included
func BuildRemindersTable(user string, rems []database.Reminder) string {
	var buf strings.Builder
	w := tabwriter.NewWriter(&buf, 1, 1, 2, ' ', 0)
	fmt.Fprintf(w, "Reminders for user %s:\n\nID\tMessage\tReminder Time\tDays\tRepeat\t\n", user)
	for _, rem := range rems {
		fmt.Fprintf(w, "%d\t%s\t%s\t%d\t%t\t\n", rem.ID, rem.Message, rem.Time, rem.Days, rem.Repeat)
	}
	w.Flush()
	return "```\n" + buf.String() + "```"
}

func BuildHelpEmbed() *discordgo.MessageEmbed {
	//Initialize embed
	embed := discordgo.MessageEmbed{
		Type:        discordgo.EmbedTypeRich,
		Title:       "Help",
		Description: "The two main functions this bot currently performs are sending reminders and tracking assignments to series.",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name: "Reminder commands:",
				Value: "Reminders to be sent out are checked hourly, so don't treat these like an exact alarm.\n\n" +
					"**add_reminder** - Basic command to add a reminder. Times are relative to time command is used\n" +
					"**set_alarm** - Set an alarm for a non-relative time. Times are evaluated in GMT\n" +
					"**my_reminders** - Check your existing reminders. Use this to get ID for removal\n" +
					"**rem_reminder** - Delete a reminder that is no longer needed. Use my_reminders to find ID",
			},
			{
				Name: "Assignment commands:",
				Value: "Many of these commands are contextual commands. That means that if you use it in a channel under a series category, it will automatically know you want that series.\n\n" +
					"**series_assignments** - Check current assignments for a series and credit colors for all users\n" +
					"**my_assignments** - Check your personal assignments\n" +
					"**user_assignments** - Check the assignments of a given user\n" +
					"**job_assignments** - Find everyone who holds a particular job\n" +
					"**tl**, **rd**, **ts**, **pr** - Ping the people assigned to this role and attach a message if so desired. To include a message, just press space after the /ts command and it will automatically select the optional message field",
			},
			{
				Name: "User customization:",
				Value: "**set_color** - Set the color you wish to have your name in on credits pages\n" +
					"**vanity_role** - Give yourself a vanity role with a custom name and color. If the server is getting close to maximum roles, you can copy someone else's role instead\n" +
					"**rem_vanity_role** - Get rid of your vanity role. Never actually deletes it from the server\n" +
					"**my_settings** - Gives you an overview of all your settings in this server",
			},
			{
				Name:  "Misc:",
				Value: "**check_db** - Check the current status of the database",
			},
		},
	}
	return &embed
}

// Standardized order based on reasonable order for scanlation roles
func less(a, b string) bool {
	order := map[string]int{"tl": 0, "tlc": 1, "cl": 2, "rd": 3, "ts": 4, "pr": 5, "qc": 6}
	a = strings.ToLower(a)
	b = strings.ToLower(b)

	_, oka := order[a]
	_, okb := order[b]
	if oka && okb {
		return order[a] < order[b]
	}
	if oka {
		return true
	}
	if okb {
		return false
	}
	return a < b
}

// Builds the embed for showing series assignments
func BuildSeriesAssignmentsEmbed(assMap map[string][]string, series string, guild string) (*discordgo.MessageEmbed, error) {
	//Initialize embed
	embed := discordgo.MessageEmbed{
		Type: discordgo.EmbedTypeRich,
	}
	//Make title using full series name
	fullName, err := database.Repo.GetSeriesFullName(series, guild)
	if err != nil {
		return nil, err
	}
	embed.Title = "Assignments for " + fullName

	//Create sorted job slice
	keys := make([]string, 0, len(assMap))
	for k := range assMap {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i int, j int) bool {
		return less(keys[i], keys[j])
	})

	//Create fields for each job
	fields := []*discordgo.MessageEmbedField{}
	for _, job := range keys {
		jobF := new(discordgo.MessageEmbedField)
		jobF.Name, err = database.Repo.GetJobFullName(job, guild)
		if err != nil {
			return nil, err
		}
		//Build string for each set of users
		usersStr := ""
		for _, user := range assMap[job] {
			name, err := GetUserPing(guild, user)
			if err != nil {
				return nil, err
			}
			usersStr += name
			color, err := database.Repo.GetUserColor(user, guild)
			if err != nil {
				return nil, err
			}
			if color != "" {
				usersStr += " (" + color + ")"
			}
			usersStr += ", "
		}
		usersStr = usersStr[:len(usersStr)-2]
		jobF.Value = usersStr
		fields = append(fields, jobF)
	}
	embed.Fields = fields
	return &embed, nil
}

// Builds the embed for showing user assignments
func BuildUserAssignmentsEmbed(assMap map[string][]string, user string, guild string) (*discordgo.MessageEmbed, error) {
	//Initialize embed
	embed := discordgo.MessageEmbed{
		Type: discordgo.EmbedTypeRich,
	}
	//Make title using user name
	name, err := GetUserName(guild, user)
	if err != nil {
		return nil, err
	}
	embed.Title = "Assignments for " + name

	//Create sorted job slice
	keys := make([]string, 0, len(assMap))
	for k := range assMap {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i int, j int) bool {
		return less(keys[i], keys[j])
	})

	//Create fields for each job
	fields := []*discordgo.MessageEmbedField{}
	for _, job := range keys {
		jobF := new(discordgo.MessageEmbedField)
		jobF.Name, err = database.Repo.GetJobFullName(job, guild)
		if err != nil {
			return nil, err
		}
		//Build string for each set of series
		seriesStr := ""
		for _, ser := range assMap[job] {
			nameSer, err := database.Repo.GetSeriesFullName(ser, guild)
			if err != nil {
				return nil, err
			}
			seriesStr += nameSer + ", "
		}
		seriesStr = seriesStr[:len(seriesStr)-2]
		jobF.Value = seriesStr
		fields = append(fields, jobF)
	}
	embed.Fields = fields
	return &embed, nil
}

// Builds the embed for showing job assignments
func BuildJobAssignmentsEmbed(assMap map[string][]string, job string, guild string) (*discordgo.MessageEmbed, error) {
	//Initialize embed
	embed := discordgo.MessageEmbed{
		Type: discordgo.EmbedTypeRich,
	}
	//Make title using job title
	job, err := database.Repo.GetJobFullName(job, guild)
	if err != nil {
		return nil, err
	}
	embed.Title = "Assignments for " + job

	//Create fields for each user
	fields := []*discordgo.MessageEmbedField{}
	for user, series := range assMap {
		userF := new(discordgo.MessageEmbedField)
		userF.Name, err = GetUserName(guild, user)
		if err != nil {
			return nil, err
		}
		//Build string for each set of series
		seriesStr := ""
		for _, ser := range series {
			nameSer, err := database.Repo.GetSeriesFullName(ser, guild)
			if err != nil {
				return nil, err
			}
			seriesStr += nameSer + ", "
		}
		seriesStr = seriesStr[:len(seriesStr)-2]
		userF.Value = seriesStr
		fields = append(fields, userF)
	}
	embed.Fields = fields
	return &embed, nil
}

// Builds the embed for showing all assignments. Hierarchy is series-job-user
func BuildFullAssignmentsEmbed(assMap map[string]map[string][]string, guild string) (*discordgo.MessageEmbed, error) {
	//Initialize embed
	embed := discordgo.MessageEmbed{
		Type: discordgo.EmbedTypeRich,
	}
	//Make title using server title
	guildVal, err := goBot.Guild(guild)
	if err != nil {
		return nil, err
	}
	embed.Title = "Assignments for " + guildVal.Name

	//Sort series alphabetically
	series := make([]string, 0, len(assMap))
	for k := range assMap {
		series = append(series, k)
	}
	sort.Strings(series)

	//Create fields for each series
	fields := []*discordgo.MessageEmbedField{}
	for _, ser := range series {
		serF := new(discordgo.MessageEmbedField)
		serF.Name, err = database.Repo.GetSeriesFullName(ser, guild)
		if err != nil {
			return nil, err
		}

		//Sort jobs
		jobs := make([]string, 0, len(assMap[ser]))
		for k := range assMap[ser] {
			jobs = append(jobs, k)
		}
		sort.Slice(jobs, func(i int, j int) bool {
			return less(jobs[i], jobs[j])
		})

		jobsStr := ""
		//Build job string
		for _, job := range jobs {
			jobStr, err := database.Repo.GetJobFullName(job, guild)
			if err != nil {
				return nil, err
			}
			jobStr += " - "

			//Build list of users assigned to job
			sort.Strings(assMap[ser][job])
			for _, user := range assMap[ser][job] {
				userN, err := GetUserPing(guild, user)
				if err != nil {
					return nil, err
				}
				jobStr += userN + ", "
			}
			jobStr = jobStr[:len(jobStr)-2] + "\n"
			jobsStr += jobStr
		}
		jobsStr = jobsStr[:len(jobsStr)-1]
		serF.Value = jobsStr
		fields = append(fields, serF)
	}
	embed.Fields = fields
	return &embed, nil
}

// Builds embed for showing user color prefs
func BuildColorsEmbed(assMap map[string]string, guild string) (*discordgo.MessageEmbed, error) {
	//Initialize embed
	embed := discordgo.MessageEmbed{
		Type:        discordgo.EmbedTypeRich,
		Title:       "Credit Color Preferences",
		Description: "This is the color you want your name to show up in on the credits. Change it to whatever you like with the /set_color command.",
	}
	fields := []*discordgo.MessageEmbedField{{}}

	//Get rid of all users without a set color
	none := []string{}
	for key, val := range assMap {
		if strings.TrimSpace(val) == "" {
			none = append(none, key)
		}
	}
	for _, rem := range none {
		delete(assMap, rem)
	}

	//Create a map of username to ID and a list of usernames
	namesToId := make(map[string]string)
	names := []string{}
	for key := range assMap {
		nm, err := GetUserName(guild, key)
		if err != nil {
			return nil, err
		}
		namesToId[nm] = key
		names = append(names, nm)
	}

	//Sort the list of usernames and use that to order list
	sort.Strings(names)
	text := ""
	for _, name := range names {
		ping, err := GetUserPing(guild, namesToId[name])
		if err != nil {
			return nil, err
		}
		text += ping + " - " + assMap[namesToId[name]] + "\n"
	}
	text = text[:len(text)-1]
	fields[0].Name = "Colors:"
	fields[0].Value = text
	embed.Fields = fields
	return &embed, nil
}

// Builds the embed for showing all server series
func BuildServerSeriesEmbed(guild string) (*discordgo.MessageEmbed, error) {
	//Initialize embed
	embed := discordgo.MessageEmbed{
		Type: discordgo.EmbedTypeRich,
	}
	//Make title using server title
	guildVal, err := goBot.Guild(guild)
	if err != nil {
		return nil, err
	}
	embed.Title = "Series for " + guildVal.Name
	embed.Description = "Top value is series full name, the first line under it is the series alias used for commands."

	//Get data
	series, err := database.Repo.GetAllSeries(guild)
	if err != nil {
		return nil, err
	}

	//Sort by full name alphabetical order
	sort.Slice(series, func(i int, j int) bool {
		return series[i].NameFull < series[j].NameFull
	})

	//Create fields for each series
	fields := []*discordgo.MessageEmbedField{}
	for _, ser := range series {
		serF := new(discordgo.MessageEmbedField)
		serF.Name = ser.NameFull

		serF.Value = ser.NameSh + "\n"

		channelsStr := ""
		//Build channels string
		channels, err := database.Repo.GetAllSeriesChannels(ser.NameSh, guild)
		if err != nil {
			log.Printf("Error getting channels for series %s: %s\n", ser.NameSh, err.Error())
		} else if len(channels) > 0 {
			channelsStr = "Registered channels:\n"
			for _, cha := range channels {
				men, err := GetChannelPing(cha)
				if err != nil {
					channelsStr += cha + "\n"
				} else {
					channelsStr += men + "\n"
				}
			}
		}
		serF.Value += channelsStr
		serF.Value = serF.Value[:len(serF.Value)-1]
		fields = append(fields, serF)
	}
	embed.Fields = fields
	return &embed, nil
}

// Builds the embed for showing registered users
func BuildServerUsersEmbed(guild string) (*discordgo.MessageEmbed, error) {
	//Initialize embed
	embed := discordgo.MessageEmbed{
		Type: discordgo.EmbedTypeRich,
	}
	//Make title using server title
	guildVal, err := goBot.Guild(guild)
	if err != nil {
		return nil, err
	}
	embed.Title = "Registered users for " + guildVal.Name
	embed.Description = "Unpingable users are represented with their ID. They may have left the server."

	//Create jobs field
	users := new(discordgo.MessageEmbedField)
	users.Name = "Users"

	//Get data
	allUsers, err := database.Repo.GetAllUsers(guild)
	if err != nil {
		return nil, err
	}

	//Convert IDs to pings where possible
	for i, user := range allUsers {
		ping, err := GetUserPing(guild, user)
		if err == nil {
			allUsers[i] = ping
		}
	}
	sort.Strings(allUsers)

	//Build body
	users.Value = ""
	for _, user := range allUsers {
		users.Value += fmt.Sprintf("%s\n", user)
	}
	users.Value = users.Value[:len(users.Value)-1]

	embed.Fields = []*discordgo.MessageEmbedField{users}
	return &embed, nil
}

// Builds the embed for showing server jobs
func BuildServerJobsEmbed(guild string) (*discordgo.MessageEmbed, error) {
	//Initialize embed
	embed := discordgo.MessageEmbed{
		Type: discordgo.EmbedTypeRich,
	}
	//Make title using server title
	guildVal, err := goBot.Guild(guild)
	if err != nil {
		return nil, err
	}
	embed.Title = "Available jobs for " + guildVal.Name
	embed.Description = "Values are separated by dashes; first value is shorthand, second value is full name. If there is a (G) at the end, job is a global job."

	//Create jobs field
	jobs := new(discordgo.MessageEmbedField)
	jobs.Name = "Jobs"

	//Get data
	allJobs, err := database.Repo.GetAllJobs(guild)
	if err != nil {
		return nil, err
	}
	//Sort alphabetically and remove any overwritten global jobs
	sort.Slice(allJobs, func(i int, j int) bool {
		if allJobs[i].JobSh == allJobs[j].JobSh {
			return allJobs[i].Guild == "GLOBAL"
		}
		return allJobs[i].JobSh < allJobs[j].JobSh
	})
	for i := len(allJobs) - 2; i >= 0; i-- {
		if allJobs[i].JobSh == allJobs[i+1].JobSh {
			//Why is this how Golang works
			allJobs = append(allJobs[:i], allJobs[i+1:]...)
		}
	}

	//Build body
	jobs.Value = ""
	for _, job := range allJobs {
		global := ""
		if job.Guild == "GLOBAL" {
			global = " (G)"
		}
		jobs.Value += fmt.Sprintf("%s - %s%s\n", job.JobSh, job.JobFull, global)
	}
	jobs.Value = jobs.Value[:len(jobs.Value)-1]

	embed.Fields = []*discordgo.MessageEmbedField{jobs}
	return &embed, nil
}

// Builds the embed for showing user settings
func BuildUserSettingsEmbed(user string, guild string) (*discordgo.MessageEmbed, error) {
	//Initialize embed
	embed := discordgo.MessageEmbed{
		Type: discordgo.EmbedTypeRich,
	}
	//Make title using user name
	name, err := GetUserName(guild, user)
	if err != nil {
		return nil, err
	}
	embed.Title = "Settings for " + name

	//Create settings field
	settings := new(discordgo.MessageEmbedField)
	settings.Name = "General Settings"

	//Get data
	color, err := database.Repo.GetUserColor(user, guild)
	if err != nil {
		color = "Error fetching - " + err.Error()
		log.Println("Error fetching color for settings embed: " + err.Error())
	} else if color == "" {
		color = "Unset"
	}

	v_role, err := database.Repo.GetUserVanityRole(user, guild)
	if err != nil {
		v_role = "Error fetching - " + err.Error()
		log.Println("Error fetching vanity role for settings embed: " + err.Error())
	} else if v_role == "" {
		v_role = "Unset"
	} else {
		v_role, err = GetRolePing(guild, v_role)
		if err != nil {
			v_role = "Error fetching - " + err.Error()
			log.Println("Error fetching vanity role ping for settings embed: " + err.Error())
		}
	}

	//Build body
	settings.Value = fmt.Sprintf("Credits color - %s\nVanity role - %s", color, v_role)

	embed.Fields = []*discordgo.MessageEmbedField{settings}
	return &embed, nil
}
