package bot

import (
	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

var keyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("/saveurl"),
		tgbotapi.NewKeyboardButton("/getall"),
		tgbotapi.NewKeyboardButton("/google"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("/start"),
		tgbotapi.NewKeyboardButton("/help"),
		tgbotapi.NewKeyboardButton("/meme"),
	),
)
