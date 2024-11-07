package urlhandler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"qr_code_scanner/internal/kafka"
	response "qr_code_scanner/internal/lib/api"
	"qr_code_scanner/internal/lib/sl"
	"qr_code_scanner/internal/storage"
	receiptrepository "qr_code_scanner/internal/storage/repository/receiptRepository"

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

		id, err := receiptrepository.CreateReceipt(strg)

		if err != nil {
			const errMsg = "failed to create receipt"

			log.Error(errMsg, sl.Err(err))

			render.JSON(w, r, response.Error(errMsg))

			return
		}

		SendUrl(log, req.Url, id)

		render.JSON(w, r, response.OK())
	}
}

func SendUrl(log *slog.Logger, url string, id int64) {
	const op = opPackage + ".SendUrl"

	log = log.With(
		slog.String("op", op),
	)

	transactionInBytes, err := json.Marshal(url)

	if err != nil {
		const errMsg = "failed to parse url to bites"

		log.Error(errMsg, sl.Err(err))
	}

	err = kafka.PushToKafkaProducer(log, kafka.TopicsEnum[kafka.Links], transactionInBytes, id)

	if err != nil {
		const errMsg = "failed to send url"

		log.Error(errMsg, sl.Err(err))
	}
}
