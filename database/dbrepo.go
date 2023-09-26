package database

import (
	"errors"
	"log"
	"math"
	"strings"
	"time"
)

//https://gosamples.dev/sqlite-intro/

// Create all used tables in SQLite if not already present
func (r *SQLiteRepository) Initialize() error {
	var err error
	for _, query := range tableQuerys {
		_, err = r.db.Exec(query)
		if err != nil {
			return err
		}
	}
	return nil
}

// Add reminder entry to DB
func (r *SQLiteRepository) AddReminder(ch chan (int), rem Reminder) {
	defer close(ch)
	//No reminder should repeat more often than once a day
	days := int64(math.Max(float64(rem.Days), 1.0))
	res := r.RemindersExec("INSERT INTO reminders(guild, channel, user, days, message, repeat, time) values(?, ?, ?, ?, ?, ?, ?)", rem.Guild, rem.Channel, rem.User, days, rem.Message, rem.Repeat, rem.Time)
	if res != nil {
		ch <- 0
	} else {
		ch <- -1
	}
}

func (r *SQLiteRepository) AddBounty(b Bounty) {
	r.BountiesExec("INSERT INTO bounties(customid, guild, job, series, expires, messageid, channel) values(?, ?, ?, ?, ?, ?, ?)", b.CustomID, b.Guild, b.Job, b.Series, b.Expires, b.MessageID, b.Channel)
}

func (r *SQLiteRepository) AddBountyInterestChannel(guild string, channel string) {
	r.BountyInterestExec("INSERT INTO bounty_interest(guild, channel) values(?, ?)", guild, channel)
}

// Add series entry to DB
func (r *SQLiteRepository) AddSeries(ch chan (int), ser Series) {
	defer close(ch)
	res := r.SeriesExec(ser.Guild, "INSERT INTO series(name_sh, name_full, guild, ping_role, repo_link) values(?, ?, ?, ?, ?)", strings.ToLower(ser.NameSh), ser.NameFull, ser.Guild, ser.PingRole, ser.RepoLink)
	if res != nil {
		ch <- 0
	} else {
		ch <- -1
	}
}

// Add channel entry to DB
func (r *SQLiteRepository) AddChannel(ch chan (int), cha Channel) {
	defer close(ch)
	res := r.ChannelsExec("INSERT INTO channels(channel, series, guild) values(?, ?, ?)", cha.Channel, strings.ToLower(cha.Series), cha.Guild)
	if res != nil {
		ch <- 0
	} else {
		ch <- -1
	}
}

// Add user entry to DB
func (r *SQLiteRepository) AddUser(ch chan (int), usr User) {
	defer close(ch)
	res := r.UsersExec(usr.Guild, "INSERT INTO users(user, color, vanity_role, guild) values(?, ?, ?, ?)", usr.User, usr.Color, usr.VanityRole, usr.Guild)
	if res != nil {
		ch <- 0
	} else {
		ch <- -1
	}
}

// Add job entry to DB
func (r *SQLiteRepository) AddJob(ch chan (int), job Job) {
	defer close(ch)
	res := r.JobsExec("INSERT INTO jobs(job_sh, job_full, guild) values(?, ?, ?)", strings.ToLower(job.JobSh), job.JobFull, job.Guild)
	if res != nil {
		ch <- 0
	} else {
		ch <- -1
	}
}

// Add member role entry to DB
func (r *SQLiteRepository) AddMemberRole(ch chan (int), mem MemberRole) {
	defer close(ch)
	res := r.MemberRoleExec("INSERT INTO member_role(guild, role_id) values(?, ?)", mem.Guild, mem.Role)
	if res != nil {
		ch <- 0
	} else {
		ch <- -1
	}
}

// Add series channel entry to DB
func (r *SQLiteRepository) AddSeriesChannels(ch chan (int), sec SeriesChannels) {
	defer close(ch)
	res := r.SeriesChannelsExec("REPLACE INTO series_channels(top, bottom, guild) values(?, ?, ?)", sec.Top, sec.Bottom, sec.Guild)
	if res != nil {
		ch <- 0
	} else {
		ch <- -1
	}
}

