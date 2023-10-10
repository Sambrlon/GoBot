package main

import (
	"bot/config"
	"bot/internal/tokens/repository"
	"bot/internal/tokens/usecase"
	"log"
)

func main() {
	cfg := config.LoadConfig(".env")
	dbConfig := repository.DBConfig{
		Host:     cfg.DBConfig.Host,
		Port:     cfg.DBConfig.Port,
		User:     cfg.DBConfig.User,
		Password: cfg.DBConfig.Password,
		Dbname:   cfg.DBConfig.Dbname,
	}

	db, err := config.SetupDatabase(dbConfig)
	if err != nil {
		log.Fatalf("Error setting up the database: %s", err)
	}

	clientBot, err := usecase.NewClientBot(cfg.ClientBotToken, cfg.ClientChatID, 0, db.DB)
	if err != nil {
		log.Fatalf("Error initializing client bot: %s", err)
	}
	go clientBot.Start()

	adminBot, err := usecase.NewAdminBot(cfg.AdminBotToken, db.DB)
	if err != nil {
		log.Fatalf("Error initializing admin bot: %s", err)
	}
	go adminBot.Start(cfg.ClientChatID)

	select {}
}
