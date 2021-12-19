package handlers

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

// Ready processes ready events emitted from Discord API
func (h *Handler) Ready(session *discordgo.Session, _ready *discordgo.Ready) {
	log.Println("[ReadyHandler] ready")
}
