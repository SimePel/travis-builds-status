package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"

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
	router := httprouter.New()
	router.GET("/", index)
	router.POST("/", travis)

	log.Fatal(http.ListenAndServe(":80", router))
}

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "There is nothing interesting!")
}

func travis(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Print("ioutil.ReadAll:", err)
		return
	}
	r.Body.Close()

	var t *Travis
	err = json.Unmarshal(b, &t)
	if err != nil {
		log.Print("json.Unmarshal:", err)
		return
	}

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
}
