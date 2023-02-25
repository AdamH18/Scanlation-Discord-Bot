package bot

import (
	"fmt"
	"scanlation-discord-bot/database"
	"strings"
	"text/tabwriter"
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
