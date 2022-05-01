package plugins

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/bwmarrin/discordgo"
	discord "github.com/cscareers-dev/cscareers-discord-v2/utils"
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

var queue *sqs.SQS

// fetchCompanies fetches companies map from cscareers api and hydrates companies map
func fetchCompanies() {
	resp, err := http.Get(COMPANIES_LIST_ENDPOINT)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var companies []string
	err = json.Unmarshal(body, &companies)
	if err != nil {
		log.Fatalln(err)
	}

	for _, company := range companies {
		companiesMap[company] = true
	}
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
	session := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	queue = sqs.New(session)

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
	isProcessCommand := strings.HasPrefix(sanitizedMessage, processCommandPrefix)
	channelName, err := discord.GetMessageChannelName(session, message)
	if err != nil {
		return false
	}
	isFromProcessChannel := strings.Contains(channelName, "process")

	return isProcessCommand && isFromProcessChannel
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

	channelName, err := discord.GetMessageChannelName(session, message)
	if err != nil {
		return false, err
	}

	_, err = queue.SendMessage(&sqs.SendMessageInput{
		MessageGroupId:         aws.String("1"),
		MessageDeduplicationId: aws.String(message.ID),
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			"channel": {
				DataType:    aws.String("String"),
				StringValue: aws.String(channelName),
			},
			"company": {
				DataType:    aws.String("String"),
				StringValue: aws.String(canonicalCompanyName),
			},
			"status": {
				DataType:    aws.String("String"),
				StringValue: aws.String(strings.ToUpper(step)),
			},
			"discordId": {
				DataType:    aws.String("String"),
				StringValue: aws.String(message.Author.ID),
			},
		},
	})

	if err != nil {
		return false, err
	}

	if err = session.MessageReactionAdd(message.ChannelID, message.ID, "✅"); err != nil {
		return false, err
	}

	return true, nil
}
