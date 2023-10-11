package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// PostgresRepository представляет репозиторий PostgreSQL.
type PostgresRepository struct {
	DB *sqlx.DB
}

// NewPostgresRepository создает новый экземпляр репозитория PostgreSQL.
func NewPostgresRepository(dbConfig DBConfig) (*PostgresRepository, error) {
	dataSourceName := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.Dbname,
	)

	db, err := sqlx.Connect("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}

	return &PostgresRepository{DB: db}, nil
}

// Close закрывает соединение с базой данных.
func (r *PostgresRepository) Close() error {
	return r.DB.Close()
}

// SaveClientChatID сохраняет идентификатор чата клиента в базе данных.
func (r *PostgresRepository) SaveClientChatID(chatID int64) error {
	_, err := r.DB.Exec(`INSERT INTO clients (chat_id) VALUES ($1)`, chatID)
	return err
}
