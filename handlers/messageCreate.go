package handlers

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/cscareers-dev/cscareers-discord-v2/plugins"
)

var messageCreatePlugins map[string]plugins.Plugin = make(map[string]plugins.Plugin)

// init initializes all message create plugins
func init() {
	resumeMessageChannelPlugin := plugins.NewResumeMessageChannelPlugin()
	messageCreatePlugins[resumeMessageChannelPlugin.Name()] = resumeMessageChannelPlugin

	banUrlMessagePlugin := plugins.NewBanUrlMessagePlugin()
	messageCreatePlugins[banUrlMessagePlugin.Name()] = banUrlMessagePlugin
}

// MessageCreate processes message create events emitted from Discord API
func (h *Handler) MessageCreate(session *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.Bot || message.Author.ID == session.State.User.ID {
		return
	}

	log.Println(fmt.Sprintf("[MessageCreateHandler] incoming message from %s - message ID %s", message.Author.Username, message.ID))

	for _, messageCreatePlugin := range messageCreatePlugins {
		if !messageCreatePlugin.Enabled() {
			continue
		}

		if messageCreatePlugin.Validate(session, message) {
			log.Println(fmt.Sprintf("[%s] executing on message ID %s", messageCreatePlugin.Name(), message.ID))
			_, err := messageCreatePlugin.Execute(session, message)
			if err != nil {
				log.Fatalln(fmt.Printf("[%s] error - %s", messageCreatePlugin.Name(), err))
			}
		}
	}
}
