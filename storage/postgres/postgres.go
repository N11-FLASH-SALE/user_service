package postgres

import (
	"auth/config"
	"auth/storage"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type postgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage(db *sql.DB) storage.IStorage {
	return &postgresStorage{
		db: db,
	}
}

func ConnectionDb() (*sql.DB, error) {
	conf := config.Load()
	conDb := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		conf.Postgres.PDB_HOST, conf.Postgres.PDB_PORT, conf.Postgres.PDB_USER, conf.Postgres.PDB_NAME, conf.Postgres.PDB_PASSWORD)
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
func (p *postgresStorage) Card() storage.CardStorage{
	return NewCardRepository(p.db)
}

func (p *postgresStorage) Notifications() storage.INotificationStorage {
	return NewNotificationsRepository(p.db)
}
