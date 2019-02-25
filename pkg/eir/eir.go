package eir

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/robfig/cron"
	tb "gopkg.in/tucnak/telebot.v2"
)

func Start() {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	chatID, err := strconv.Atoi(os.Getenv("TELEGRAM_BOT_CHAT_ID"))
	if err != nil {
		log.Fatal("TELEGRAM_BOT_CHAT_ID is invalid")
	}
	chat := &tb.User{ID: chatID}

	poller := &tb.LongPoller{Timeout: 10 * time.Second}
	logger := tb.NewMiddlewarePoller(poller, func(upd *tb.Update) bool {
		if upd.Message != nil {
			log.Printf("[%d] @%s: %s", upd.Message.Chat.ID, upd.Message.Sender.Username, upd.Message.Text)
		}
		return true
	})

	b, err := tb.NewBot(tb.Settings{
		Token:  token,
		Poller: logger,
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	say := func(msg string) {
		log.Printf("[%d] @sawmillbot: %s", chat.ID, msg)
		b.Send(chat, msg, tb.ModeMarkdown)
	}

	reply := func(chat *tb.Chat, msg string) {
		log.Printf("[%d] @sawmillbot: %s", chat.ID, msg)
		b.Send(chat, msg, tb.ModeMarkdown)
	}

	error := func(msg string) {
		log.Printf("[%d] ERROR: %s", chat.ID, msg)
		b.Send(chat, fmt.Sprintf("ERROR: ```%s```", msg), tb.ModeMarkdown)
	}

	b.Handle(tb.OnAddedToGroup, func(m *tb.Message) {
		msg := fmt.Sprintf("Hello, your chat ID is %d", m.Chat.ID)
		log.Printf("[%d] @sawmillbot: %s", m.Chat.ID, msg)
		reply(m.Chat, msg)
	})

	b.Handle("/start", func(m *tb.Message) {
		msg := fmt.Sprintf("Hello, your chat ID is %d", m.Chat.ID)
		log.Printf("[%d] @sawmillbot: %s", m.Chat.ID, msg)
		reply(m.Chat, msg)
	})

	b.Handle("/weather", func(m *tb.Message) {
		weather, err := fetchWeather()
		if err != nil {
			error(err.Error())
			return
		}
		reply(m.Chat, fmt.Sprintf("```\n%s\n```", weather))
	})

	b.Handle("/suckmydick", func(m *tb.Message) {
		pedo := os.Getenv("PEDO_TELEGRAM_USERNAME")
		var msg string
		if m.Sender.Username == pedo {
			msg = "Cannot be found"
		} else {
			msg = fmt.Sprintf("@%s, please suck @%s's dick", pedo, m.Sender.Username)
		}
		log.Printf("[%d] @sawmillbot: %s", m.Chat.ID, msg)
		reply(m.Chat, msg)
	})

	scheduleJobs(say, error)

	log.Printf("Starting bot...")
	b.Start()
}

func scheduleJobs(say func(string), error func(string)) {
	c := cron.New()

	c.AddFunc(
		"0 0 22 * * WED",
		func() {
			now := time.Now()
			_, week := now.ISOWeek()
			isRecyclingWeek := week%2 == 0
			if isRecyclingWeek {
				say("*💩 Garbage & ♻️ Recycling day!*")
			} else {
				say("*💩 Garbage day!*")
			}
		},
	)

	c.Start()
	log.Printf("Scheduled jobs!")
}

func fetchWeather() (string, error) {
	coordinates := os.Getenv("COORDINATES")

	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://wttr.in/%s?m0QT", coordinates), nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("User-Agent", "curl")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("HTTP status not OK")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
