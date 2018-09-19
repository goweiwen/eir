package eir

import (
	"fmt"
	"log"
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

	b, err := tb.NewBot(tb.Settings{
		Token:  token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle(tb.OnAddedToGroup, func(m *tb.Message) {
		b.Send(m.Chat, fmt.Sprintf("Hello, your chat ID is %d", m.Chat.ID))
	})

	b.Handle("/start", func(m *tb.Message) {
		b.Send(m.Chat, fmt.Sprintf("Hello, your chat ID is %d", m.Chat.ID))
	})

	scheduleJobs(func(msg string) {
		log.Printf(msg)
		b.Send(chat, msg, tb.ModeMarkdown)
	})

	log.Printf("Starting bot...")
	b.Start()
}

func scheduleJobs(say func(string)) {
	c := cron.New()

	c.AddFunc(
		"0 55 21 * * WED",
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
		`)
		},
	)

	c.AddFunc(
		"0 0 23 * * TUE",
		func() {
			now := time.Now()
			_, week := now.ISOWeek()
			isRecyclingWeek := week%2 == 0
			if isRecyclingWeek {
				say("*💩♻ Garbage & Recycling day!*")
			} else {
				say("*💩 Garbage day!*")
			}
		},
	)

	c.AddFunc(
		"0 0 7 25 * *",
		func() {
			say("*🏠 Remember to pay the rent and utilities!")
		},
	)

	c.Start()
	log.Printf("Scheduled jobs!")
}
