package transactionhandler

import (
	"log/slog"
	"net/http"
	response "qr_code_scanner/internal/lib/api"
	"qr_code_scanner/internal/lib/sl"
	"qr_code_scanner/internal/storage"
	"qr_code_scanner/internal/storage/repository/productRepository"
	receiptrepository "qr_code_scanner/internal/storage/repository/receiptRepository"
	transactionsrepository "qr_code_scanner/internal/storage/repository/transactionsRepository"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

const opPackage = "internal.http-server.handlers.url.transactionHandler"

type SingleTransactionRequest struct {
	Name           string  `json:"name" validate:"required"`
	TranslatedName string  `json:"translatedName"`
	Price          float32 `json:"cost_unit" validate:"required"`
	Amount         float32 `json:"amount" validate:"required"`
}

type CreateTransactionRequest struct {
	Transactions []SingleTransactionRequest `json:"transactions"`
}

type CreateTransactionResponse struct {
	response.Response
	Error string `json:"error,omitempty"`
}

func createTransaction(strg *storage.Storage, transaction SingleTransactionRequest, receiptId int64) (int64, error) {
	if transaction.TranslatedName != "" {
		productRepository.CreateProduct(
			&productRepository.CreateProductDto{Name: transaction.Name, TranslatedName: transaction.TranslatedName},
			strg,
		)
	} else {
		//todo *todo2*
		productRepository.CreateProduct(
			&productRepository.CreateProductDto{Name: transaction.Name, TranslatedName: transaction.Name},
			strg,
		)
	}
	transactionId, err := transactionsrepository.CreateTransaction(&transactionsrepository.CreateTransactionDto{
		Price:       transaction.Price,
		Amount:      transaction.Amount,
		ProductName: transaction.Name,
		ReceiptId:   receiptId,
	}, strg)

	if err != nil {
		return 0, err
	}

	return transactionId, nil
}

func CreateTransactionHandler(log *slog.Logger, strg *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = opPackage + ".CreateTransactionHandler"

		log := log.With(
			slog.String("op", op),
		)

		var req CreateTransactionRequest

		err := render.DecodeJSON(r.Body, &req)

		if err != nil {
			const errMsg = "failed to decode request body"

			log.Error(errMsg, sl.Err(err))

			render.JSON(w, r, response.Error(errMsg))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		//todo *todo1*
		if err := validator.New().Struct(req); err != nil {
			const errMsg = "invalid request"

			log.Error(errMsg, sl.Err(err))

			render.JSON(w, r, response.Error(errMsg))

			return
		}

		receiptId, err := receiptrepository.CreateReceipt(strg)

		if err != nil {
			const errMsg = "failed to add new receipt -> failed to add new transaction"

			log.Error(errMsg)

			render.JSON(w, r, response.Error(errMsg))

			return
		}

		for _, transaction := range req.Transactions {
			_, err := createTransaction(strg, transaction, receiptId)

			if err != nil {
				const errMsg = "failed to add new transaction"

				log.Error(errMsg, slog.String("product", transaction.Name), slog.Int64("receipt", receiptId))

				render.JSON(w, r, response.Error(errMsg))

				return
			}
		}

		log.Info("receipt added", slog.Int64("id", receiptId))

		render.JSON(w, r, response.OK())
	}
}

func CreateTransactionsBatch(log *slog.Logger, strg *storage.Storage, transactions []SingleTransactionRequest) {
	receiptId, err := receiptrepository.CreateReceipt(strg)

	if err != nil {
		const errMsg = "failed to add new receipt -> failed to add new transaction"

		log.Error(errMsg)

		return
	}

	for _, transaction := range transactions {
		_, err := createTransaction(strg, transaction, receiptId)

		if err != nil {
			const errMsg = "failed to add new transaction"

			log.Error(errMsg, slog.String("product", transaction.Name), slog.Int64("receipt", receiptId))

			return
		}
	}
}
