package plugins

import "github.com/bwmarrin/discordgo"

// Plugin interface is blueprint for a Discord plugin
type Plugin interface {
	Name() string
	Validate(session *discordgo.Session, message *discordgo.MessageCreate) bool
	Execute(session *discordgo.Session, message *discordgo.MessageCreate) (bool, error)
}
