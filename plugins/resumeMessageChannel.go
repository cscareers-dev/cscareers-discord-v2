package plugins

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type ResumeMessageChannelPlugin struct{}

const (
	resumeMessageChannelPlugin = "ResumeMessageChannelPlugin"
	resumeMessageChannelId     = "699532322948513832"
)

// NewResumeMessageChannelPlugin returns a new ResumeMessageChannelPlugin
func NewResumeMessageChannelPlugin() *ResumeMessageChannelPlugin {
	return &ResumeMessageChannelPlugin{}
}

// Name returns name of ResumeMessageChannelPlugin
func (r *ResumeMessageChannelPlugin) Name() string {
	return resumeMessageChannelPlugin
}

// Enabled returns if ResumeMessageChannelPlugin is enabled
func (b *ResumeMessageChannelPlugin) Enabled() bool {
	// This should be enabled once discordgo adds thread functionality
	// See https://github.com/bwmarrin/discordgo/pull/1058
	return false
}

// Validate determines if incoming message should be executed by ResumeMessageChannelPlugin
func (r *ResumeMessageChannelPlugin) Validate(session *discordgo.Session, message *discordgo.MessageCreate) bool {
	return message.ChannelID == resumeMessageChannelId
}

// Execute runs ResumeMessageChannelPlugin on incoming message
func (r *ResumeMessageChannelPlugin) Execute(session *discordgo.Session, message *discordgo.MessageCreate) (bool, error) {
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
	err := session.ChannelMessageDelete(message.ChannelID, message.ID)
	if err != nil {
		return false, err
	}

	privateMessageChannel, err := session.UserChannelCreate(message.Author.ID)
	if err != nil {
		return false, err
	}

	_, err = session.ChannelMessageSend(privateMessageChannel.ID, fmt.Sprintf(`
		We only allow top level messages that have a PDF attached in the resume channel. Your message content:
		%s
	`, messageContent))

	if err != nil {
		return false, err
	}

	return true, nil
}
