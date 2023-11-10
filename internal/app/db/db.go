package db

import (
	"fmt"

	mg "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DB       string
}

func New(cfg Config) (*sqlx.DB, error) {
	connURL := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DB)
	db, err := sqlx.Connect("postgres", connURL)
	if err != nil {
		return nil, fmt.Errorf("couldn't connect to database: %f", err)
	}
	if err = migrate(cfg); err != nil {
		return nil, fmt.Errorf("couldn't make migrations: %f", err)
	}
	return db, nil
}

func migrate(cfg Config) error {
	connURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DB)
	m, err := mg.New("file://migrations", connURL)
	if err != nil {
		return err
	}
	if err = m.Up(); err != nil && err != mg.ErrNoChange {
		return err
	}
	return nil
}
