package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/golang-migrate/migrate/v4"
	driver "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type DB struct {
	*sql.DB
}

func Open(databaseURL string) (*DB, error) {
	var database DB
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, err
	}
	database.DB = db
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	driver, err := driver.WithInstance(db, &driver.Config{})
	if err != nil {
		return nil, err
	}
	_, cwd, _, _ := runtime.Caller(0)
	path := filepath.Join(filepath.Dir(cwd), "..", "..", "migrations")
	migrator, err := migrate.NewWithDatabaseInstance("file://"+path, "postgres", driver)
	if err != nil {
		return nil, err
	}
	if err := migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, err
	}
	return &database, nil

}
