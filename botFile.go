package main

import (
	"./almaany"
	"encoding/json"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	strip "github.com/grokify/html-strip-tags-go"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

type Config struct {
	BotUsername string
	Dev         string
}

func handleTextUpdates(bot *tgbotapi.BotAPI, update tgbotapi.Update) {

	firstWord := strings.Fields(update.Message.Text)[0]
	dbResults := almaany.GetSearchedWord(firstWord)
	if len(dbResults) == 0 {
		results := almaany.ScrapePages(firstWord)
		if len(results) == 0 {
			_, _ = bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "لم يتم العثور على اي نتيجة"))
			return
		} else {
			almaany.SaveWords(firstWord, results)
			dbResults = almaany.GetSearchedWord(firstWord)
		}
	}
	if len(dbResults) == 0 {
		_, _ = bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "لم يتم العثور على اي نتيجة"))
		return
	}

	var divided [][]tgbotapi.InlineKeyboardButton

	chunkSize := 2

	for i := 0; i < len(dbResults); i += chunkSize {
		var temp []tgbotapi.InlineKeyboardButton
		to := chunkSize
		if len(dbResults)-i < chunkSize {
			to = len(dbResults) - i
		}
		for j := 0; j < to; j += 1 {

			temp = append(temp, tgbotapi.NewInlineKeyboardButtonData(dbResults[j+i], dbResults[j+i]))
		}

		divided = append(divided, temp)
	}

	numericKeyboard := tgbotapi.NewInlineKeyboardMarkup(divided...)

	text := "تم العثور على " + strconv.Itoa(len(dbResults)) + ".\nالرجاء النقر على الكلمة التي تريدها"
	text += "\n..........................................................................." +
		"..............................................................................." +
		"................................................................................"
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	msg.ReplyMarkup = numericKeyboard
	_, _ = bot.Send(msg)

}

const start = `حياكم الله في قاموس المعاني

💠 ابحث عن أي كلمة عربية في التلغرام لتعرف عن معناها و تفاسيرها.

للمساعدة وطريقة الاستخدام اضغط /help`

const help = `📚 المعاني - عربي عربي

💠 يمكنك البحث بطريقتين:

🔍 داخل البوت:
اكتب الكلمة في المحادثة الخاصة  مع البوت وسيعطيك البوت أقرب النتائج لكلمتك، تستطيع الضغط على أي منها لمعرفة تفاسيرها.

🔍 من خارج البوت في أي محادثة:
في أي محادثة يمكنك كتابة اسم البوت في صندوق الكتابة مبتوعا بكلمة البحث مع مسافة فاصلة بينهما هكذا:

@%s كلمة

سيعطيك البوت الكثير من النتائج وتستطيع اختيار احداها لإرسالها للشخص الذي تتحدث معه.

للتواصل مع المطور في حال حدوث خلل مع البوت
@%s`

func main() {
	almaany.InitDatabase()
	var config Config
	data, err := ioutil.ReadFile("./configs.json")
	if err != nil {
		fmt.Print(err)
	}

	_ = json.Unmarshal(data, &config)
	help := fmt.Sprintf(help, config.BotUsername, config.Dev)

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

		if update.Message != nil { // handles text updates
			if update.Message.Chat.ID > 0 {
				if update.Message.Text == "/start" {
					msg := tgbotapi.NewMessage(int64(update.Message.From.ID), start)
					_, _ = bot.Send(msg)
					almaany.AddUser(update.Message.From)
				} else if update.Message.Text == "/help" {
					msg := tgbotapi.NewMessage(int64(update.Message.From.ID), help)
					_, _ = bot.Send(msg)
				} else {
					handleTextUpdates(bot, update)
				}

			}
		} else if update.CallbackQuery != nil { // handles inline callbacks
			handleCallbackQueryUpdates(bot, update)
		} else if update.InlineQuery != nil {
			handleInlineQueryUpdates(bot, update)

		}
	}
}

