package receipt_service

import (
	"log/slog"
	receiptrepository "qr_code_scanner/internal/storage/repository/receiptRepository"
)

type receiptService struct {
	receiptRepository receiptrepository.ReceiptRepository
	log               *slog.Logger
}

type ReceiptService interface {
	GetReceiptsByOffsetAndLimit(offset string, limit string) ([]receiptrepository.Receipt, error)
}

func GetReceiptService(receiptRepository receiptrepository.ReceiptRepository, log *slog.Logger) ReceiptService {
	return &receiptService{
		receiptRepository: receiptRepository,
		log:               log,
	}
}

func (rs receiptService) GetReceiptsByOffsetAndLimit(offset string, limit string) ([]receiptrepository.Receipt, error) {
	receipts, err := rs.receiptRepository.GetReceiptsByOffsetAndLimit(offset, limit)

	if err != nil {
		return nil, err
	}

	return receipts, nil
}
