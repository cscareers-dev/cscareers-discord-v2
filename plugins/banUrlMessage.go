package plugins

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type BanUrlMessagePlugin struct{}

const (
	banUrlMessagePlugin        = "BanUrlMessagePlugin"
	messageDeleteHistoryInDays = 7
)

// key = domain, value = type of url (phishing url)
var bannedUrlPatterns map[string]string = make(map[string]string)

// init initalizes bannedUrlPatterns
func init() {
	bannedUrlPatterns["discordn.gift"] = "phishing"
}

// NewBanUrlMessagePlugin returns a new BanUrlMessagePlugin
func NewBanUrlMessagePlugin() *BanUrlMessagePlugin {
	return &BanUrlMessagePlugin{}
}

// Name returns name of BanUrlMessagePlugin
func (b *BanUrlMessagePlugin) Name() string {
	return banUrlMessagePlugin
}

// Validate determines if incoming message should be executed by BanUrlMessagePlugin
func (b *BanUrlMessagePlugin) Validate(session *discordgo.Session, message *discordgo.MessageCreate) bool {
	for url := range bannedUrlPatterns {
		if strings.Contains(strings.ToLower(message.Content), url) {
			return true
		}
	}
	return false
}

// Execute runs BanUrlMessagePlugin on incoming message
func (b *BanUrlMessagePlugin) Execute(session *discordgo.Session, message *discordgo.MessageCreate) (bool, error) {
	var reason string

	for url := range bannedUrlPatterns {
		if strings.Contains(strings.ToLower(message.Content), url) {
			reason = bannedUrlPatterns[url]
		}
	}

	privateMessageChannel, err := session.UserChannelCreate(message.Author.ID)
	if err != nil {
		return false, err
	}

	_, err = session.ChannelMessageSend(privateMessageChannel.ID, fmt.Sprintf(`
		Your account has been banned from cscareers discord for posting a %s link - please secure your account and email contact@joey.dev to be unbanned
	`, reason))
	if err != nil {
		return false, err
	}

	log.Println(fmt.Sprintf("Banning %s for %s URL", message.Author.Username, reason))

	err = session.GuildBanCreateWithReason(message.GuildID, message.Author.ID, reason, messageDeleteHistoryInDays)
	if err != nil {
		return false, err
	}

	return true, nil
}
