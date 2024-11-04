package productRepository

import (
	"errors"
	"fmt"
	"qr_code_scanner/internal/storage"
)

var (
	ErrProductNotFound = errors.New("product not found")
	ErrProductExists   = errors.New("product already exists")
)

type CreateProductRequest struct {
	Name           float32 `json:"name" validate:"required"`
	TranslatedName float32 `json:"translatedName" validate:"required"`
}

const opPackage = "internal.storage.repository.productRepository"

type CreateProductDto struct {
	Name           string
	TranslatedName string
}

func CreateProduct(product *CreateProductDto, strg *storage.Storage) (int64, error) {
	const op = opPackage + ".AddNewProduct"

	const query = `INSERT INTO products (product_name, translated_name)
	VALUES ($1, $2) RETURNING product_name`

	var pk int64

	err := strg.Db.QueryRow(query, product.Name, product.TranslatedName).Scan(&pk)

	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return pk, nil
}
