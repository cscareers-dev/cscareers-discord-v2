package handlers

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/cscareers-dev/cscareers-discord-v2/plugins"
)

var messageCreatePlugins map[string]plugins.Plugin = make(map[string]plugins.Plugin)

// init initializes all message create plugins
func init() {
	resumeMessageChannelPlugin := plugins.NewResumeMessageChannelPlugin()
	messageCreatePlugins[resumeMessageChannelPlugin.Name()] = resumeMessageChannelPlugin
}

// MessageCreate processes message create events emitted from Discord API
func (h *Handler) MessageCreate(session *discordgo.Session, message *discordgo.MessageCreate) {
	log.Println("[message-create] incoming message " + message.ID)

	if message.Author.Bot || message.Author.ID == session.State.User.ID {
		return
	}

	for _, messageCreatePlugin := range messageCreatePlugins {
		if messageCreatePlugin.Validate(session, message) {
			_, err := messageCreatePlugin.Execute(session, message)
			if err != nil {
				log.Fatalln(err)
			}
		}
	}
}
