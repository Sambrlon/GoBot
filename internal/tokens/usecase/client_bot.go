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
	DB        *sqlx.DB
}

func NewClientBot(token string, chatID int64, adminChat int64, db *sqlx.DB) (*ClientBot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	return &ClientBot{
		Bot:       bot,
		ChatID:    chatID,
		AdminChat: adminChat,
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

			c.sendToAdmin(message)

			err := c.saveMessage(update.Message.Chat.ID, update.Message.From.UserName, message, false)
			if err != nil {
				log.Printf("Error saving message to DB: %s", err)
			}

			c.sendToClient("Ваш вопрос был отправлен администратору для обработки. Ожидайте ответа.")
		}
	}
}

func (c *ClientBot) sendToAdmin(message string) {
	msg := tgbotapi.NewMessage(c.AdminChat, message)
	_, err := c.Bot.Send(msg)
	if err != nil {
		log.Printf("Error sending message to admin: %s", err)
	}
}

func (c *ClientBot) sendToClient(message string) {
	msg := tgbotapi.NewMessage(c.ChatID, message)
	_, err := c.Bot.Send(msg)
	if err != nil {
		log.Printf("Error sending message to client: %s", err)
	}
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
