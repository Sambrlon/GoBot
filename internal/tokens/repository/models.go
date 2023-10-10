package repository

import "time"

type Message struct {
	ID        int64     `db:"id"`
	ChatID    int64     `db:"chat_id"`
	Username  string    `db:"username"`
	Text      string    `db:"text"`
	IsAdmin   bool      `db:"is_admin"`
	Timestamp time.Time `db:"timestamp"`
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Dbname   string
}
