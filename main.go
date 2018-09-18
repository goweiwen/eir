package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/jasonlvhit/gocron"
	tb "gopkg.in/tucnak/telebot.v2"
)

func main() {
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

	b.Send(chat, "Hello!")

	gocron.Every(1).Wednesday().At("22:00").Do(NotifyCleaningDay, b, chat)
	gocron.Every(1).Tuesday().At("22:30").Do(NotifyGarbageDay, b, chat)

	log.Printf("Starting bot...")

	b.Start()
}

func NotifyCleaningDay(b *tb.Bot, chat tb.Recipient) {
	log.Printf("Cleaning day!")
	b.Send(chat, "Time to clean the house!")
}

func NotifyGarbageDay(b *tb.Bot, chat tb.Recipient) {
	log.Printf("Garbage day!")
	b.Send(chat, "Remember to take out the trash tonight!")
}
