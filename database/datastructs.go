package database

import "fmt"

//Contains all data corresponding to a reminder
type Reminder struct {
	ID      int
	Guild   string
	Channel string
	User    string
	Days    int64
	Message string
	Repeat  bool
	Time    string
}

func (r Reminder) String() string {
	return fmt.Sprintf("%d;%s;%s;%s;%d;%s;%t;%s", r.ID, r.Guild, r.Channel, r.User, r.Days, r.Message, r.Repeat, r.Time)
}

//Contains all data corresponding to a bounty
type Bounty struct {
	CustomID    string
	Guild       string
	Job         string
	Series      string
	Expires     int64
	MessageID   string
	Channel     string
	Description string
}

type BountyInterest struct {
	ChannelID string
	Guild     string
}

type Series struct {
	NameSh   string
	NameFull string
	Guild    string
	PingRole string
	RepoLink string
}

type Channel struct {
	Channel string
	Series  string
	Guild   string
}

type User struct {
	User       string
	Color      string
	VanityRole string
	Guild      string
}

type Job struct {
	JobSh   string
	JobFull string
	Guild   string
}

type MemberRole struct {
	Guild string
	Role  string
}

type SeriesChannels struct {
	Top    string
	Bottom string
	Guild  string
}

type SeriesAssignment struct {
	User   string
	Series string
	Job    string
	Guild  string
}

type JobBB struct {
	Guild   string
	Channel string
	Message string
}

type ColorBB struct {
	Guild   string
	Channel string
	Message string
}

type NotificationChannel struct {
	Guild   string
	Channel string
}
