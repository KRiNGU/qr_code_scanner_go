package receiptrepository

import (
	"fmt"
	"qr_code_scanner/internal/storage"
)

const opPackage = "internal.storage.repository.receiptRepository"

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
