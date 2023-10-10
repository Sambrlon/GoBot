package repository

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PostgresRepository struct {
	DB *sqlx.DB
}

func NewPostgresRepository(dbConfig DBConfig) (*PostgresRepository, error) {
	connStr := "user=" + dbConfig.User + " dbname=" + dbConfig.Dbname + " host=" + dbConfig.Host + " port=" + dbConfig.Port + " password=" + dbConfig.Password + " sslmode=disable"
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, err
	}

	return &PostgresRepository{DB: db}, nil
}
