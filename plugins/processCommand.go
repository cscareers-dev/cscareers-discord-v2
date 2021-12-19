package plugins

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type ProcessCommandPlugin struct{}

const (
	processCommandPlugin = "ProcessCommandPlugin"
	processCommandPrefix = "!process"
	processCommandUsage  = "!process Company Apply|Reject|OA|Phone|Final|Offer"
)

const (
	process2022NewGradChannelID      = "855856926235951124"
	process2022SummerInternChannelID = "856027992701927444"
)

const (
	COMPANIES_CACHE_SECONDS_TTL = 3600 // 1 hour
	COMPANIES_LIST_ENDPOINT     = "https://www.cscareers.dev/api/_utils/getProcessTrackingCompanies"
)

var stepsMap map[string]bool
var companiesMap map[string]bool
var _lastHydrateTime time.Time

// fetchCompanies fetches companies map from cscareers api and hydrates companies map
func fetchCompanies() {

}

// shouldFetchCompanies determines if the companies map cache is stale
func shouldFetchCompanies() bool {
	if _lastHydrateTime.IsZero() {
		return true
	}

	diff := time.Now().Sub(_lastHydrateTime)
	return diff.Seconds() >= COMPANIES_CACHE_SECONDS_TTL
}

// init sets up valid steps and hydrates companies map
func init() {
	// hydrate steps map
	stepsMap["oa"] = true
	stepsMap["reject"] = true
	stepsMap["phone"] = true
	stepsMap["final"] = true
	stepsMap["offer"] = true

	fetchCompanies()
}

// NewProcessCommandPlugin returns a new ProcessCommandPlugin
func NewProcessCommandPlugin() *ProcessCommandPlugin {
	return &ProcessCommandPlugin{}
}

// Name returns name of ProcessCommandPlugin
func (p *ProcessCommandPlugin) Name() string {
	return processCommandPlugin
}

// Enabled returns if ProcessCommandPlugin is enabled
func (p *ProcessCommandPlugin) Enabled() bool {
	return true
}

// Validate determines if incoming message should be executed by ProcessCommandPlugin
func (p *ProcessCommandPlugin) Validate(session *discordgo.Session, message *discordgo.MessageCreate) bool {
	sanitizedMessage := strings.ToLower(message.Content)
	return strings.HasPrefix(sanitizedMessage, processCommandPrefix)
}

// Execute runs ProcessCommandPlugin on incoming message
func (p *ProcessCommandPlugin) Execute(session *discordgo.Session, message *discordgo.MessageCreate) (bool, error) {
	sanitizedMessage := strings.Replace(strings.ToLower(message.Content), processCommandPrefix, "", 1)
	segments := strings.Split(sanitizedMessage, " ")
	if len(segments) < 2 {
		_, err := session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Invalid usage: %s", processCommandUsage))
		return false, err
	}

	step := segments[0]
	if _, valid := stepsMap[step]; !valid {
		_, err := session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Invalid step: %s", processCommandUsage))
		return false, err
	}

	if shouldFetchCompanies() {
		fetchCompanies()
	}

	canonicalCompanyName := strings.Join(segments[1:], " ")
	if _, valid := companiesMap[canonicalCompanyName]; !valid {
		_, err := session.ChannelMessageSend(message.ChannelID, "Company is not recognized. Please submit a company request <https://cscareers.dev/company/add>")
		return false, err
	}

	// TODO:
	// - submit to queue
	// - react to sender's message

	return true, nil
}
