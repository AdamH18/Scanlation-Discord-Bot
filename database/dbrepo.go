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
func (r *SQLiteRepository) AddReminder(rem Reminder) error {
	//No reminder should repeat more often than once a day
	days := int64(math.Max(float64(rem.Days), 1.0))
	M.Lock()
	_, err := r.db.Exec("INSERT INTO reminders(guild, channel, user, days, message, repeat, time) values(?, ?, ?, ?, ?, ?, ?)", rem.Guild, rem.Channel, rem.User, days, rem.Message, rem.Repeat, rem.Time)
	M.Unlock()
	if err != nil {
		return err
	}

	return nil
}

// Add series entry to DB
func (r *SQLiteRepository) AddSeries(ser Series) error {
	M.Lock()
	_, err := r.db.Exec("INSERT INTO series(name_sh, name_full, guild, ping_role) values(?, ?, ?, ?)", strings.ToLower(ser.NameSh), ser.NameFull, ser.Guild, ser.PingRole)
	M.Unlock()
	if err != nil {
		return err
	}

	return nil
}

// Add channel entry to DB
func (r *SQLiteRepository) AddChannel(cha Channel) error {
	M.Lock()
	_, err := r.db.Exec("INSERT INTO channels(channel, series, guild) values(?, ?, ?)", cha.Channel, strings.ToLower(cha.Series), cha.Guild)
	M.Unlock()
	if err != nil {
		return err
	}

	return nil
}

// Add user entry to DB
func (r *SQLiteRepository) AddUser(usr User) error {
	M.Lock()
	_, err := r.db.Exec("INSERT INTO users(user, color, guild) values(?, ?, ?)", usr.User, usr.Color, usr.Guild)
	M.Unlock()
	if err != nil {
		return err
	}

	return nil
}

// Add job entry to DB
func (r *SQLiteRepository) AddJob(job Job) error {
	M.Lock()
	_, err := r.db.Exec("INSERT INTO jobs(job_sh, job_full, guild) values(?, ?, ?)", strings.ToLower(job.JobSh), job.JobFull, job.Guild)
	M.Unlock()
	if err != nil {
		return err
	}

	return nil
}

// Add member role entry to DB
func (r *SQLiteRepository) AddMemberRole(mem MemberRole) error {
	M.Lock()
	_, err := r.db.Exec("INSERT INTO member_role(guild, role_id) values(?, ?)", mem.Guild, mem.Role)
	M.Unlock()
	if err != nil {
		return err
	}

	return nil
}

// Add series channel entry to DB
func (r *SQLiteRepository) AddSeriesChannels(sec SeriesChannels) error {
	M.Lock()
	_, err := r.db.Exec("REPLACE INTO series_channels(top, bottom, guild) values(?, ?, ?)", sec.Top, sec.Bottom, sec.Guild)
	M.Unlock()
	if err != nil {
		return err
	}

	return nil
}

// Add assignment entry to DB
func (r *SQLiteRepository) AddAssignment(sea SeriesAssignment) error {
	M.Lock()
	_, err := r.db.Exec("INSERT INTO series_assignments(user, series, job, guild) values(?, ?, ?, ?)", sea.User, strings.ToLower(sea.Series), strings.ToLower(sea.Job), sea.Guild)
	M.Unlock()
	if err != nil {
		return err
	}

	return nil
}

