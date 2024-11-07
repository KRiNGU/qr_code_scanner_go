package receiptrepository

import (
	"fmt"
	"qr_code_scanner/internal/storage"
	"time"
)

const opPackage = "internal.storage.repository.receiptRepository"

type Receipt struct {
	Id        int64
	CreatedAt time.Time
}

type receiptRepository struct {
	storage *storage.Storage
}

type ReceiptRepository interface {
	CreateReceipt()
	GetReceiptsByOffsetAndLimit(offset string, limit string) ([]Receipt, error)
}

func GetReceiptRepository(strg *storage.Storage) ReceiptRepository {
	return &receiptRepository{
		storage: strg,
	}
}

func CreateReceipt(strg *storage.Storage) (int64, error) {
	const op = opPackage + ".CreateReceipt"

	const query = `INSERT INTO receipts DEFAULT VALUES RETURNING id`

	var pk int64

	err := strg.Db.QueryRow(query).Scan(&pk)

	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return pk, nil
}

func (receiptRepository receiptRepository) GetReceiptsByOffsetAndLimit(offset string, limit string) ([]Receipt, error) {
	const op = opPackage + ".GetReceiptsByOffsetAndLimit"

	const query = `SELECT * FROM receipts ORDER BY created_at LIMIT $1 OFFSET $2`

	rows, err := receiptRepository.storage.Db.Query(query, limit, offset)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer rows.Close()

	var receipts []Receipt

	for rows.Next() {
		var receipt Receipt
		err := rows.Scan(
			&receipt.Id,
			&receipt.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		receipts = append(receipts, receipt)
	}

	return receipts, nil
}

func (receiptRepository receiptRepository) CreateReceipt() {}
