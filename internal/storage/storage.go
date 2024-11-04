package storage

import (
	"database/sql"
	"fmt"
	"qr_code_scanner/internal/config"
)

type Storage struct {
	Db *sql.DB
}

const opPackage = "internal.storage.repository.productRpository"

func DbInit(cfg *config.Config) (*Storage, error) {
	const op = "storage.storage.DbInit"
	db, err := sql.Open("postgres", cfg.Database.Url)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{Db: db}, nil
}
