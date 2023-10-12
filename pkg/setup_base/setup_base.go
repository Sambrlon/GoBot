package setup_base

import "bot/internal/tokens/repository"

func SetupDatabase(dbConfig repository.DBConfig) (*repository.PostgresRepository, error) {
	db, err := repository.NewPostgresRepository(dbConfig)
	if err != nil {
		return nil, err
	}
	return db, nil
}
