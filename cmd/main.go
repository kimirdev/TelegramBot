package main

import (
	"TelegramBot/bot"
	"TelegramBot/utils"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	db, err := utils.NewMongoDB()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Disconnect()

	tgbot, err := bot.NewTelegramBot(db)
	if err != nil {
		log.Fatal(err)
	}

	tgbot.Start()
}
