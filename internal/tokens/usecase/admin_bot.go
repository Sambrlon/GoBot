package usecase

import (
	"log"
	"time"

	"bot/internal/tokens/repository"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jmoiron/sqlx"
)

type AdminBot struct {
	Bot             *tgbotapi.BotAPI
	DB              *sqlx.DB
	AdminChatID     int64
	LastClientMsgID map[int64]int
}

func NewAdminBot(token string, db *sqlx.DB, adminChatID int64) (*AdminBot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	return &AdminBot{
		Bot:             bot,
		DB:              db,
		AdminChatID:     adminChatID,
		LastClientMsgID: make(map[int64]int),
	}, nil
}

func (a *AdminBot) Start() {
	updates := a.Bot.GetUpdatesChan(tgbotapi.NewUpdate(0))

	for update := range updates {
		if update.Message != nil {
			message := update.Message.Text
			adminChatID := update.Message.Chat.ID
			clientChatID := update.Message.Chat.ID
			clientMessageID := update.Message.MessageID

			a.sendToClient(clientChatID, message, clientMessageID)
			a.LastClientMsgID[clientChatID] = clientMessageID

			err := a.saveMessage(adminChatID, update.Message.From.UserName, message, true)

			if err != nil {
				log.Printf("Error saving message to DB: %s", err)
			}
		}
	}
}

func (a *AdminBot) ForwardToAdmin(clientChatID int64, clientMessageID int, message string) {
	// Отправляем сообщение администратора как ответ на сообщение клиента
	msg := tgbotapi.NewMessage(clientChatID, message)
	msg.ReplyToMessageID = clientMessageID
	_, err := a.Bot.Send(msg)

	// Сохраняем сообщение в базе данных
	err = a.saveMessage(clientChatID, "Admin", message, true)

	if err != nil {
		log.Printf("Error forwarding message to admin: %s", err)
	}
}

func (a *AdminBot) sendToClient(clientChatID int64, message string, clientMessageID int) {
	msg := tgbotapi.NewMessage(clientChatID, message)
	msg.ReplyToMessageID = clientMessageID
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

	return err
}