func handleInlineQueryUpdates(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if len(update.InlineQuery.Query) == 0 {
		inline := tgbotapi.NewInlineQueryResultArticleHTML(update.InlineQuery.ID, "إكتب كلمة ليتم البحث عن معناها",
			"الرجاء كتابة كلمة ليتم البحث عن معناها")
		inlineConf := tgbotapi.InlineConfig{
			InlineQueryID: update.InlineQuery.ID,
			IsPersonal:    false,
			CacheTime:     0,
			Results:       []interface{}{inline},
		}
		_, _ = bot.AnswerInlineQuery(inlineConf)
		return

	}
	firstWord := strings.Fields(update.InlineQuery.Query)[0]
	if len(firstWord) == 0 {
		inline := tgbotapi.NewInlineQueryResultArticleHTML(update.InlineQuery.ID, "إكتب كلمة ليتم البحث عن معناها",
			"الرجاء كتابة كلمة ليتم البحث عن معناها")
		inlineConf := tgbotapi.InlineConfig{
			InlineQueryID: update.InlineQuery.ID,
			IsPersonal:    false,
			CacheTime:     0,
			Results:       []interface{}{inline},
		}
		_, _ = bot.AnswerInlineQuery(inlineConf)

		return
	}
	dbResults := almaany.GetSearchedWord(firstWord)
	if len(dbResults) == 0 {
		results := almaany.ScrapePages(firstWord)
		if len(results) == 0 {
			inline := tgbotapi.NewInlineQueryResultArticleHTML(update.InlineQuery.ID, "لم يتم العثور على الكلمة", "لم أقدر على العثور"+
				" على الكلمة "+firstWord+" الرجاء البحث عن كلمة أخرى")
			inlineConf := tgbotapi.InlineConfig{
				InlineQueryID: update.InlineQuery.ID,
				IsPersonal:    false,
				CacheTime:     0,
				Results:       []interface{}{inline},
			}
			_, _ = bot.AnswerInlineQuery(inlineConf)
			return
		} else {
			almaany.SaveWords(firstWord, results)
			dbResults = almaany.GetSearchedWord(firstWord)
		}
	}
	if len(dbResults) == 0 {
		inline := tgbotapi.NewInlineQueryResultArticleHTML(update.InlineQuery.ID, "لم يتم العثور على الكلمة",
			"لم أقدر على العثور"+" على الكلمة "+firstWord+" الرجاء البحث عن كلمة أخرى")
		inlineConf := tgbotapi.InlineConfig{
			InlineQueryID: update.InlineQuery.ID,
			IsPersonal:    false,
			CacheTime:     0,
			Results:       []interface{}{inline},
		}
		_, _ = bot.AnswerInlineQuery(inlineConf)
		return
	}
	var results []interface{}

	for index, result := range dbResults {
		explanation := almaany.GetExplanation(result)
		article := tgbotapi.NewInlineQueryResultArticleHTML(strconv.Itoa(index), result, almaany.FormatMaany(explanation))
		article.Description = strip.StripTags(explanation.Explanations[0])
		results = append(results, article)
	}

	inline := tgbotapi.InlineConfig{
		InlineQueryID: update.InlineQuery.ID,
		IsPersonal:    true,
		CacheTime:     0,
		Results:       results,
	}

	bot.AnswerInlineQuery(inline)

}

func handleCallbackQueryUpdates(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	_, _ = bot.AnswerCallbackQuery(
		tgbotapi.NewCallback(update.CallbackQuery.ID,
			"تم إرسال معنى : "+update.CallbackQuery.Data))
	query := update.CallbackQuery.Data
	from := update.CallbackQuery.From.ID
	res := almaany.GetExplanation(query)
	formattedString := almaany.FormatMaany(res)
	msg := tgbotapi.NewMessage(int64(from), formattedString)
	msg.ParseMode = "html"
	_, _ = bot.Send(msg)

}
