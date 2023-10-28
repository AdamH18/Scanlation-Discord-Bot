package bot

import (
	"errors"
	"log"
	"scanlation-discord-bot/database"
	"sort"
	"strings"

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

// Handles response to a slash command, non ephemeral
func RespondNonEph(s *discordgo.Session, i *discordgo.InteractionCreate, response string) {
	log.Printf("Response: %s\n", response)
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: response,
		},
	})
}

// Handles embed response to a slash command
func RespondEmbed(s *discordgo.Session, i *discordgo.InteractionCreate, response *discordgo.MessageEmbed) {
	log.Printf("Response: %v\n", response)
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{response},
			Flags:  discordgo.MessageFlagsEphemeral,
		},
	})
}

// Standardized logging of command calls
func LogCommand(i *discordgo.InteractionCreate, name string) {
	var g1, g2, c1, c2 string
	gld, err := GetGuildName(i.GuildID)
	if err != nil {
		g1 = i.GuildID
		g2 = "Error fetching name: " + err.Error()
	} else {
		g1 = gld
		g2 = i.GuildID
	}
	chn, err := GetChannelName(i.ChannelID)
	if err != nil {
		c1 = i.ChannelID
		c2 = "Error fetching name: " + err.Error()
	} else {
		c1 = chn
		c2 = i.ChannelID
	}
	log.Printf("User %s (%s) in guild %s (%s) and channel %s (%s) used %s command with options:\n", i.Member.User.Username, i.Member.User.ID, g1, g2, c1, c2, name)
}

// Checks if command registered with Discord is in bot
func IsCommand(n string) bool {
	for _, cmd := range commands {
		if cmd.Name == n {
			return true
		}
	}
	return false
}

// Checks if command in bot is registered with Discord
func DiscordCommand(cmds []*discordgo.ApplicationCommand, n string) bool {
	for _, cmd := range cmds {
		if cmd.Name == n {
			return true
		}
	}
	return false
}

// Creates new channels for a series
func CreateChannels(ser database.Series) error {
	//Get registered bounds for series channels
	top, bottom, err := database.Repo.GetSeriesChannels(ser.Guild)
	if err != nil {
		return err
	}

	//Get all channels in guild
	channels, err := goBot.GuildChannels(ser.Guild)
	if err != nil {
		return err
	}

	//Filter down to only categories
	categories := []*discordgo.Channel{}
	for _, channel := range channels {
		if channel.Type == 4 {
			categories = append(categories, channel)
		}
	}

	//Sort categories based on position in server
	sort.Slice(categories, func(i, j int) bool {
		return categories[i].Position < categories[j].Position
	})

	//Determine position of top and bottom of series channels
	topInd := -1
	botInd := -1
	for i, cat := range categories {
		if cat.ID == top {
			topInd = i
		}
		if cat.ID == bottom {
			botInd = i
		}
	}
	if botInd < topInd || topInd < 0 || botInd < 0 {
		return errors.New("series channel bounds improperly defined")
	}
	seriesCats := categories[topInd : botInd+1]

	//Determine insertion location
	loc := -1
	newStart := false
	newEnd := false
	for i, existing := range seriesCats {
		if strings.ToLower(ser.NameSh) < strings.ToLower(existing.Name) && loc < 0 {
			loc = existing.Position
			if i == 0 {
				newStart = true
			}
		}
		//Increase position value of all categories after new series
		if loc >= 0 {
			categories[topInd+i].Position++
		}
	}
	if loc < 0 {
		loc = seriesCats[len(seriesCats)-1].Position + 1
		newEnd = true
	}
	//Increase position value of all categories after series categories
	for i := range categories[botInd+1:] {
		categories[botInd+1+i].Position++
	}

	//Member role visible, everyone invisible perms
	mem := database.Repo.GetMemberRole(ser.Guild)
	perms := []*discordgo.PermissionOverwrite{}
	if mem != "" {
		perms = []*discordgo.PermissionOverwrite{
			{
				ID:   ser.Guild,
				Type: discordgo.PermissionOverwriteTypeRole,
				Deny: (1 << 10),
			},
			{
				ID:    mem,
				Type:  discordgo.PermissionOverwriteTypeRole,
				Allow: (1 << 10),
			},
		}
	}

	//Create the new channels
	catData := discordgo.GuildChannelCreateData{
		Name:                 ser.NameSh,
		Type:                 discordgo.ChannelTypeGuildCategory,
		Position:             loc,
		PermissionOverwrites: perms,
	}
	newCat, err := goBot.GuildChannelCreateComplex(ser.Guild, catData)
	if err != nil {
		return err
	}
	genData := discordgo.GuildChannelCreateData{
		Name:                 ser.NameSh,
		Type:                 discordgo.ChannelTypeGuildText,
		Position:             2,
		PermissionOverwrites: perms,
		ParentID:             newCat.ID,
	}
	pr, err := goBot.GuildChannelCreateComplex(ser.Guild, genData)
	if err != nil {
		return err
	}
	prData := discordgo.GuildChannelCreateData{
		Name:                 ser.NameSh + "-pr",
		Type:                 discordgo.ChannelTypeGuildText,
		Position:             3,
		PermissionOverwrites: perms,
		ParentID:             newCat.ID,
	}
	gen, err := goBot.GuildChannelCreateComplex(ser.Guild, prData)
	if err != nil {
		return err
	}
	if len(perms) > 0 {
		perms[0].Deny = (1 << 10) | (1 << 11)
	} else {
		perms = append(perms, &discordgo.PermissionOverwrite{
			ID:   ser.Guild,
			Type: 0,
			Deny: (1 << 11),
		})
	}
	infData := discordgo.GuildChannelCreateData{
		Name:                 ser.NameSh + "-info",
		Type:                 discordgo.ChannelTypeGuildText,
		Position:             1,
		PermissionOverwrites: perms,
		ParentID:             newCat.ID,
	}
	inf, err := goBot.GuildChannelCreateComplex(ser.Guild, infData)
	if err != nil {
		return err
	}

	//Reorder categories after all channels have been created
	categories = append(categories, newCat)
	goBot.GuildChannelsReorder(ser.Guild, categories)

	//Register channels as belonging to series
	go database.Repo.AddChannel(database.Channel{Channel: inf.ID, Series: ser.NameSh, Guild: ser.Guild})
	go database.Repo.AddChannel(database.Channel{Channel: gen.ID, Series: ser.NameSh, Guild: ser.Guild})
	go database.Repo.AddChannel(database.Channel{Channel: pr.ID, Series: ser.NameSh, Guild: ser.Guild})

	//Update database with new top or bottom of channel bounds if necessary
	if newStart {
		go database.Repo.UpdateSeriesChannelsTop(newCat.ID, ser.Guild)
	}
	if newEnd {
		go database.Repo.UpdateSeriesChannelsBottom(newCat.ID, ser.Guild)
	}

	return nil
}

// Creates pingable role for a series
func CreatePingRole(ser database.Series) (string, error) {
	params := discordgo.RoleParams{
		Name: ser.NameSh,
	}
	role, err := goBot.GuildRoleCreate(ser.Guild, &params)
	if err != nil {
		return "", err
	}
	return role.ID, nil
}
