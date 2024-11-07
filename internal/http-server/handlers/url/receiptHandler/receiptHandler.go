package receipthandler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	response "qr_code_scanner/internal/lib/api"
	"qr_code_scanner/internal/lib/sl"
	receipt_service "qr_code_scanner/internal/service/receiptService"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

const opPackage = "internal.http-server.handlers.url.receiptHandler"

type ReceiptHandler struct {
	receiptService receipt_service.ReceiptService
	log            *slog.Logger
}

type GetReceiptsByOffsetAndLimitRequest struct {
	Limit  string `json:"limit" validate:"required"`
	Offset string `json:"offset" validate:"required"`
}

func CreateReceiptHandler(router *chi.Mux, receiptService receipt_service.ReceiptService, log *slog.Logger) *chi.Mux {
	handler := &ReceiptHandler{
		receiptService: receiptService,
		log:            log,
	}

	router.Get("/receipts_list", handler.GetReceiptsByOffsetAndLimit)

	return router
}

func (rh *ReceiptHandler) GetReceiptsByOffsetAndLimit(w http.ResponseWriter, r *http.Request) {
	const op = opPackage + ".GetReceiptsByOffsetAndLimit"

	log := rh.log.With(
		slog.String("op", op),
	)

	var req GetReceiptsByOffsetAndLimitRequest

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

	receipts, err := rh.receiptService.GetReceiptsByOffsetAndLimit(req.Offset, req.Limit)

	if err != nil {
		const errMsg = "failed to get receipts by offset and limit"

		log.Error(errMsg, sl.Err(err))

		render.JSON(w, r, response.Error(errMsg))

		return
	}

	err = json.NewEncoder(w).Encode(receipts)

	if err != nil {
		const errMsg = "failed to encode receipts"

		log.Error(errMsg, sl.Err(err))

		render.JSON(w, r, response.Error(errMsg))

		return
	}

	var receiptsIds []int64

	for _, receipt := range receipts {
		receiptsIds = append(receiptsIds, receipt.Id)
	}

	log.Info("receipt found by offset and limit", slog.Any("ids", receiptsIds))

	render.JSON(w, r, response.OK())
}
