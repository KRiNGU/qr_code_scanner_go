package transactionsrepository

import (
	"fmt"
	"qr_code_scanner/internal/storage"
	"time"
)

const opPackage = "internal.storage.repository.transactionsRepository"

type CreateTransactionDto struct {
	Price       float32
	Amount      float32
	ReceiptId   int64
	ProductName string
}

type TransactionDto struct {
	CreateTransactionDto
	Id        int64
	CreatedAt time.Time
}

func CreateTransaction(transaction *CreateTransactionDto, strg *storage.Storage) (int64, error) {
	const op = opPackage + ".CreateTransaction"

	const query = `INSERT INTO transactions (price, amount, receipt_id, product_name_fk)
	VALUES ($1, $2, $3, $4) RETURNING id`

	var pk int64

	err := strg.Db.QueryRow(
		query,
		transaction.Price,
		transaction.Amount,
		transaction.ReceiptId,
		transaction.ProductName).Scan(&pk)

	if err != nil {
		fmt.Println(err)
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return pk, nil
}
