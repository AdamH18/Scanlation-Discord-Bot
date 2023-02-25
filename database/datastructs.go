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
