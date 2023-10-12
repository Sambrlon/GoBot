package main

import (
	"bot/config"
	"bot/internal/tokens/repository"
	"bot/internal/tokens/usecase"
	"bot/pkg/setup_base"
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
	} // Зачем в маине находится структура подклчения к бд

	db, err := setup_base.SetupDatabase(dbConfig)
	if err != nil {
		log.Fatalf("Error setting up the database: %s", err)
	}

	adminBot, err := usecase.NewAdminBot(cfg.AdminBotToken, db.DB, cfg.AdminChatID)
	if err != nil {
		log.Fatalf("Error initializing admin bot: %s", err)
	}
	go adminBot.Start()

	clientBot, err := usecase.NewClientBot(cfg.ClientBotToken, cfg.ClientChatID, cfg.AdminChatID, adminBot, db.DB)

	if err != nil {
		log.Fatalf("Error initializing client bot: %s", err)
	}
	go clientBot.Start()

	select {}
}
