package bot

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

// Helper command to turn Discord Options array to a map of labels to options
func OptionsToMap(options []*discordgo.ApplicationCommandInteractionDataOption) map[string]*discordgo.ApplicationCommandInteractionDataOption {
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}
	return optionMap
}

// Takes a user ID and returns the username
func GetUserName(guildID string, userID string) (string, error) {
	usr, err := goBot.GuildMember(guildID, userID)
	if err != nil {
		return "", err
	}
	return usr.User.Username, nil
}

// Takes a user ID and returns a ping string
func GetUserPing(guildID string, userID string) (string, error) {
	usr, err := goBot.GuildMember(guildID, userID)
	if err != nil {
		return "", err
	}
	return usr.Mention(), nil
}

// Handles response to a slash command
func Respond(s *discordgo.Session, i *discordgo.InteractionCreate, response string) {
	log.Printf("Response: %s\n", response)
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: response,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}

func LogCommand(i *discordgo.InteractionCreate, name string) {
	log.Printf("User %s (%s) in guild %s and channel %s used %s command with options:\n", i.Member.User.Username, i.Member.User.ID, i.GuildID, i.ChannelID, name)
}

func IsCommand(n string) bool {
	for _, cmd := range commands {
		if cmd.Name == n {
			return true
		}
	}
	return false
}