// Add roles billboard entry to DB
func (r *SQLiteRepository) AddRolesBillboard(ch chan (int), bb JobBB) {
	defer close(ch)
	res := r.RolesBillboardsExec("INSERT INTO roles_billboards(guild, channel, message) values(?, ?, ?)", bb.Guild, bb.Channel, bb.Message)
	if res != nil {
		ch <- 0
	} else {
		ch <- -1
	}
}

// Add colors billboard entry to DB
func (r *SQLiteRepository) AddColorsBillboard(ch chan (int), bb ColorBB) {
	defer close(ch)
	res := r.RolesBillboardsExec("INSERT INTO colors_billboards(guild, channel, message) values(?, ?, ?)", bb.Guild, bb.Channel, bb.Message)
	if res != nil {
		ch <- 0
	} else {
		ch <- -1
	}
}

// Add assignment entry to DB
func (r *SQLiteRepository) AddAssignment(ch chan (int), sea SeriesAssignment) {
	defer close(ch)
	res := r.SeriesAssignmentsExec(sea.Guild, "INSERT INTO series_assignments(user, series, job, guild) values(?, ?, ?, ?)", sea.User, strings.ToLower(sea.Series), strings.ToLower(sea.Job), sea.Guild)
	if res != nil {
		ch <- 0
	} else {
		ch <- -1
	}
}

// Add notification channel entry to DB
func (r *SQLiteRepository) AddNotificationChannel(ch chan (int), cha NotificationChannel) {
	defer close(ch)
	res := r.NotificationChannelsExec("REPLACE INTO notification_channels(guild, channel) values(?, ?)", cha.Guild, cha.Channel)
	if res != nil {
		ch <- 0
	} else {
		ch <- -1
	}
}

// fix this
func (r *SQLiteRepository) RemoveBounty(customID string, guild string) int64 {

	res, err := r.BountiesExec("DELETE FROM bounties WHERE customid = ? AND guild = ?", customID, guild)
	if err != nil {
		log.Printf("Error removing bounty: %s\n", err)
		return -1
	}

	rows, err := res.RowsAffected()
	if err != nil {
		log.Println("RowsAffected() error: " + err.Error())
		return -1
	}
	return rows
}

func (r *SQLiteRepository) RemoveInterestedUser(b Bounty, userid string) {
	r.BountyInterestExec("DELETE FROM bounty_interest WHERE customid = ? AND user = ?", b.CustomID, userid)
}

// Remove reminder entry by ID
func (r *SQLiteRepository) RemoveReminder(ch chan (int), id int64, guild string) {
	defer close(ch)
	res := r.RemindersExec("DELETE FROM reminders WHERE ROWID = ? AND guild = ?", id, guild)

	rows, err := res.RowsAffected()
	if err != nil {
		log.Println("RowsAffected() error: " + err.Error())
		ch <- -1
		return
	}

	ch <- int(rows)
}

// Remove reminder entry only if it belongs to specified user
func (r *SQLiteRepository) RemoveUserReminder(ch chan (int), id int64, userID string, guild string) {
	defer close(ch)
	res := r.RemindersExec("DELETE FROM reminders WHERE ROWID = ? AND user = ? AND guild = ?", id, userID, guild)

	rows, err := res.RowsAffected()
	if err != nil {
		log.Println("RowsAffected() error: " + err.Error())
		ch <- -1
		return
	}

	ch <- int(rows)
}

// Remove series entry and all references to series in other tables
func (r *SQLiteRepository) RemoveSeries(ch chan (int), nameSh string, nameFull string, guildId string) {
	defer close(ch)
	res := r.SeriesExec(guildId, "DELETE FROM series WHERE name_sh = ? AND name_full = ? AND guild = ?", nameSh, nameFull, guildId)

	rows, err := res.RowsAffected()
	if err != nil {
		log.Println("RowsAffected() error: " + err.Error())
		ch <- -1
		return
	}

	ch <- int(rows)
	// If a series was removed, remove all references to it from other tables. Can't be assed to error check, not a big deal if this fails
	if rows > 0 {
		go r.ChannelsExec("DELETE FROM channels WHERE series = ? AND guild = ?", nameSh, guildId)
		go r.SeriesAssignmentsExec(guildId, "DELETE FROM series_assignments WHERE series = ? AND guild = ?", nameSh, guildId)
		go r.SeriesBillboardsExec("DELETE FROM series_billboard WHERE series = ? AND guild = ?", nameSh, guildId)
	}
}