// Remove reminder entry by ID
func (r *SQLiteRepository) RemoveReminder(id int64) (int64, error) {
	M.Lock()
	res, err := r.db.Exec("DELETE FROM reminders WHERE ROWID = ?", id)
	M.Unlock()
	if err != nil {
		return 0, err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rows, nil
}

// Remove series entry and all references to series in other tables
func (r *SQLiteRepository) RemoveSeries(nameSh string, nameFull string, guildId string) (bool, error) {
	M.Lock()
	res, err := r.db.Exec("DELETE FROM series WHERE name_sh = ? AND name_full = ? AND guild = ?", nameSh, nameFull, guildId)
	M.Unlock()
	if err != nil {
		return false, err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	done := rows > 0
	// If a series was removed, remove all references to it from other tables. Can't be assed to error check, not a big deal if this fails
	if done {
		r.db.Exec("DELETE FROM channels WHERE series = ? AND guild = ?", nameSh, guildId)
		r.db.Exec("DELETE FROM series_assignments WHERE series = ? AND guild = ?", nameSh, guildId)
		r.db.Exec("DELETE FROM series_billboard WHERE series = ? AND guild = ?", nameSh, guildId)
	}

	return done, nil
}

// Remove channel
func (r *SQLiteRepository) RemoveChannel(channel string) (bool, error) {
	M.Lock()
	res, err := r.db.Exec("DELETE FROM channels WHERE channel = ?", channel)
	M.Unlock()
	if err != nil {
		return false, err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	return rows > 0, nil
}

// Remove user and all references to user in other tables
func (r *SQLiteRepository) RemoveUser(userId string, guildId string) (bool, error) {
	M.Lock()
	res, err := r.db.Exec("DELETE FROM users WHERE user = ? AND guild = ?", userId, guildId)
	M.Unlock()
	if err != nil {
		return false, err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	done := rows > 0
	// If a user was removed, remove all references to them from other tables. Can't be assed to error check, not a big deal if this fails
	if done {
		r.db.Exec("DELETE FROM reminders WHERE user = ? AND guild = ?", userId, guildId)
		r.db.Exec("DELETE FROM series_assignments WHERE user = ? AND guild = ?", userId, guildId)
	}

	return done, nil
}

// Remove job and all references to job in other tables
func (r *SQLiteRepository) RemoveJob(nameSh string, guildId string) (bool, error) {
	M.Lock()
	res, err := r.db.Exec("DELETE FROM jobs WHERE job_sh = ? AND guild = ?", nameSh, guildId)
	M.Unlock()
	if err != nil {
		return false, err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	done := rows > 0
	// If a job was removed, remove all references to it from other tables. Can't be assed to error check, not a big deal if this fails
	if done {
		r.db.Exec("DELETE FROM series_assignments WHERE job = ? AND guild = ?", nameSh, guildId)
	}

	return done, nil
}

// Remove member role
func (r *SQLiteRepository) RemoveMemberRole(guild string) (bool, error) {
	M.Lock()
	res, err := r.db.Exec("DELETE FROM member_role WHERE guild = ?", guild)
	M.Unlock()
	if err != nil {
		return false, err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	return rows > 0, nil
}

// Remove series assignment
func (r *SQLiteRepository) RemoveSeriesAssignment(user string, series string, job string, guild string) (bool, error) {
	M.Lock()
	res, err := r.db.Exec("DELETE FROM series_assignments WHERE user = ? AND series = ? AND job = ? AND guild = ?", user, series, job, guild)
	M.Unlock()
	if err != nil {
		return false, err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	return rows > 0, nil
}

// Remove all assignments for a user
func (r *SQLiteRepository) RemoveAllAssignments(user string, guild string) (bool, error) {
	M.Lock()
	res, err := r.db.Exec("DELETE FROM series_assignments WHERE user = ? AND guild = ?", user, guild)
	M.Unlock()
	if err != nil {
		return false, err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	return rows > 0, nil
}

// Update user's color to new one
func (r *SQLiteRepository) UpdateColor(user string, color string, guild string) error {
	M.Lock()
	_, err := r.db.Exec("UPDATE users SET color = ? WHERE user = ? AND guild = ?", color, user, guild)
	M.Unlock()
	if err != nil {
		return err
	}

	return nil
}

// Add update top of series channels entry
func (r *SQLiteRepository) UpdateSeriesChannelsTop(cat string, guild string) {
	M.Lock()
	_, err := r.db.Exec("UPDATE series_channels SET top = ? WHERE guild = ?", cat, guild)
	M.Unlock()
	if err != nil {
		log.Print("Error updating top of series channels: " + err.Error())
	}
}

// Add update bottom of series channels entry
func (r *SQLiteRepository) UpdateSeriesChannelsBottom(cat string, guild string) {
	M.Lock()
	_, err := r.db.Exec("UPDATE series_channels SET bottom = ? WHERE guild = ?", cat, guild)
	M.Unlock()
	if err != nil {
		log.Print("Error updating bottom of series channels: " + err.Error())
	}
}

// Take expired reminder and add days field to alarm time to set next reminder
func (r *SQLiteRepository) ResetReminder(id int64) error {
	M.Lock()
	_, err := r.db.Exec("UPDATE reminders SET time = datetime(time, '+' || (SELECT CAST(days AS varchar(20))) || ' days') WHERE ROWID = ?", id)
	M.Unlock()
	if err != nil {
		return err
	}

	return nil
}

// Remove reminder entry only if it belongs to specified user
func (r *SQLiteRepository) RemoveUserReminder(id int64, userID string) (int64, error) {
	M.Lock()
	res, err := r.db.Exec("DELETE FROM reminders WHERE ROWID = ? AND user = ?", id, userID)
	M.Unlock()
	if err != nil {
		return 0, err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rows, nil
}

// Return all reminders belonging to a specific user
func (r *SQLiteRepository) GetUserReminders(userID string) ([]Reminder, error) {
	M.Lock()
	res, err := r.db.Query("SELECT ROWID, * FROM reminders WHERE user = ?", userID)
	M.Unlock()
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
func (r *SQLiteRepository) GetAllReminders() ([]Reminder, error) {
	M.Lock()
	res, err := r.db.Query("SELECT ROWID, * FROM reminders")
	M.Unlock()
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

// Return all reminders for which current time is after time field
func (r *SQLiteRepository) GetActiveReminders() ([]Reminder, error) {
	M.Lock()
	res, err := r.db.Query("SELECT ROWID, * FROM reminders WHERE time < ?", time.Now().Format("2006-01-02 15:04:05"))
	M.Unlock()
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
	M.Lock()
	res, err := r.db.Query("SELECT * FROM series WHERE name_sh = ? AND guild = ?", ser, guild)
	M.Unlock()
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
	M.Lock()
	res, err := r.db.Query("SELECT * FROM users WHERE user = ? AND guild = ?", usr, guild)
	M.Unlock()
	if err != nil {
		log.Println("Failed to retrieve user info: " + err.Error())
		return false
	}
	defer res.Close()

	//If there was a result, return true
	return res.Next()
}

// Check if the given job exists in the database
func (r *SQLiteRepository) RegisteredJob(job string, guild string) bool {
	M.Lock()
	res, err := r.db.Query("SELECT * FROM jobs WHERE (job_sh = ?) AND (guild = ? OR guild = 'GLOBAL')", job, guild)
	M.Unlock()
	if err != nil {
		log.Println("Failed to retrieve job info: " + err.Error())
		return false
	}
	defer res.Close()

	//If there was a result, return true
	return res.Next()
}

// Get a server's member role if registered
func (r *SQLiteRepository) GetMemberRole(guild string) string {
	M.Lock()
	res, err := r.db.Query("SELECT role_id FROM member_role WHERE guild = ?", guild)
	M.Unlock()
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
	M.Lock()
	res, err := r.db.Query("SELECT top, bottom FROM series_channels WHERE guild = ?", guild)
	M.Unlock()
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
	M.Lock()
	res, err := r.db.Query("SELECT series FROM channels WHERE channel = ?", channel)
	M.Unlock()
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
	M.Lock()
	res, err := r.db.Query("SELECT user, job FROM series_assignments WHERE series = ? AND guild = ?", series, guild)
	M.Unlock()
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
	M.Lock()
	res, err := r.db.Query("SELECT series, job FROM series_assignments WHERE user = ? AND guild = ?", user, guild)
	M.Unlock()
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

// Get all assignments for a given series and job
func (r *SQLiteRepository) GetSeriesJobAssignments(series string, job string, guild string) ([]string, error) {
	M.Lock()
	res, err := r.db.Query("SELECT user FROM series_assignments WHERE series = ? AND job = ? AND guild = ?", series, job, guild)
	M.Unlock()
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
	M.Lock()
	res, err := r.db.Query("SELECT name_full FROM series WHERE name_sh = ? AND guild = ?", seriesSh, guild)
	M.Unlock()
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

// Get the full name of a job from its shorthand
func (r *SQLiteRepository) GetJobFullName(jobSh string, guild string) (string, error) {
	M.Lock()
	res, err := r.db.Query("SELECT job_full, guild FROM jobs WHERE (job_sh = ?) AND (guild = ? OR guild = 'GLOBAL')", jobSh, guild)
	M.Unlock()
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
	M.Lock()
	res, err := r.db.Query("SELECT color FROM users WHERE user = ? AND guild = ?", user, guild)
	M.Unlock()
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
