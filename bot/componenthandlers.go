package bot

import (
	"log"
	"scanlation-discord-bot/database"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

// Handler for note_add_button
func NoteAddButtonHandler(s *discordgo.Session, i *discordgo.InteractionCreate, identifier string) {
	LogCommand(i, "note_add_button")

	if !database.Repo.RegisteredUser(i.Member.User.ID, i.GuildID) {
		Respond(s, i, "You are not registered with this group, please get registered before adding notes.")
		return
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: "note_add_modal " + identifier,
			Title:    "Add Note",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "note",
							Label:       "Add Note",
							Style:       discordgo.TextInputParagraph,
							Placeholder: "Write your note here",
							Required:    true,
						},
					},
				},
			},
		},
	})
	if err != nil {
		log.Printf("Add note modal interaction failed for guild %s and identifier %s: %s\n", i.GuildID, identifier, err.Error())
	}
}

// Handler for note_add_modal
func NoteAddModalHandler(s *discordgo.Session, i *discordgo.InteractionCreate, identifier string) {
	LogCommand(i, "note_add_modal")

	data := i.ModalSubmitData()
	note := data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value

	serNote := database.SeriesNote{
		Series: identifier,
		Note:   note,
		Guild:  i.GuildID,
	}
	err := database.Repo.AddSeriesNote(serNote)
	response := ""
	if err != nil {
		response = "Failed to add note to database. Error: " + err.Error()
	} else {
		response = "Note successfully added"
	}

	Respond(s, i, response)
}

// Handler for note_remove_button
func NoteRemoveButtonHandler(s *discordgo.Session, i *discordgo.InteractionCreate, identifier string) {
	LogCommand(i, "note_remove_button")

	if !database.Repo.RegisteredUser(i.Member.User.ID, i.GuildID) {
		Respond(s, i, "You are not registered with this group, please get registered before removing notes.")
		return
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: "note_remove_modal " + identifier,
			Title:    "Remove Note",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "note",
							Label:       "Remove Note",
							Style:       discordgo.TextInputShort,
							Placeholder: "Input the number of the note you wish to remove here",
							Required:    true,
						},
					},
				},
			},
		},
	})
	if err != nil {
		log.Printf("Remove note modal interaction failed for guild %s and identifier %s: %s\n", i.GuildID, identifier, err.Error())
	}
}

// Handler for note_remove_modal
func NoteRemoveModalHandler(s *discordgo.Session, i *discordgo.InteractionCreate, identifier string) {
	LogCommand(i, "note_remove_modal")

	data := i.ModalSubmitData()
	note := data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value

	num, err := strconv.Atoi(note)
	if err != nil {
		Respond(s, i, "Input value could not be parsed to an integer")
		return
	}

	_, ids, err := database.Repo.GetSeriesNotes(identifier, i.GuildID)
	if err != nil {
		Respond(s, i, "Error getting notes from database: "+err.Error())
		return
	}
	if num < 1 || num > len(ids) {
		Respond(s, i, "Provided number is outside of range of notes on this series")
		return
	}
	done, err := database.Repo.RemoveSeriesNote(identifier, i.GuildID, ids[num-1])
	response := ""
	if err != nil {
		response = "Failed to remove note from database. Error: " + err.Error()
	} else if !done {
		response = "Note was not found in database"
	} else {
		response = "Note successfully removed"
	}

	Respond(s, i, response)
}