// Remove channel
func (r *SQLiteRepository) RemoveChannel(ch chan (int), channel string) {
	defer close(ch)
	res := r.ChannelsExec("DELETE FROM channels WHERE channel = ?", channel)

	rows, err := res.RowsAffected()
	if err != nil {
		log.Println("RowsAffected() error: " + err.Error())
		ch <- -1
		return
	}

	ch <- int(rows)
}

// Remove user and all references to user in other tables
func (r *SQLiteRepository) RemoveUser(ch chan (int), userId string, guildId string) {
	defer close(ch)
	res := r.UsersExec(guildId, "DELETE FROM users WHERE user = ? AND guild = ?", userId, guildId)

	rows, err := res.RowsAffected()
	if err != nil {
		log.Println("RowsAffected() error: " + err.Error())
		ch <- -1
		return
	}

	ch <- int(rows)
	// If a user was removed, remove all references to them from other tables. Can't be assed to error check, not a big deal if this fails
	if rows > 0 {
		go r.RemindersExec("DELETE FROM reminders WHERE user = ? AND guild = ?", userId, guildId)
		go r.SeriesAssignmentsExec(guildId, "DELETE FROM series_assignments WHERE user = ? AND guild = ?", userId, guildId)
	}
}

// Remove job and all references to job in other tables
func (r *SQLiteRepository) RemoveJob(ch chan (int), nameSh string, guildId string) {
	defer close(ch)
	res := r.JobsExec("DELETE FROM jobs WHERE job_sh = ? AND guild = ?", nameSh, guildId)

	rows, err := res.RowsAffected()
	if err != nil {
		log.Println("RowsAffected() error: " + err.Error())
		ch <- -1
		return
	}

	ch <- int(rows)
	// If a job was removed, remove all references to it from other tables. Can't be assed to error check, not a big deal if this fails
	if rows > 0 {
		go r.SeriesAssignmentsExec(guildId, "DELETE FROM series_assignments WHERE job = ? AND guild = ?", nameSh, guildId)
	}
}

// Remove member role
func (r *SQLiteRepository) RemoveMemberRole(ch chan (int), guild string) {
	defer close(ch)
	res := r.MemberRoleExec("DELETE FROM member_role WHERE guild = ?", guild)

	rows, err := res.RowsAffected()
	if err != nil {
		log.Println("RowsAffected() error: " + err.Error())
		ch <- -1
		return
	}

	ch <- int(rows)
}

// Remove series assignment
func (r *SQLiteRepository) RemoveSeriesAssignment(ch chan (int), user string, series string, job string, guild string) {
	defer close(ch)
	res := r.SeriesAssignmentsExec(guild, "DELETE FROM series_assignments WHERE user = ? AND series = ? AND job = ? AND guild = ?", user, series, job, guild)

	rows, err := res.RowsAffected()
	if err != nil {
		log.Println("RowsAffected() error: " + err.Error())
		ch <- -1
		return
	}

	ch <- int(rows)
}

// Remove all assignments for a user
func (r *SQLiteRepository) RemoveAllAssignments(ch chan (int), user string, guild string) {
	defer close(ch)
	res := r.SeriesAssignmentsExec(guild, "DELETE FROM series_assignments WHERE user = ? AND guild = ?", user, guild)

	rows, err := res.RowsAffected()
	if err != nil {
		log.Println("RowsAffected() error: " + err.Error())
		ch <- -1
		return
	}

	ch <- int(rows)
}

// Remove roles billboard
func (r *SQLiteRepository) RemoveRolesBillboard(ch chan (int), guild string) {
	defer close(ch)
	res := r.RolesBillboardsExec("DELETE FROM roles_billboards WHERE guild = ?", guild)

	rows, err := res.RowsAffected()
	if err != nil {
		log.Println("RowsAffected() error: " + err.Error())
		ch <- -1
		return
	}

	ch <- int(rows)
}

// Remove colors billboard
func (r *SQLiteRepository) RemoveColorsBillboard(ch chan (int), guild string) {
	defer close(ch)
	res := r.ColorsBillboardsExec("DELETE FROM colors_billboards WHERE guild = ?", guild)

	rows, err := res.RowsAffected()
	if err != nil {
		log.Println("RowsAffected() error: " + err.Error())
		ch <- -1
		return
	}

	ch <- int(rows)
}

