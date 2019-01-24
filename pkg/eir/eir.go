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

	error := func(msg string) {
		log.Printf("[%d] ERROR: %s", chat.ID, msg)
		b.Send(chat, fmt.Sprintf("ERROR: ```%s```", msg), tb.ModeMarkdown)
	}

	b.Handle(tb.OnAddedToGroup, func(m *tb.Message) {
		msg := fmt.Sprintf("Hello, your chat ID is %d", m.Chat.ID)
		log.Printf("[%d] @sawmillbot: %s", chat.ID, msg)
		b.Send(m.Chat, msg)
	})

	b.Handle("/start", func(m *tb.Message) {
		msg := fmt.Sprintf("Hello, your chat ID is %d", m.Chat.ID)
		log.Printf("[%d] @sawmillbot: %s", chat.ID, msg)
		b.Send(m.Chat, msg)
	})

	b.Handle("/weather", func(m *tb.Message) {
		weather, err := fetchWeather()
		if err != nil {
			error(err.Error())
			return
		}
		say(fmt.Sprintf("```\n%s\n```", weather))
	})

	b.Handle("/suckmydick", func(m *tb.Message) {
		pedo := os.Getenv("PEDO_TELEGRAM_USERNAME")
		var msg string
		if m.Sender.Username == pedo {
			msg = "Cannot be found"
		} else {
			msg = fmt.Sprintf("@%s, please suck @%s's dick", pedo, m.Sender.Username)
		}
		log.Printf("[%d] @sawmillbot: %s", chat.ID, msg)
		b.Send(m.Chat, msg)
	})

	scheduleJobs(say, error)

	log.Printf("Starting bot...")
	b.Start()
}

func scheduleJobs(say func(string), error func(string)) {
	c := cron.New()

	c.AddFunc(
		"0 0 9 * * *",
		func() {
			weather, err := fetchWeather()
			if err != nil {
				error(err.Error())
				return
			}
			say(fmt.Sprintf("Good morning!\n```\n%s\n```", weather))
		},
	)

	c.AddFunc(
		"0 55 21 * * TUE",
		func() {
			say(`*✨ House Cleaning*
- Organize the fridge
- Clean and wipe the sink
- Wipe down the stove
- Wipe the tables
- Mop the kitchen floor
- Vacuum the living room
- Sweep the backyard
- Put shoes neatly
- Take out the trash
		`)
		},
	)

	c.AddFunc(
		"0 0 18 * * WED",
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

	c.AddFunc(
		"0 0 7 25,26,27 * *",
		func() {
			say("*🏠 Remember to pay the rent and utilities!*")
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
