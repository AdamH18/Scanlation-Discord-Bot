package bot

import (
	"fmt"
	"log"
	"os"
	"scanlation-discord-bot/config"
	"scanlation-discord-bot/database"
	"time"
)

// Handles the sending of reminders identified from CheckReminders
func SendReminder(rem database.Reminder) error {
	//Get string for user ping
	ping, err := GetUserPing(rem.Guild, rem.User)
	if err != nil {
		return err
	}

	message := fmt.Sprintf("Reminder for %s: %s\n\n", ping, rem.Message)
	if rem.Repeat {
		//If supposed to repeat, add defined number of days to DB time for next reminder
		message = message + fmt.Sprintf("Message is set to repeat. If no longer needed, delete using ID %d", rem.ID)
		err = database.Repo.ResetReminder(int64(rem.ID))
	} else {
		//If not supposed to repeat, just delete
		message = message + "Message is not set to repeat. Removing from database"
		_, err = database.Repo.RemoveReminder(int64(rem.ID))
	}
	if err != nil {
		return err
	}
	//Actually send the message
	_, err = goBot.ChannelMessageSend(rem.Channel, message)
	if err != nil {
		return err
	}
	return nil
}

// Runs every hour to check what reminders are ready to send
func CheckReminders() {
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-quit:
			return
		case <-ticker.C:
			log.Println("Checking for active reminders")
			rems, err := database.Repo.GetActiveReminders()
			if err != nil {
				log.Printf("Error checking active reminders: %s\n", err)
				continue
			}
			if len(rems) == 0 {
				log.Println("No reminders to send")
				continue
			}
			//Send every reminder identified
			for _, rem := range rems {
				log.Printf("Sending reminder %s\n", rem)
				err := SendReminder(rem)
				if err != nil {
					log.Printf("Error sending reminder: %s\n", err)
				}
			}
		}
	}
}

// Send DB backups to channel identified in config
func BackupDB() {
	ticker := time.NewTicker(time.Hour * 24)
	defer ticker.Stop()

	for {
		select {
		case <-quit:
			return
		case <-ticker.C:
			log.Println("Backing up DB")
			name := "DB_" + time.Now().Format(time.RFC3339) + ".db"
			database.M.Lock()
			r, err := os.Open(config.DatabaseFile)
			if err != nil {
				log.Printf("Error opening DB file: %s\n", err)
				database.M.Unlock()
				continue
			}
			_, err = goBot.ChannelFileSend(config.DatabaseBackupChannel, name, r)
			r.Close()
			database.M.Unlock()
			if err != nil {
				log.Printf("Error backing up file: %s\n", err)
			}
		}
	}
}
