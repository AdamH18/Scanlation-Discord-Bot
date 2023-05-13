package bot

import (
	"fmt"
	"scanlation-discord-bot/database"
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
					"**tl**, **rd**, **ts**, **pr** - Ping the people assigned to this role and attach a message if so desired. To include a message, just press space after the /ts command and it will automatically select the optional message field",
			},
			{
				Name:  "User customization:",
				Value: "**set_color** - Set the color you wish to have your name in on credits pages",
			},
		},
	}
	return &embed
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
	//Create fields for each job
	fields := []*discordgo.MessageEmbedField{}
	for job, users := range assMap {
		jobF := new(discordgo.MessageEmbedField)
		jobF.Name, err = database.Repo.GetJobFullName(job, guild)
		if err != nil {
			return nil, err
		}
		//Build string for each set of users
		usersStr := ""
		for _, user := range users {
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
	//Create fields for each job
	fields := []*discordgo.MessageEmbedField{}
	for job, series := range assMap {
		jobF := new(discordgo.MessageEmbedField)
		jobF.Name, err = database.Repo.GetJobFullName(job, guild)
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
		jobF.Value = seriesStr
		fields = append(fields, jobF)
	}
	embed.Fields = fields
	return &embed, nil
}
