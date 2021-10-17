package plugins

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type ResumeMessageChannelPlugin struct{}

const (
	resumeMessageChannelPlugin = "ResumeMessageChannelPlugin"
	// TODO(corgi): update with correct channel id
	resumeMessageChannelId = "761269221001920552"
)

// NewResumeMessageChannelPlugin returns a new ResumeMessageChannelPlugin
func NewResumeMessageChannelPlugin() *ResumeMessageChannelPlugin {
	return &ResumeMessageChannelPlugin{}
}

// Name returns name of ResumeMessageChannelPlugin
func (r *ResumeMessageChannelPlugin) Name() string {
	return resumeMessageChannelPlugin
}

// Validate determines if incoming message should be executed by ResumeMessageChannelPlugin
func (r *ResumeMessageChannelPlugin) Validate(session *discordgo.Session, message *discordgo.MessageCreate) bool {
	return message.ChannelID == resumeMessageChannelId
}

// Execute runs ResumeMessageChannelPlugin on incoming message
func (r *ResumeMessageChannelPlugin) Execute(session *discordgo.Session, message *discordgo.MessageCreate) (bool, error) {
	log.Println("Executing ResumeMessageChannelPlugin on " + message.ID)
	hasPDFAttachment := false

	for _, attachment := range message.Attachments {
		if strings.HasSuffix(attachment.Filename, ".pdf") {
			hasPDFAttachment = true
			break
		}
	}

	if hasPDFAttachment {
		return true, nil
	}

	messageContent := message.Content
	session.ChannelMessageDelete(message.ChannelID, message.ID)
	privateMessageChannel, err := session.UserChannelCreate(message.Author.ID)
	if err != nil {
		return false, err
	}

	_, err = session.ChannelMessageSend(privateMessageChannel.ID, fmt.Sprintf(`
		We only allow messages that have a PDF attached in the resume channel. Your message content:
		%s
	`, messageContent))

	if err != nil {
		return false, err
	}

	return true, nil
}
