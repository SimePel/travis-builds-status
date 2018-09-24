package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"

	"github.com/bwmarrin/discordgo"
)

// Repository
type Repository struct {
	Name string `json:"name"`
	URL  string `json:"url,omitempty"`
}

// Travis
type Travis struct {
	State      string     `json:"state"`
	BuildURL   string     `json:"build_url"`
	Message    string     `json:"message"`
	AuthorName string     `json:"author_name"`
	Repo       Repository `json:"repository"`
}

func main() {
	token := os.Getenv("BOT_TOKEN")
	bot, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal(err)
	}

	webhookID := os.Getenv("WEBHOOK_ID")
	wb, err := bot.Webhook(webhookID)
	if err != nil {
		log.Fatal(err)
	}

	ln, err := net.Listen("tcp", ":80")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handle(conn, bot, wb)
	}
}

func handle(c net.Conn, bot *discordgo.Session, wb *discordgo.Webhook) {
	b, err := ioutil.ReadAll(c)
	if err != nil {
		log.Print("ioutil.ReadAll:", err)
		return
	}

	var t *Travis
	err = json.Unmarshal(b, &t)
	if err != nil {
		log.Print("json.Unmarshal:", err)
		return
	}

	d := discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{{
			URL:         t.BuildURL,
			Title:       fmt.Sprintf("Status of your repo - %s", t.Repo.Name),
			Description: fmt.Sprintf("Author: %s\nStatus - %s\nLatest commit: %s", t.AuthorName, t.State, t.Message),
		},
		},
	}
	err = bot.WebhookExecute(wb.ID, wb.Token, false, &d)
	if err != nil {
		log.Print(err)
	}

	c.Close()
}
