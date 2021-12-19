package plugins

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type SalaryCommandPlugin struct{}

var salaryUrls map[string]string

const (
	salaryCommandPlugin = "SalaryCommandPlugin"
	salaryCommandPrefix = "!salary"
)

const (
	SALARIES_ENDPOINT = "https://cscareers.dev/api/discord/getSalaryUrls"
)

// init hydrates salaryUrls with salary url information
func init() {
	resp, err := http.Get(SALARIES_ENDPOINT)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	err = json.Unmarshal(body, &salaryUrls)
	if err != nil {
		log.Fatalln(err)
	}
}

// NewSalaryCommandPlugin returns a new ResumeMessageChannelPlugin
func NewSalaryCommandPlugin() *SalaryCommandPlugin {
	return &SalaryCommandPlugin{}
}

// Name returns name of ResumeMessageChannelPlugin
func (s *SalaryCommandPlugin) Name() string {
	return salaryCommandPlugin
}

// Enabled returns if SalaryCommandPlugin is enabled
func (s *SalaryCommandPlugin) Enabled() bool {
	return true
}

// Validate determines if incoming message should be executed by SalaryCommandPlugin
func (s *SalaryCommandPlugin) Validate(session *discordgo.Session, message *discordgo.MessageCreate) bool {
	sanitizedMessage := strings.ToLower(message.Content)
	return strings.HasPrefix(sanitizedMessage, salaryCommandPrefix)
}

// Execute runs SalaryCommandPlugin on incoming message
func (s *SalaryCommandPlugin) Execute(session *discordgo.Session, message *discordgo.MessageCreate) (bool, error) {
	sanitizedMessage := strings.Replace(strings.ToLower(message.Content), salaryCommandPrefix, "", 1)
	sanitizedMessage = strings.TrimSpace(sanitizedMessage)
	companyName := strings.Title(sanitizedMessage)
	url, ok := salaryUrls[companyName]
	if !ok {
		session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Unable to locate company salary info on %s", companyName))
		return true, nil
	}

	session.ChannelMessageSend(message.ChannelID, url)
	return true, nil
}
