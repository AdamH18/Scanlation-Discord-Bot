package bot

import (
	"log"
	"scanlation-discord-bot/database"

	"github.com/bwmarrin/discordgo"
)

func UpdateSeriesBillboard(series string, guild string) {

}

func UpdateAllSeriesBillboards(guild string) {

}

// Edit existing assignments billboard message to reflect new data
func UpdateAssignmentsBillboard(guild string) {
	log.Println("Updating assignments billboard for guild " + guild)
	bill, channel, err := database.Repo.GetRolesBillboard(guild)
	if err != nil {
		log.Println("Error getting billboard message: " + err.Error())
	} else if bill == "" {
		log.Println("Server does not have an assignments billboard")
		return
	}

	//Billboard should be edited, so gather data
	assMap, err := database.Repo.GetAllAssignments(guild)
	if err != nil {
		log.Println("Error getting assignments data: " + err.Error())
		return
	}

	embed, err := BuildFullAssignmentsEmbed(assMap, guild)
	if err != nil {
		log.Println("Error building embed: " + err.Error())
		return
	}

	message := discordgo.MessageEdit{
		Embeds:  []*discordgo.MessageEmbed{embed},
		ID:      bill,
		Channel: channel,
	}
	_, err = goBot.ChannelMessageEditComplex(&message)
	if err != nil {
		log.Println("Error editing message: " + err.Error())
	}
}

// Edit existing colors billboard message to reflect new data
func UpdateColorsBillboard(guild string) {
	log.Println("Updating colors billboard for guild " + guild)
	bill, channel, err := database.Repo.GetColorsBillboard(guild)
	if err != nil {
		log.Println("Error getting billboard message: " + err.Error())
	} else if bill == "" {
		log.Println("Server does not have a colors billboard")
		return
	}

	//Billboard should be edited, so gather data
	assMap, err := database.Repo.GetAllColors(guild)
	if err != nil {
		log.Println("Error getting user color data: " + err.Error())
		return
	}

	embed, err := BuildColorsEmbed(assMap, guild)
	if err != nil {
		log.Println("Error building embed: " + err.Error())
		return
	}

	message := discordgo.MessageEdit{
		Embeds:  []*discordgo.MessageEmbed{embed},
		ID:      bill,
		Channel: channel,
	}
	_, err = goBot.ChannelMessageEditComplex(&message)
	if err != nil {
		log.Println("Error editing message: " + err.Error())
	}
}
