package postgres

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/dotchev/sm/storage"
	"github.com/jmoiron/sqlx"

	"github.com/golang-migrate/migrate"
	migratepg "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/lib/pq"
)

type pgStorage struct {
	db *sqlx.DB
}

func NewStorage(dbUrl string) (*pgStorage, error) {
	db, err := sqlx.Connect("postgres", dbUrl)
	if err != nil {
		return nil, err
	}

	store := &pgStorage{db}
	if err := store.migrate(); err != nil {
		return nil, err
	}
	return store, nil
}

func (store *pgStorage) migrate() error {
	driver, err := migratepg.WithInstance(store.db.DB, &migratepg.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://storage/postgres/migrations",
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

func (store *pgStorage) AddPlatform(platform *storage.Platform) error {
	_, err := store.db.NamedExec(`INSERT INTO platforms (id, name, type, description)
		VALUES (:id, :name, :type, :description)`,
		platform)
	return err
}

func (store *pgStorage) UpdatePlatform(platform *storage.Platform) error {
	set := make([]string, 0, 5)
	if platform.Name != "" {
		set = append(set, "name = :name")
	}
	if platform.Type != "" {
		set = append(set, "type = :type")
	}
	if platform.Description != "" {
		set = append(set, "description = :description")
	}
	if len(set) == 0 {
		return nil // nothing to update
	}
	update := fmt.Sprintf("UPDATE platforms SET %s WHERE id = :id",
		strings.Join(set, ", "))
	println(update)
	result, err := store.db.NamedExec(update, platform)
	if err != nil {
		return err
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return storage.ErrNotFound
	}
	return nil
}

func (store *pgStorage) GetPlatforms() (platforms []storage.Platform, err error) {
	err = store.db.Select(&platforms, "SELECT * FROM platforms ORDER BY name ASC")
	return
}

func (store *pgStorage) GetPlatform(id string) (*storage.Platform, error) {
	var platform storage.Platform
	err := store.db.Get(&platform, "SELECT * FROM platforms WHERE id = $1", id)
	if err == sql.ErrNoRows {
		err = storage.ErrNotFound
	}
	return &platform, err
}

func (store *pgStorage) DeletePlatform(id string) (deleted bool, err error) {
	result, err := store.db.Exec("DELETE FROM platforms WHERE id = $1", id)
	if err == nil {
		n, _ := result.RowsAffected()
		if n > 0 {
			deleted = true
		}
	}
	return
}
