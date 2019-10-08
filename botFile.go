package main

import (
	"./almaany"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"strconv"
	"strings"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("986995701:AAHIyuq1Nj8uc92rWYrsDhgM20zfIu6ZZRk")
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		if update.Message.Text == "/start" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "الرجاء ادخال كلمة واحدة. شكرا")
			_, _ = bot.Send(msg)
			continue
		}
		firstWord := strings.Fields(update.Message.Text)[0]
		results := almaany.ScrapePages(firstWord)
		if len(results) == 0 {
			_, _ = bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "لم يتم العثور على اي نتيجة"))
			continue
		} else {
			_, _ = bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "تم العثور على "+strconv.Itoa(len(results))))
		}
		var divided [][]almaany.Manaa

		chunkSize := 4

		for i := 0; i < len(results); i += chunkSize {
			end := i + chunkSize

			if end > len(results) {
				end = len(results)
			}

			divided = append(divided, results[i:end])
		}

		for _, element := range results {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, element.Word)
			_, _ = bot.Send(msg)

		}

	}
}
