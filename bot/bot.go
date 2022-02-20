package bot

import (
	"TelegramBot/utils"
	"errors"
	"fmt"
	"os"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

const (
	TOKEN = "5117279371:AAHKyOdJOF0VIIkkMrPHioabMbtosxgf9AM"
)

type TelegramBot struct {
	bot         *tgbotapi.BotAPI
	userUpdates map[int64]chan tgbotapi.Update
	updates     tgbotapi.UpdatesChannel
	db          *utils.MongoDB
}

func NewTelegramBot(db *utils.MongoDB) (*TelegramBot, error) {
	token, exists := os.LookupEnv("TOKEN")

	if !exists {
		return nil, errors.New("TOKEN not found in file .env")
	}

	bot, err := tgbotapi.NewBotAPI(token)

	if err != nil {
		return nil, err
	}

	ucfg := tgbotapi.NewUpdate(0)
	ucfg.Timeout = 60
	updates, err := bot.GetUpdatesChan(ucfg)

	if err != nil {
		return nil, err
	}

	return &TelegramBot{
		bot:         bot,
		userUpdates: make(map[int64]chan tgbotapi.Update),
		updates:     updates,
		db:          db,
	}, nil
}

func (t *TelegramBot) Start() {
	for update := range t.updates {
		if update.Message != nil {
			t.commandHandler(update)
		} else if update.CallbackQuery != nil {
			t.callbackHandler(update)
		}
	}
}

func (t *TelegramBot) SendMsg(chatId int64, content string) {
	msg := tgbotapi.NewMessage(chatId, content)
	t.bot.Send(msg)
}

func (t *TelegramBot) SendInlineMsg(chatId int64, title, link, action string) {
	inlineButton := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(action, action),
		),
	)

	msg := tgbotapi.NewMessage(chatId, fmt.Sprintf("%s\n%s", title, link))
	msg.ReplyMarkup = inlineButton
	t.bot.Send(msg)
}
