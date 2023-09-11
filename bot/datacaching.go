package bot

import (
	"errors"
	"hash/fnv"
)

var (
	Usernames    map[KeyStruct]string
	UserPings    map[KeyStruct]string
	RolePings    map[KeyStruct]string
	ChannelPings map[string]string
	ChannelNames map[string]string
	GuildNames   map[string]string
)

// Needed to use guild-user/role IDs together for map keys
type KeyStruct struct {
	Guild string
	Data  string
}

func (k KeyStruct) Equals(other KeyStruct) bool {
	return k.Guild == other.Guild && k.Data == other.Data
}

func (k KeyStruct) HashCode() uint32 {
	h := fnv.New32a()
	h.Write([]byte(k.Guild))
	h.Write([]byte(k.Data))
	return h.Sum32()
}

func InitializeCache() {
	Usernames = make(map[KeyStruct]string)
	UserPings = make(map[KeyStruct]string)
	RolePings = make(map[KeyStruct]string)
	ChannelPings = make(map[string]string)
	ChannelNames = make(map[string]string)
	GuildNames = make(map[string]string)
}

// Takes a user ID and returns the username
func GetUserName(guildID string, userID string) (string, error) {
	//Check if cached before asking Discord
	key := KeyStruct{
		Guild: guildID,
		Data:  userID,
	}
	if _, ok := Usernames[key]; ok {
		return Usernames[key], nil
	}
	usr, err := goBot.GuildMember(guildID, userID)
	if err != nil {
		return "", err
	}
	//Cache returned name
	Usernames[KeyStruct{Guild: guildID, Data: userID}] = usr.User.Username
	UserPings[KeyStruct{Guild: guildID, Data: userID}] = usr.Mention()
	return usr.User.Username, nil
}

// Takes a user ID and returns a ping string
func GetUserPing(guildID string, userID string) (string, error) {
	//Check if cached before asking Discord
	key := KeyStruct{
		Guild: guildID,
		Data:  userID,
	}
	if _, ok := UserPings[key]; ok {
		return UserPings[key], nil
	}
	usr, err := goBot.GuildMember(guildID, userID)
	if err != nil {
		return "", err
	}
	//Cache returned ping
	UserPings[KeyStruct{Guild: guildID, Data: userID}] = usr.Mention()
	Usernames[KeyStruct{Guild: guildID, Data: userID}] = usr.User.Username
	return usr.Mention(), nil
}

// Takes a role ID and returns a ping string
func GetRolePing(guildID string, roleID string) (string, error) {
	//Check if cached before asking Discord
	key := KeyStruct{
		Guild: guildID,
		Data:  roleID,
	}
	if _, ok := RolePings[key]; ok {
		return RolePings[key], nil
	}
	roles, err := goBot.GuildRoles(guildID)
	if err != nil {
		return "", err
	}
	for _, role := range roles {
		if role.ID == roleID {
			//Cache returned role
			UserPings[KeyStruct{Guild: guildID, Data: roleID}] = role.Mention()
			return role.Mention(), nil
		}
	}
	return "", errors.New("role not found")
}

// Takes a channel ID and returns a ping string
func GetChannelPing(channelID string) (string, error) {
	//Check if cached before asking Discord
	if _, ok := ChannelPings[channelID]; ok {
		return ChannelPings[channelID], nil
	}
	channel, err := goBot.Channel(channelID)
	if err != nil {
		return "", err
	}
	//Cache returned channel
	ChannelPings[channelID] = channel.Mention()
	ChannelNames[channelID] = channel.Name
	return channel.Mention(), nil
}

// Takes a channel ID and returns a ping string
func GetGuildName(guildID string) (string, error) {
	//Check if cached before asking Discord
	if _, ok := GuildNames[guildID]; ok {
		return GuildNames[guildID], nil
	}
	guild, err := goBot.Guild(guildID)
	if err != nil {
		return "", err
	}
	//Cache returned guild
	GuildNames[guildID] = guild.Name
	return guild.Name, nil
}

// Takes a channel ID and returns a ping string
func GetChannelName(channelID string) (string, error) {
	//Check if cached before asking Discord
	if _, ok := ChannelNames[channelID]; ok {
		return ChannelNames[channelID], nil
	}
	channel, err := goBot.Channel(channelID)
	if err != nil {
		return "", err
	}
	//Cache returned channel
	ChannelNames[channelID] = channel.Name
	ChannelPings[channelID] = channel.Mention()
	return channel.Name, nil
}
