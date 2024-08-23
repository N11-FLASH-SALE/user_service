package postgres

import (
	"auth/config"
	"auth/storage"
	"database/sql"
	"fmt"
	"log/slog"

	_ "github.com/lib/pq"
)

type postgresStorage struct {
	db  *sql.DB
	log *slog.Logger
}

func NewPostgresStorage(db *sql.DB, log *slog.Logger) storage.IStorage {
	return &postgresStorage{
		db:  db,
		log: log,
	}
}

func ConnectionDb() (*sql.DB, error) {
	conf := config.Load()
	conDb := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		conf.Postgres.DB_HOST, conf.Postgres.DB_PORT, conf.Postgres.DB_USER, conf.Postgres.DB_NAME, conf.Postgres.DB_PASSWORD)
	db, err := sql.Open("postgres", conDb)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func (p *postgresStorage) Close() {
	p.db.Close()
}

func (p *postgresStorage) User() storage.IUserStorage {
	return NewUserRepository(p.db)
}
