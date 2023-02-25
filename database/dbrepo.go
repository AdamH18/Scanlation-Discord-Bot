package database

import "time"

//https://gosamples.dev/sqlite-intro/

// Create all used tables in SQLite if not already present
func (r *SQLiteRepository) Initialize() error {
	query := `
    CREATE TABLE IF NOT EXISTS reminders(
		guild VARCHAR(20),
		channel VARCHAR(20),
		user VARCHAR(20),
		days INT,
		message VARCHAR(100),
		repeat BOOLEAN,
		time DATETIME
    );
    `

	_, err := r.db.Exec(query)
	return err
}

// Add reminder entry to DB
func (r *SQLiteRepository) AddReminder(rem Reminder) error {
	M.Lock()
	_, err := r.db.Exec("INSERT INTO reminders(guild, channel, user, days, message, repeat, time) values(?, ?, ?, ?, ?, ?, ?)", rem.Guild, rem.Channel, rem.User, rem.Days, rem.Message, rem.Repeat, rem.Time)
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
	res, err := r.db.Query("SELECT ROWID, * FROM reminders WHERE  time < ?", time.Now().Format("2006-01-02 15:04:05"))
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
