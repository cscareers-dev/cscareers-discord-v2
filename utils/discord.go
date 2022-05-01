package utils

import (
	"errors"

	"github.com/bwmarrin/discordgo"
)

// GetMessageChannelName returns the name of the channel that the message was sent to
func GetMessageChannelName(session *discordgo.Session, message *discordgo.MessageCreate) (string, error) {
	channels, err := session.GuildChannels(message.GuildID)
	if err != nil {
		return "", err
	}

	for _, channel := range channels {
		if message.ChannelID == channel.ID {
			return channel.Name, nil
		}
	}

	return "", errors.New("Unable to find channel name")
}
