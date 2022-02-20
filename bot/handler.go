package bot

import (
	"TelegramBot/utils"
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

func (t *TelegramBot) commandHandler(update tgbotapi.Update) {
	chatId := update.Message.Chat.ID

	isCmd := update.Message.IsCommand()

	userUpdate, exist := t.userUpdates[chatId]

	if !exist && !isCmd {
		return
	}

	if exist && !isCmd {
		userUpdate <- update
		return
	}

	command := update.Message.Command()

	userUpdate = make(chan tgbotapi.Update)
	t.userUpdates[chatId] = userUpdate

	log.Printf("switch command: [%v]", command)
	switch command {
	case "start", "help":
		go t.sendHelpInfo(userUpdate)
	case "saveurl":
		go t.saveMessageUrl(userUpdate)
	case "getall":
		go t.getAll(userUpdate)
	case "google":
		go t.googleSearch(userUpdate)
	case "meme":
		go t.sendMeme(userUpdate)
	default:
		t.SendMsg(chatId, "Wrong command")
	}
	userUpdate <- update
}

func (t *TelegramBot) callbackHandler(update tgbotapi.Update) {
	data := update.CallbackQuery.Data

	switch data {
	case "Save":
		go t.saveCallbackUrl(update)
	case "Delete":
		go t.deleteCallbackUrl(update)
	}
}

func (t *TelegramBot) deleteCallbackUrl(update tgbotapi.Update) {
	chatId := update.CallbackQuery.Message.Chat.ID

	titleAndLink := strings.Split(update.CallbackQuery.Message.Text, "\n")

	title := titleAndLink[0]
	link := titleAndLink[1]

	if err := t.db.DeleteLink(chatId, title, link); err != nil {
		msg := fmt.Sprintf("Failed to delete: %s", err.Error())
		log.Print(msg)
		t.SendMsg(chatId, msg)
	}

	msgId := update.CallbackQuery.Message.MessageID
	deleteMsg := tgbotapi.NewDeleteMessage(chatId, msgId)
	t.bot.Send(deleteMsg)

	t.SendMsg(chatId, fmt.Sprintf("Link with title %q was deleted", title))
}

func (t *TelegramBot) saveCallbackUrl(update tgbotapi.Update) {
	chatId := update.CallbackQuery.Message.Chat.ID

	text := update.CallbackQuery.Message.Text

	titleAndLink := strings.Split(text, "\n")
	title := titleAndLink[0]
	url := titleAndLink[1]

	t.saveUrl(chatId, title, url)

	msgId := update.CallbackQuery.Message.MessageID
	inlineButton := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Delete", "Delete"),
		),
	)
	editMsg := tgbotapi.NewEditMessageReplyMarkup(chatId, msgId, inlineButton)
	t.bot.Send(editMsg)
}

func (t *TelegramBot) saveMessageUrl(userUpdate chan tgbotapi.Update) {
	update := <-userUpdate

	chatId := update.Message.Chat.ID

	defer delete(t.userUpdates, chatId)

	t.SendMsg(chatId, "Enter title:")
	update = <-userUpdate
	title := update.Message.Text

	t.SendMsg(chatId, "Enter url:")
	update = <-userUpdate
	url := update.Message.Text

	t.saveUrl(chatId, title, url)
}

func (t *TelegramBot) saveUrl(chatId int64, title, url string) {
	err := t.db.InsertLink(chatId, title, url)
	if err != nil {
		msg := fmt.Sprintf("Error occurred while writing in db: %s", err.Error())
		log.Print(msg)
		t.SendMsg(chatId, msg)
		return
	}
	t.SendMsg(chatId, fmt.Sprintf("Link with title %q was saved", title))
}

func (t *TelegramBot) sendHelpInfo(userUpdate chan tgbotapi.Update) {
	update := <-userUpdate

	chatId := update.Message.Chat.ID

	defer delete(t.userUpdates, chatId)

	text := "/start - start bot\n/help - help info\n/google - google search\n" +
		"/saveurl - save your link\n/getall - get all saved urls"

	msg := tgbotapi.NewMessage(chatId, text)
	msg.ReplyMarkup = keyboard
	t.bot.Send(msg)
}

func (t *TelegramBot) getAll(userUpdate chan tgbotapi.Update) {
	update := <-userUpdate
	chatId := update.Message.Chat.ID

	defer delete(t.userUpdates, chatId)

	links, err := t.db.GetAllLinks(chatId)

	if err != nil {
		msg := fmt.Sprintf("Error occurred while reading from db: %s", err.Error())
		log.Print(msg)
		t.SendMsg(chatId, msg)
		return
	}
	for _, link := range links {
		t.SendInlineMsg(chatId, link.Title, link.Url, "Delete")
	}
}

func (t *TelegramBot) googleSearch(userUpdate chan tgbotapi.Update) {
	update := <-userUpdate

	chatId := update.Message.Chat.ID

	defer delete(t.userUpdates, chatId)

	t.SendMsg(chatId, "Enter search query:")

	update = <-userUpdate

	sr, err := utils.GoogleCustomSearchRequest(update.Message.Text)
	if err != nil {
		log.Printf("Google search failed")
		t.SendMsg(chatId, "Google search failed")
		return
	}

	for _, item := range sr.Items {
		t.SendInlineMsg(chatId, item.Title, item.Link, "Save")
	}
}

func (t *TelegramBot) sendMeme(userUpdate chan tgbotapi.Update) {
	update := <-userUpdate

	chatId := update.Message.Chat.ID

	defer delete(t.userUpdates, chatId)

	meme, err := utils.GetMeme()

	if err != nil {
		log.Println("No memes")
		t.SendMsg(chatId, "No memes")
		return
	}
	fmt.Printf("%s\n%s\n", meme.Title, meme.Url)

	if err != nil {
		log.Println("No memes")
		t.SendMsg(chatId, "No memes")
		return
	}

	if strings.HasSuffix(meme.Url, ".gif") {
		msg := tgbotapi.NewAnimationUpload(chatId, nil)
		msg.FileID = meme.Url
		msg.UseExisting = true
		t.bot.Send(msg)
	} else {
		msg := tgbotapi.NewPhotoUpload(chatId, nil)
		msg.FileID = meme.Url
		msg.UseExisting = true
		t.bot.Send(msg)
	}
}
