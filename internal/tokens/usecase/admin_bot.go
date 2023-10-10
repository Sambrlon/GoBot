package usecase

import (
	"log"
	"time"

	"bot/internal/tokens/repository"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jmoiron/sqlx"
)

type AdminBot struct {
	Bot *tgbotapi.BotAPI
	DB  *sqlx.DB
}

func NewAdminBot(token string, db *sqlx.DB) (*AdminBot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	return &AdminBot{
		Bot: bot,
		DB:  db,
	}, nil
}

func (a *AdminBot) Start(clientChatID int64) {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates := a.Bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message != nil {
			message := update.Message.Text

			a.sendToClient(clientChatID, message)

			err := a.saveMessage(update.Message.Chat.ID, update.Message.From.UserName, message, true)
			if err != nil {
				log.Printf("Error saving message to DB: %s", err)
			}
		}
	}
}

func (a *AdminBot) sendToClient(clientChatID int64, message string) {
	msg := tgbotapi.NewMessage(clientChatID, message)
	_, err := a.Bot.Send(msg)
	if err != nil {
		log.Printf("Error sending message to client: %s", err)
	}
}

func (a *AdminBot) saveMessage(chatID int64, username, text string, isAdmin bool) error {
	message := repository.Message{
		ChatID:    chatID,
		Username:  username,
		Text:      text,
		IsAdmin:   isAdmin,
		Timestamp: time.Now(),
	}

	_, err := a.DB.NamedExec(`
        INSERT INTO messages (chat_id, username, text, is_admin, timestamp)
        VALUES (:chat_id, :username, :text, :is_admin, :timestamp)
    `, message)

	if err != nil {
		log.Printf("Error saving message to DB: %s", err)
	}
	return err
}
