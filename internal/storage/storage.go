package storage

import (
	"database/sql"
	"fmt"
	"qr_code_scanner/internal/config"
)

type Storage struct {
	Db *sql.DB
}

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

type CreateProductDto struct {
	Name           string
	TranslatedName string
}

func (strg *Storage) CreateProduct(product *CreateProductDto) (int64, error) {
	const op = "internal.storage.repository.productRpository.AddNewProduct"

	const query = `INSERT INTO products (product_name, translated_name)
	VALUES ($1, $2) RETURNING id`

	var pk int64

	err := strg.Db.QueryRow(query, product.Name, product.TranslatedName).Scan(&pk)

	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return pk, nil
}
