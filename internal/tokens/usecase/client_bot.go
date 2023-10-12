package usecase

import (
	"log"
	"time"

	"bot/internal/tokens/repository"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jmoiron/sqlx"
)

type ClientBot struct {
	Bot       *tgbotapi.BotAPI
	ChatID    int64
	AdminChat int64
	AdminBot  *AdminBot
	DB        *sqlx.DB
}

func NewClientBot(token string, chatID int64, adminChat int64, adminBot *AdminBot, db *sqlx.DB) (*ClientBot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	return &ClientBot{
		Bot:       bot,
		ChatID:    chatID,
		AdminChat: adminChat,
		AdminBot:  adminBot,
		DB:        db,
	}, nil
}

func (c *ClientBot) Start() {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates := c.Bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message != nil {
			message := update.Message.Text
			clientChatID := update.Message.Chat.ID
			clientMessageID := update.Message.MessageID

			err := c.saveClientChatID(clientChatID)

			if err != nil {
				log.Printf("Error saving client chat ID: %s", err)
			}

			c.AdminBot.sendToClient(clientChatID, message, clientMessageID)

			err = c.saveMessage(clientChatID, update.Message.From.UserName, message, false)
			if err != nil {
				log.Printf("Error saving message to DB: %s", err)
			}
		}
	}
}

func (c *ClientBot) sendToAdmin(message string, clientChatID int64, clientMessageID int) {
	msg := tgbotapi.NewMessage(c.AdminChat, message)
	msg.ReplyToMessageID = clientMessageID
	msg.BaseChat.ChatID = clientChatID
	_, err := c.Bot.Send(msg)
	if err != nil {
		log.Printf("Error sending message to admin: %s", err)
	}
}

func (c *ClientBot) saveClientChatID(chatID int64) error {
	_, err := c.DB.Exec(`
        INSERT INTO messages (chat_id)
        VALUES ($1)
        ON CONFLICT (chat_id) DO NOTHING
    `, chatID)

	return err
}

func (c *ClientBot) saveMessage(chatID int64, username, text string, isAdmin bool) error {
	message := repository.Message{
		ChatID:    chatID,
		Username:  username,
		Text:      text,
		IsAdmin:   isAdmin,
		Timestamp: time.Now(),
	}

	_, err := c.DB.NamedExec(`
        INSERT INTO messages (chat_id, username, text, is_admin, timestamp)
        VALUES (:chat_id, :username, :text, :is_admin, :timestamp)
    `, message)

	return err
}
