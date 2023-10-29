package bot

import (
	"fmt"
	"log"
	"scanlation-discord-bot/database"

	"github.com/bwmarrin/discordgo"
)

func UpdateSeriesBillboard(series string, guild string) {
	log.Printf("Updating series billboard for series %s and guild %s\n", series, guild)
	bill, channel, err := database.Repo.GetSeriesBillboard(series, guild)
	if err != nil {
		log.Println("Error getting billboard message: " + err.Error())
		return
	} else if bill == "" {
		log.Println("Server does not have a billboard for this series")
		return
	}

	//Billboard should be edited, so gather data
	serData, err := database.Repo.GetAllSeriesInfo(series, guild)
	if err != nil {
		log.Println("Error getting series info for billboard: " + err.Error())
		return
	}
	notes, _, err := database.Repo.GetSeriesNotes(series, guild)
	if err != nil {
		log.Println("Error getting notes info for billboard: " + err.Error())
		notes = []string{}
	}
	assMap, err := database.Repo.GetSeriesAssignments(series, guild)
	if err != nil {
		log.Println("Error getting series assignment info: " + err.Error())
		return
	}

	//Build embeds
	serInfoEmb := BuildSeriesInfoEmbed(serData, notes)
	serAssEmb, err := BuildSeriesAssignmentsEmbed(assMap, series, guild)
	if err != nil {
		log.Println("Error building assignments embed: " + err.Error())
		return
	}

	message := discordgo.MessageEdit{
		Content: &emptyStr,
		Embeds:  []*discordgo.MessageEmbed{serInfoEmb, serAssEmb},
		ID:      bill,
		Channel: channel,
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label:    "Add Note",
						Style:    discordgo.PrimaryButton,
						CustomID: fmt.Sprintf("note_add_button %s", series),
					},
					discordgo.Button{
						Label:    "Remove Note",
						Style:    discordgo.DangerButton,
						CustomID: fmt.Sprintf("note_remove_button %s", series),
					},
				},
			},
		},
	}
	_, err = goBot.ChannelMessageEditComplex(&message)
	if err != nil {
		log.Println("Error editing message: " + err.Error())
	}
}

func UpdateAllSeriesBillboards(guild string) {
	log.Printf("Updating all series billboards in guild %s\n", guild)
	allBill, err := database.Repo.GetAllSeriesBillboards(guild)
	if err != nil {
		log.Println("Error getting all billboards: " + err.Error())
		return
	}
	for _, bill := range allBill {
		go UpdateSeriesBillboard(bill, guild)
	}
}

// Edit existing assignments billboard message to reflect new data
func UpdateAssignmentsBillboard(guild string) {
	log.Println("Updating assignments billboard for guild " + guild)
	bill, channel, err := database.Repo.GetRolesBillboard(guild)
	if err != nil {
		log.Println("Error getting billboard message: " + err.Error())
		return
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
		Content: &emptyStr,
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
		return
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

	embed := BuildColorsEmbed(assMap, guild)

	message := discordgo.MessageEdit{
		Content: &emptyStr,
		Embeds:  []*discordgo.MessageEmbed{embed},
		ID:      bill,
		Channel: channel,
	}
	_, err = goBot.ChannelMessageEditComplex(&message)
	if err != nil {
		log.Println("Error editing message: " + err.Error())
	}
}
