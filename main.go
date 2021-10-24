package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/cscareers-dev/cscareers-discord-v2/config"
	discordHandlers "github.com/cscareers-dev/cscareers-discord-v2/handlers"
)

var session *discordgo.Session

func init() {
	var err error
	config := config.Load()
	session, err = discordgo.New("Bot " + config.Token)

	if err != nil {
		panic("failed to initalize bot")
	}

	handlers := discordHandlers.New()
	session.AddHandler(handlers.MessageCreate)
	session.AddHandler(handlers.Ready)

	err = session.Open()
	if err != nil {
		panic(err)
	}
}

func main() {
	fmt.Println("Bot is now running. Press CTRL-C to exit.")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGTERM)
	<-stop

	// clean up
	session.Close()
}
