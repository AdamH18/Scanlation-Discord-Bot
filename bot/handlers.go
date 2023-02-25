package bot

import (
	"log"
	"scanlation-discord-bot/database"

	"github.com/bwmarrin/discordgo"
)

func OptionsToMap(options []*discordgo.ApplicationCommandInteractionDataOption) map[string]*discordgo.ApplicationCommandInteractionDataOption {
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}
	return optionMap
}

func TestHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Response to test",
		},
	})
	log.Println("Test command used.")
}

func DBHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := OptionsToMap(i.ApplicationCommandData().Options)
	text := options["text"].StringValue()
	database.Repo.Create(text)
}

// Creates handlers for all slash commands based on relationship defined in commandHandlers
func CreateHandlers() {
	goBot.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}
