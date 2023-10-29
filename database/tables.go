package database

// Specs for all tables used by bot
var tableQuerys = []string{
	// Holds all data for user-created reminders
	`
    CREATE TABLE IF NOT EXISTS reminders(
		guild VARCHAR(20),
		channel VARCHAR(20),
		user VARCHAR(20),
		days INT,
		message VARCHAR(100),
		repeat BOOLEAN,
		time DATETIME
    );
    `,
	// Registers a series with a short and long name as well as its ping role
	`
    CREATE TABLE IF NOT EXISTS series(
		name_sh VARCHAR(100) COLLATE NOCASE,
		name_full VARCHAR(100) COLLATE NOCASE,
		guild VARCHAR(20),
		ping_role VARCHAR(20),
		repo_link VARCHAR(100),
		PRIMARY KEY(name_sh, guild)
    );
    `,
	// Registers a channel as belonging to a particular series
	`
    CREATE TABLE IF NOT EXISTS channels(
		channel VARCHAR(20) PRIMARY KEY,
		series VARCHAR(100) COLLATE NOCASE,
		guild VARCHAR(20)
    );
    `,
	// Registers a user as belonging to the group. Keeps track of personalized settings
	`
    CREATE TABLE IF NOT EXISTS users(
		user VARCHAR(20),
		color VARCHAR(6),
		vanity_role VARCHAR(20),
		guild VARCHAR(20),
		PRIMARY KEY(user, guild)
    );
    `,
	// Registers all user assignments to a given series
	`
    CREATE TABLE IF NOT EXISTS series_assignments(
		user VARCHAR(20),
		series VARCHAR(100) COLLATE NOCASE,
		job VARCHAR(20) COLLATE NOCASE,
		guild VARCHAR(20),
		UNIQUE(user, series, job, guild)
    );
    `,
	// Registers notes for series
	`
    CREATE TABLE IF NOT EXISTS series_notes(
		series VARCHAR(100) COLLATE NOCASE,
		note VARCHAR(1000),
		guild VARCHAR(20)
    );
    `,
	// Keeps track of all job types like TS, PR, and RD
	`CREATE TABLE IF NOT EXISTS jobs(
		job_sh VARCHAR(20) COLLATE NOCASE,
		job_full VARCHAR(100) COLLATE NOCASE,
		guild VARCHAR(20),
		PRIMARY KEY(job_sh, guild)
    );
    `,
	// Registers a group's member role
	`CREATE TABLE IF NOT EXISTS member_role(
		guild VARCHAR(20) PRIMARY KEY,
		role_id VARCHAR(20)
    );
    `,
	// Marks the start and end of work channels in the server for use in new channel generation
	`CREATE TABLE IF NOT EXISTS series_channels(
		top VARCHAR(20),
		bottom VARCHAR(20),
		guild VARCHAR(20) PRIMARY KEY
    );
    `,
	// Message for series billboard
	`CREATE TABLE IF NOT EXISTS series_billboards(
		series VARCHAR(100) COLLATE NOCASE,
		guild VARCHAR(20),
		channel VARCHAR(30),
		message VARCHAR(30),
		PRIMARY KEY(series, guild)
    );
    `,
	// Message for roles billboard
	`CREATE TABLE IF NOT EXISTS roles_billboards(
		guild VARCHAR(20) PRIMARY KEY,
		channel VARCHAR(30),
		message VARCHAR(30)
    );
    `,
	// Message for colors billboard
	`CREATE TABLE IF NOT EXISTS colors_billboards(
		guild VARCHAR(20) PRIMARY KEY,
		channel VARCHAR(30),
		message VARCHAR(30)
    );
    `,
	// Channel for bot notifications
	`CREATE TABLE IF NOT EXISTS notification_channels(
		guild VARCHAR(20) PRIMARY KEY,
		channel VARCHAR(30)
    );
    `,
}
