package urlhandler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	response "qr_code_scanner/internal/lib/api"
	"qr_code_scanner/internal/lib/sl"
	"qr_code_scanner/internal/storage"
	transactionsrepository "qr_code_scanner/internal/storage/repository/transactionsRepository"
	"time"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

const opPackage = "internal.hppt-server.handlers.url.urlHandler"

type ScanUrl struct {
	Url string `json:"url" validate:"required,url"`
}

func ScanUrlHandler(log *slog.Logger, strg *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = opPackage + ".ScanUrlHandler"

		log := log.With(
			slog.String("op", op),
		)

		var req ScanUrl

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

		td := transactionsrepository.TransactionDto{
			CreateTransactionDto: transactionsrepository.CreateTransactionDto{
				Price:       1.0,
				Amount:      2.0,
				ReceiptId:   1,
				ProductName: "Product",
			},
			Id:        1,
			CreatedAt: time.Now(),
		}

		//todo *todo2*

		err = json.NewEncoder(w).Encode(td)

		if err != nil {
			const errMsg = "failed to encode response json"

			log.Error(errMsg, sl.Err(err))

			render.JSON(w, r, response.Error(errMsg))

			return
		}

		render.JSON(w, r, response.OK())
	}
}
