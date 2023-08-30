# Scanlation-Discord-Bot
This is a utility bot for scanlation tracking and workflow organization.

## Deployment (using docker)

To deploy the bot using Docker, follow these steps:

1. Clone this repository:

```sh
git clone https://github.com/AdamH18/Scanlation-Discord-Bot.git
cd ./Scanlation-Discord-Bot
```

2. Make sure the configuration is available:

```sh
mkdir -p ./app
touch ./app/sqlite.db
cp ./config-template.json ./app/config.json
```

The configuration scheme is as follows:

| config | default | value | description |
|--------|---------|-------|-------------|
| Token | none | string | Required, [bot's access token](https://discord.com/developers/docs/reference#authentication) |
| RemoveCommands | false | bool | Choose whether or not to remove registered slash commands from Discord upon shutdown |
| DatabaseFile | sqlite.db | string | Location of database file for bot to use |
| DatabaseBackupChannel | none | string | Discord channel in which to dump DB backups |
| Owner | none | string | User ID of bot owner |

3. Launch the bot:

```sh
docker compose up -d
```

## Available commands

| Command | Permission | Description |
|---------|------------|-------------|
| /help | all | Show help |
| /add_reminder | all | Add reminder for yourself |
| /my_reminders | all | Show all personal reminders |
| /rem_reminder | all | Remove reminder for yourself |
| /set_alarm | all | Set an alarm for yourself |
| /user_reminders | admin only |Show all reminders for a user |
| /all_reminders | admin only | Show all reminders |
| /set_any_alarm | admin only | Set alarm for any user |
| /add_any_reminder | admin only | Add reminder for a user |
| /rem_any_reminder | admin only | Remove reminder for any user |
| /add_series | admin only | Register new series for group |
| /remove_series | admin only | Removes series for group, including all related settings. Channels are not deleted. |
| /server_series | admin only | See all existing series on the server |
| /change_series_title | admin only | Changes the full name of the series. Shorthand name is unchanged |
| /change_series_repo | admin only | Changes the repo link of the series |
| /add_series_channel | admin only | Register a channel with a given series |
| /remove_series_channel | admin only | Deregister a channel with a given series, channel is not deleted |
| /add_user | admin only | Register a user as a member of the group |
| /remove_user | admin only | Remove a user from the group, deletes all related settings. User is not kicked |
| /server_users | admin only | See all registered users on the server |
| /add_job | admin only | Register a new job type for the group |
| /add_global_job | owner only | Register a new job type for all users |
| /remove_job | admin only | Remove a job type for the group, including all assignments to that job |
| /server_jobs | admin only | See all existing jobs on the server |
| /add_member_role | admin only | Registers the role used to determine group members |
| /remove_member_role | admin only | Deregisters the role used to determine group members. Role is not deleted |
| /reg_series_channels | admin only | Registers bounds for series channels. Should be IDs of first and last categories |
| /add_series_assignment | admin only | Register an assignment to a series for a group member |
| /remove_series_assignment | admin only | Remove an assignment to a series for a group member |
| /remove_all_assignments | admin only | Remove all assignments for a group member. Does not kick member from group |
| /series_assignments | admin only | See the assignments and user colors for a given series |
| /my_assignments | all | See your personal assignments |
| /user_assignments | all | See the assignments of a given user |
| /job_assignments | all | See everyone assigned to a given job |
| /tl | all | Ping the translator(s) assigned to a series |
| /rd | all | Ping the redrawer(s) assigned to a series |
| /ts | all | Ping the typesetter(s) assigned to a series |
| /pr | all | Ping the proofreader(s) assigned to a series |
| /my_settings | all | See your server settings |
| /user_settings | admin only | See user's server settings |
| /set_color | all | Set your color for credits pages |
| /set_user_color | admin only | Set user color for credits pages |
| /vanity_role | all | Give yourself a vanity role |
| /rem_vanity_role | all | Removes your vanity role |
| /create_series_billboard | admin only | Create a billboard showcasing series information in this channel |
| /delete_series_billboard | admin only | Deregister the billboard showcasing series information. Does not delete message |
| /create_assignments_billboard | admin only | Create a billboard showcasing all assignments in this channel |
| /delete_assignments_billboard | admin only | Deregister the billboard showcasing all assignments. Does not delete message |
| /create_colors_billboard | admin only | Create a billboard showcasing all color prefs in this channel |
| /delete_colors_billboard | admin only | Deregister the billboard showcasing all color prefs. Does not delete message | 
| /refresh_all_billboards | admin only | Refreshes all billboards on the server |
| /add_notification_channel | admin only | Sets channel to receive messages containing updates from bot owner |
| /send_notification | owner only | Send message to all registered notification channels |
