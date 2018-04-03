package postgres

import (
	"github.com/jmoiron/sqlx"

	"github.com/golang-migrate/migrate"
	migratepg "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/lib/pq"
)

type storage struct {
	db *sqlx.DB
}

func NewStorage(dbUrl string) (*storage, error) {
	db, err := sqlx.Connect("postgres", dbUrl)
	if err != nil {
		return nil, err
	}

	storage := &storage{db}
	if err := storage.migrate(); err != nil {
		return nil, err
	}
	return storage, nil
}

func (storage *storage) migrate() error {
	driver, err := migratepg.WithInstance(storage.db.DB, &migratepg.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://postgres/migrations",
		"postgres", driver)
	if err != nil {
		return err
	}
	m.Log = Logger{}
	err = m.Up()
	if err == migrate.ErrNoChange {
		err = nil
	}
	return err
}
