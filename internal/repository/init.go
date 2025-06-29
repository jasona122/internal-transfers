package repository

import (
	"database/sql"
	"fmt"

	"internal-transfers/internal/config"

	_ "github.com/lib/pq"
)

func InitDB(cfg config.DBConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// verify connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