// Update series name
func (r *SQLiteRepository) UpdateSeriesName(ch chan (int), nameSh string, newName string, guild string) {
	defer close(ch)
	res := r.SeriesExec(guild, "UPDATE series SET name_full = ? WHERE name_sh = ? AND guild = ?", newName, nameSh, guild)

	rows, err := res.RowsAffected()
	if err != nil {
		log.Println("RowsAffected() error: " + err.Error())
		ch <- -1
		return
	}

	ch <- int(rows)
}

func (r *SQLiteRepository) UpdateBounty(customID string, newJob string, newSeries string, newExpires int64, guild string) (bool, error) {
	res, err := r.BountiesExec("UPDATE bounties SET job = ?, series = ?, expires = ? WHERE customid = ? AND guild = ?", newJob, newSeries, newExpires, customID, guild)

	if err != nil {
		return false, err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	return rows > 0, nil
}

func (r *SQLiteRepository) DisableBounty(customID string, guild string) (bool, error) {
	res, err := r.BountiesExec("UPDATE bounties SET disabled = 1 WHERE customid = ? AND guild = ?", customID, guild)

	if err != nil {
		return false, err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	return rows > 0, nil
}

// Update series repo link
func (r *SQLiteRepository) UpdateSeriesRepoLink(ch chan (int), nameSh string, newLink string, guild string) {
	defer close(ch)
	res := r.SeriesExec(guild, "UPDATE series SET repo_link = ? WHERE name_sh = ? AND guild = ?", newLink, nameSh, guild)

	rows, err := res.RowsAffected()
	if err != nil {
		log.Println("RowsAffected() error: " + err.Error())
		ch <- -1
		return
	}

	ch <- int(rows)
}

// Update user's color to new one
func (r *SQLiteRepository) UpdateColor(ch chan (int), user string, color string, guild string) {
	defer close(ch)
	res := r.UsersExec(guild, "UPDATE users SET color = ? WHERE user = ? AND guild = ?", color, user, guild)

	rows, err := res.RowsAffected()
	if err != nil {
		log.Println("RowsAffected() error: " + err.Error())
		ch <- -1
		return
	}

	ch <- int(rows)
}

// Update user's vanity role to new one
func (r *SQLiteRepository) UpdateVanityRole(ch chan (int), user string, role string, guild string) {
	defer close(ch)
	res := r.UsersExec(guild, "UPDATE users SET vanity_role = ? WHERE user = ? AND guild = ?", role, user, guild)

	rows, err := res.RowsAffected()
	if err != nil {
		log.Println("RowsAffected() error: " + err.Error())
		ch <- -1
		return
	}

	ch <- int(rows)
}

// Add update top of series channels entry
func (r *SQLiteRepository) UpdateSeriesChannelsTop(ch chan (int), cat string, guild string) {
	defer close(ch)
	res := r.SeriesChannelsExec("UPDATE series_channels SET top = ? WHERE guild = ?", cat, guild)

	rows, err := res.RowsAffected()
	if err != nil {
		log.Println("RowsAffected() error: " + err.Error())
		ch <- -1
		return
	}

	ch <- int(rows)
}

// Add update bottom of series channels entry
func (r *SQLiteRepository) UpdateSeriesChannelsBottom(ch chan (int), cat string, guild string) {
	defer close(ch)
	res := r.SeriesChannelsExec("UPDATE series_channels SET bottom = ? WHERE guild = ?", cat, guild)

	rows, err := res.RowsAffected()
	if err != nil {
		log.Println("RowsAffected() error: " + err.Error())
		ch <- -1
		return
	}

	ch <- int(rows)
}

// Take expired reminder and add days field to alarm time to set next reminder
func (r *SQLiteRepository) ResetReminder(id int64) {
	r.RemindersExec("UPDATE reminders SET time = datetime(time, '+' || (SELECT CAST(days AS varchar(20))) || ' days') WHERE ROWID = ?", id)
}

// Return all reminders belonging to a specific user
func (r *SQLiteRepository) GetUserReminders(userID string, guild string) ([]Reminder, error) {
	res, err := r.db.Query("SELECT ROWID, * FROM reminders WHERE user = ? AND guild = ?", userID, guild)

	if err != nil {
		return nil, err
	}
	defer res.Close()

	var all []Reminder
	for res.Next() {
		var rem Reminder
		if err := res.Scan(&rem.ID, &rem.Guild, &rem.Channel, &rem.User, &rem.Days, &rem.Message, &rem.Repeat, &rem.Time); err != nil {
			return nil, err
		}
		all = append(all, rem)
	}
	return all, nil
}

// Return all reminders
func (r *SQLiteRepository) GetAllReminders(guild string) ([]Reminder, error) {
	res, err := r.db.Query("SELECT ROWID, * FROM reminders WHERE guild = ?", guild)

	if err != nil {
		return nil, err
	}
	defer res.Close()

	var all []Reminder
	for res.Next() {
		var rem Reminder
		if err := res.Scan(&rem.ID, &rem.Guild, &rem.Channel, &rem.User, &rem.Days, &rem.Message, &rem.Repeat, &rem.Time); err != nil {
			return nil, err
		}
		all = append(all, rem)
	}
	return all, nil
}

func (r *SQLiteRepository) GetAllBounties(guild string) ([]Bounty, error) {
	res, err := r.db.Query("SELECT * FROM bounties WHERE guild = ?", guild)

	if err != nil {
		return nil, err
	}
	defer res.Close()

	var all []Bounty
	for res.Next() {
		var b Bounty
		if err := res.Scan(&b.CustomID, &b.Guild, &b.Job, &b.Series, &b.Expires, &b.MessageID, &b.Channel); err != nil {
			return nil, err
		}
		all = append(all, b)
	}
	return all, nil
}

func (r *SQLiteRepository) GetBountyInterestChannel(guild string) (string, error) {
	res, err := r.db.Query("SELECT channel FROM bounty_interest WHERE guild = ?", guild)

	if err != nil {
		return "", err
	}
	defer res.Close()

	var channel string
	if res.Next() {
		if err := res.Scan(&channel); err != nil {
			return "", err
		}
	}
	return channel, nil
}

// Return all reminders for which current time is after time field
func (r *SQLiteRepository) GetActiveReminders() ([]Reminder, error) {
	res, err := r.db.Query("SELECT ROWID, * FROM reminders WHERE time < ?", time.Now().Format("2006-01-02 15:04:05"))

	if err != nil {
		return nil, err
	}
	defer res.Close()

	var all []Reminder
	for res.Next() {
		var rem Reminder
		if err := res.Scan(&rem.ID, &rem.Guild, &rem.Channel, &rem.User, &rem.Days, &rem.Message, &rem.Repeat, &rem.Time); err != nil {
			return nil, err
		}
		all = append(all, rem)
	}
	return all, nil
}

// Check if the given series shorthand exists in the database
func (r *SQLiteRepository) RegisteredSeries(ser string, guild string) bool {
	res, err := r.db.Query("SELECT * FROM series WHERE name_sh = ? AND guild = ?", ser, guild)

	if err != nil {
		log.Println("Failed to retrieve series info: " + err.Error())
		return false
	}
	defer res.Close()

	//If there was a result, return true
	return res.Next()
}

// Check if the given user exists in the database
func (r *SQLiteRepository) RegisteredUser(usr string, guild string) bool {
	res, err := r.db.Query("SELECT * FROM users WHERE user = ? AND guild = ?", usr, guild)

	if err != nil {
		log.Println("Failed to retrieve user info: " + err.Error())
		return false
	}
	defer res.Close()

	//If there was a result, return true
	return res.Next()
}

func (r *SQLiteRepository) GetBounty(customID string, guild string) (Bounty, error) {
	res, err := r.db.Query("SELECT * FROM bounties WHERE customid = ? AND guild = ?", customID, guild)

	if err != nil {
		return Bounty{}, err
	}
	defer res.Close()

	var b Bounty
	if res.Next() {
		if err := res.Scan(&b.CustomID, &b.Guild, &b.Job, &b.Series, &b.Expires, &b.MessageID, &b.Channel); err != nil {
			return Bounty{}, err
		}
	}
	return b, nil
}

// Check if the given job exists in the database
func (r *SQLiteRepository) RegisteredJob(job string, guild string) bool {
	res, err := r.db.Query("SELECT * FROM jobs WHERE (job_sh = ?) AND (guild = ? OR guild = 'GLOBAL')", job, guild)

	if err != nil {
		log.Println("Failed to retrieve job info: " + err.Error())
		return false
	}
	defer res.Close()

	//If there was a result, return true
	return res.Next()
}

// Get all registered users in server
func (r *SQLiteRepository) GetAllUsers(guild string) ([]string, error) {
	res, err := r.db.Query("SELECT user FROM users WHERE guild = ?", guild)

	if err != nil {
		return nil, err
	}
	defer res.Close()

	//Return all users
	users := []string{}
	for res.Next() {
		var user string
		if err := res.Scan(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

// Get a server's member role if registered
func (r *SQLiteRepository) GetMemberRole(guild string) string {
	res, err := r.db.Query("SELECT role_id FROM member_role WHERE guild = ?", guild)

	if err != nil {
		log.Println("Failed to retrieve role info: " + err.Error())
		return ""
	}
	defer res.Close()

	//If there was a result, return it
	var roleID string
	if res.Next() {
		err = res.Scan(&roleID)
		if err != nil {
			log.Println("Failed to read role ID: " + err.Error())
			return ""
		}
		return roleID
	} else {
		return ""
	}
}

// Get a server's series channel bounds if registered
func (r *SQLiteRepository) GetSeriesChannels(guild string) (string, string, error) {
	res, err := r.db.Query("SELECT top, bottom FROM series_channels WHERE guild = ?", guild)

	if err != nil {
		return "", "", err
	}
	defer res.Close()

	//If there was a result, return it
	var top, bottom string
	if res.Next() {
		err = res.Scan(&top, &bottom)
		if err != nil {
			return "", "", err
		}
		return top, bottom, nil
	} else {
		return "", "", errors.New("server has not registered series channel bounds")
	}
}

// Get the registered channel of a given series
func (r *SQLiteRepository) GetLocalSeries(channel string) (string, error) {
	res, err := r.db.Query("SELECT series FROM channels WHERE channel = ?", channel)

	if err != nil {
		return "", err
	}
	defer res.Close()

	//If there was a result, return it
	var series string
	if res.Next() {
		err = res.Scan(&series)
		if err != nil {
			return "", err
		}
		return series, nil
	} else {
		return "", errors.New("channel is not registered to a series")
	}
}

// Get all assignments for a given series
func (r *SQLiteRepository) GetSeriesAssignments(series string, guild string) (map[string][]string, error) {
	res, err := r.db.Query("SELECT user, job FROM series_assignments WHERE series = ? AND guild = ?", series, guild)

	if err != nil {
		return nil, err
	}
	defer res.Close()

	//Return all assignments to the series
	assignments := make(map[string][]string)
	for res.Next() {
		var user, job string
		if err := res.Scan(&user, &job); err != nil {
			return nil, err
		}
		if _, ok := assignments[job]; !ok {
			assignments[job] = []string{user}
		} else {
			assignments[job] = append(assignments[job], user)
		}
	}
	return assignments, nil
}

// Get all assignments for a given user
func (r *SQLiteRepository) GetUserAssignments(user string, guild string) (map[string][]string, error) {
	res, err := r.db.Query("SELECT series, job FROM series_assignments WHERE user = ? AND guild = ?", user, guild)

	if err != nil {
		return nil, err
	}
	defer res.Close()

	//Return all assignments to the series
	assignments := make(map[string][]string)
	for res.Next() {
		var series, job string
		if err := res.Scan(&series, &job); err != nil {
			return nil, err
		}
		if _, ok := assignments[job]; !ok {
			assignments[job] = []string{series}
		} else {
			assignments[job] = append(assignments[job], series)
		}
	}
	return assignments, nil
}

// Get all assignments for a given job
func (r *SQLiteRepository) GetJobAssignments(job string, guild string) (map[string][]string, error) {
	res, err := r.db.Query("SELECT user, series FROM series_assignments WHERE job = ? AND guild = ?", job, guild)

	if err != nil {
		return nil, err
	}
	defer res.Close()

	//Return all assignments to the job
	assignments := make(map[string][]string)
	for res.Next() {
		var user, series string
		if err := res.Scan(&user, &series); err != nil {
			return nil, err
		}
		if _, ok := assignments[user]; !ok {
			assignments[user] = []string{series}
		} else {
			assignments[user] = append(assignments[user], series)
		}
	}
	return assignments, nil
}

// Get all assignments in guild. Hierarchy is series-job-user
func (r *SQLiteRepository) GetAllAssignments(guild string) (map[string]map[string][]string, error) {
	res, err := r.db.Query("SELECT user, series, job FROM series_assignments WHERE guild = ?", guild)

	if err != nil {
		return nil, err
	}
	defer res.Close()

	//Return all assignments
	assignments := make(map[string]map[string][]string)
	for res.Next() {
		var user, series, job string
		if err := res.Scan(&user, &series, &job); err != nil {
			return nil, err
		}
		//If new series, add to map
		if _, ok := assignments[series]; !ok {
			assignments[series] = make(map[string][]string)
		}
		//If new job within series, add to jobs
		if _, ok := assignments[series][job]; !ok {
			assignments[series][job] = []string{user}
		} else {
			assignments[series][job] = append(assignments[series][job], user)
		}
	}
	return assignments, nil
}

// Get all assignments for a given series and job
func (r *SQLiteRepository) GetSeriesJobAssignments(series string, job string, guild string) ([]string, error) {
	res, err := r.db.Query("SELECT user FROM series_assignments WHERE series = ? AND job = ? AND guild = ?", series, job, guild)

	if err != nil {
		return nil, err
	}
	defer res.Close()

	//Return all assignments to the series
	users := []string{}
	for res.Next() {
		var user string
		if err := res.Scan(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

// Get the full name of a series from its shorthand
func (r *SQLiteRepository) GetSeriesFullName(seriesSh string, guild string) (string, error) {
	res, err := r.db.Query("SELECT name_full FROM series WHERE name_sh = ? AND guild = ?", seriesSh, guild)

	if err != nil {
		return "", err
	}
	defer res.Close()

	//If there was a result, return it
	var seriesFull string
	if res.Next() {
		err = res.Scan(&seriesFull)
		if err != nil {
			return "", err
		}
		return seriesFull, nil
	} else {
		return "", errors.New("failed to get full name from DB")
	}
}

// Get full name and short name of all series in server
func (r *SQLiteRepository) GetAllSeries(guild string) ([]Series, error) {
	res, err := r.db.Query("SELECT name_sh, name_full FROM series WHERE guild = ?", guild)

	if err != nil {
		return nil, err
	}
	defer res.Close()

	//Return array of all results
	var nameSh, nameFull string
	series := []Series{}
	for res.Next() {
		err = res.Scan(&nameSh, &nameFull)
		if err != nil {
			return nil, err
		}
		ser := Series{
			NameSh:   nameSh,
			NameFull: nameFull,
		}
		series = append(series, ser)
	}
	return series, nil
}

// Get all channels registered to a series
func (r *SQLiteRepository) GetAllSeriesChannels(series string, guild string) ([]string, error) {
	res, err := r.db.Query("SELECT channel FROM channels WHERE series = ? AND guild = ?", series, guild)

	if err != nil {
		return nil, err
	}
	defer res.Close()

	//Return array of all results
	var channel string
	channels := []string{}
	for res.Next() {
		err = res.Scan(&channel)
		if err != nil {
			return nil, err
		}
		channels = append(channels, channel)
	}
	return channels, nil
}

// Get get full name, short name, and locality of all jobs in server
func (r *SQLiteRepository) GetAllJobs(guild string) ([]Job, error) {
	res, err := r.db.Query("SELECT * FROM jobs WHERE guild = ? OR guild = 'GLOBAL'", guild)

	if err != nil {
		return nil, err
	}
	defer res.Close()

	//Return array of all results
	var jobSh, jobFull, gld string
	jobs := []Job{}
	for res.Next() {
		err = res.Scan(&jobSh, &jobFull, &gld)
		if err != nil {
			return nil, err
		}
		job := Job{
			JobSh:   jobSh,
			JobFull: jobFull,
			Guild:   gld,
		}
		jobs = append(jobs, job)
	}
	return jobs, nil
}

// Get the full name of a job from its shorthand
func (r *SQLiteRepository) GetJobFullName(jobSh string, guild string) (string, error) {
	res, err := r.db.Query("SELECT job_full, guild FROM jobs WHERE (job_sh = ?) AND (guild = ? OR guild = 'GLOBAL')", jobSh, guild)

	if err != nil {
		return "", err
	}
	defer res.Close()

	//If there was a result, return it
	var jobFull, gld string
	if res.Next() {
		err = res.Scan(&jobFull, &gld)
		if err != nil {
			return "", err
		}
		//Prioritize local role if global overridden
		if res.Next() && guild == "GLOBAL" {
			err = res.Scan(&jobFull, &gld)
			if err != nil {
				return "", err
			}
			return jobFull, nil
		}
		return jobFull, nil
	} else {
		return "", errors.New("failed to get full name from DB")
	}
}

// Get preferred color of a user
func (r *SQLiteRepository) GetUserColor(user string, guild string) (string, error) {
	res, err := r.db.Query("SELECT color FROM users WHERE user = ? AND guild = ?", user, guild)

	if err != nil {
		return "", err
	}
	defer res.Close()

	//If there was a result, return it
	var color string
	if res.Next() {
		err = res.Scan(&color)
		if err != nil {
			return "", err
		}
		return color, nil
	} else {
		return "", errors.New("failed to get color from DB")
	}
}

// Get preferred color of all users in server
func (r *SQLiteRepository) GetAllColors(guild string) (map[string]string, error) {
	res, err := r.db.Query("SELECT user, color FROM users WHERE guild = ?", guild)

	if err != nil {
		return nil, err
	}
	defer res.Close()

	//Return all color preferences
	assignments := make(map[string]string)
	for res.Next() {
		var user, color string
		if err := res.Scan(&user, &color); err != nil {
			return nil, err
		}
		assignments[user] = color
	}
	return assignments, nil
}

// Get vanity role of a user
func (r *SQLiteRepository) GetUserVanityRole(user string, guild string) (string, error) {
	res, err := r.db.Query("SELECT vanity_role FROM users WHERE user = ? AND guild = ?", user, guild)

	if err != nil {
		return "", err
	}
	defer res.Close()

	//If there was a result, return it
	var role string
	if res.Next() {
		err = res.Scan(&role)
		if err != nil {
			return "", err
		}
		return role, nil
	} else {
		return "", errors.New("failed to get role from DB")
	}
}

// Get number of users with a given vanity role
func (r *SQLiteRepository) NumUsersWithVanity(role string, guild string) (int, error) {
	res, err := r.db.Query("SELECT user FROM users WHERE vanity_role = ? AND guild = ?", role, guild)

	if err != nil {
		return 0, err
	}
	defer res.Close()

	//If there was a result, count up and return it
	count := 0
	for res.Next() {
		count++
	}
	return count, nil
}

// Get roles billboard message ID in guild
func (r *SQLiteRepository) GetRolesBillboard(guild string) (string, string, error) {
	res, err := r.db.Query("SELECT message, channel FROM roles_billboards WHERE guild = ?", guild)

	if err != nil {
		return "", "", err
	}
	defer res.Close()

	//If there was a result, return it
	var message, channel string
	if res.Next() {
		err = res.Scan(&message, &channel)
		if err != nil {
			return "", "", err
		}
		return message, channel, nil
	} else {
		return "", "", nil
	}
}

// Get colors billboard message ID in guild
func (r *SQLiteRepository) GetColorsBillboard(guild string) (string, string, error) {
	res, err := r.db.Query("SELECT message, channel FROM colors_billboards WHERE guild = ?", guild)

	if err != nil {
		return "", "", err
	}
	defer res.Close()

	//If there was a result, return it
	var message, channel string
	if res.Next() {
		err = res.Scan(&message, &channel)
		if err != nil {
			return "", "", err
		}
		return message, channel, nil
	} else {
		return "", "", nil
	}
}

// Return all notification channels
func (r *SQLiteRepository) GetAllNotificationChannels() ([]NotificationChannel, error) {
	res, err := r.db.Query("SELECT * FROM notification_channels")

	if err != nil {
		return nil, err
	}
	defer res.Close()

	var all []NotificationChannel
	for res.Next() {
		var notif NotificationChannel
		if err := res.Scan(&notif.Guild, &notif.Channel); err != nil {
			return nil, err
		}
		all = append(all, notif)
	}
	return all, nil
}
