package tgbot

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TgBot struct {
	bot             *tgbotapi.BotAPI
	channelUsername string
}

func NewBot(token string, channelUsername string) (*TgBot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	return &TgBot{
		bot:             bot,
		channelUsername: channelUsername,
	}, nil
}

func (t *TgBot) SendMessage(message string) error {
	msg := tgbotapi.NewMessageToChannel(t.channelUsername, message)
	msg.ParseMode = "MarkdownV2"
	msg.DisableWebPagePreview = true
	fmt.Printf("%+v\n", msg)
	_, err := t.bot.Send(msg)
	return err
}

func EscapeText(text string) string {
	return tgbotapi.EscapeText("MarkdownV2", text)
}
