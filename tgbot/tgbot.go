package tgbot

import (
	"encoding/json"
	"fmt"

	"github.com/RaghavSood/collectibles/clogger"
	"github.com/RaghavSood/collectibles/storage"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var log = clogger.NewLogger("tgbot")

type TgBot struct {
	bot             *tgbotapi.BotAPI
	channelUsername string
	db              storage.Storage
}

func NewBot(token string, channelUsername string, db storage.Storage) (*TgBot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	tgBot := &TgBot{
		bot:             bot,
		channelUsername: channelUsername,
		db:              db,
	}

	go tgBot.Run()

	return tgBot, nil
}

func (t *TgBot) Run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30

	updates := t.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Info().
			Int64("user", update.Message.From.ID).
			Str("username", update.Message.From.UserName).
			Bool("is_bot", update.Message.From.IsBot).
			Str("text", update.Message.Text).
			Msg("Received message")

		err := t.SaveMessage(update.Message.Chat.ID, update)
		if err != nil {
			log.Println("Error saving message:", err)
		}

		if !update.Message.IsCommand() {
			t.SendChatMessage(update.Message.Chat.ID, "I only understand commands. Please use /help to see the list of commands.")
			continue
		}

		var msg string

		switch update.Message.Command() {
		case "help":
			msg = "I understand the following commands:\n" +
				"/help - Show this help message\n" +
				"/subscribe - Subscribe to updates for a creator, series, or item by sending it's link\n" +
				"/unsubscribe - Unsubscribe from updates for a creator, series, or item by sending it's link\n" +
				"/admin - Talk to the boss\n"
		case "admin":
			msg = "Chat with @RaghavSood if you're facing issues with the bot, or would like to submit collectibles or learn about them."
		case "subscribe":
			payload := update.Message.CommandArguments()
			if len(payload) == 0 {
				msg = "Please provide a link to the creator, series, or item you'd like to subscribe to.\n\nExample: /subscribe https://collectible.money/creator/casascius"
				break
			}

			// Parse the URL
			urlType, slug, err := ParseURL(payload)
			if err != nil {
				msg = fmt.Sprintf("Failed to parse link: %s", err)
				break
			}

			err = t.db.UpsertTelegramSubscription(update.Message.Chat.ID, urlType, slug)
			if err != nil {
				msg = fmt.Sprintf("Failed to subscribe: %s", err)
			}

			msg = "Subscribed successfully! You'll receive updates for this " + urlType + " from now on."
		default:
			msg = "I don't know that command. Try /help"
		}

		err = t.SendChatMessage(update.Message.Chat.ID, msg)
		if err != nil {
			log.Error().Err(err).Msg("Error sending message")
			continue
		}
	}
}

func (t *TgBot) SendChatMessage(chatID int64, message string) error {
	msg := tgbotapi.NewMessage(chatID, message)
	err := t.SaveMessage(chatID, msg)
	if err != nil {
		fmt.Println("Error saving message:", err)
	}

	_, err = t.bot.Send(msg)
	return err
}

func (t *TgBot) SendMessage(message string) error {
	msg := tgbotapi.NewMessageToChannel(t.channelUsername, message)
	msg.ParseMode = "MarkdownV2"
	msg.DisableWebPagePreview = true
	fmt.Printf("%+v\n", msg)
	_, err := t.bot.Send(msg)
	return err
}

func (t *TgBot) SaveMessage(chatID int64, msg any) error {
	msgJson, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message to JSON: %w", err)
	}

	return t.db.InsertMessage(chatID, string(msgJson))
}

func EscapeText(text string) string {
	return tgbotapi.EscapeText("MarkdownV2", text)
}
