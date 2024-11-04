package productHandler

import (
	"errors"
	"log/slog"
	"net/http"
	response "qr_code_scanner/internal/lib/api"
	"qr_code_scanner/internal/lib/sl"
	"qr_code_scanner/internal/storage"
	productRepository "qr_code_scanner/internal/storage/repository"

	"github.com/go-chi/render"
	"github.com/go-playground/validator"
)

type Request struct {
	Name           string `json:"name" validate:"required"`
	TranslatedName string `json:"translatedName" validate:"required"`
}

type Response struct {
	response.Response
	Error string `json:"error,omitempty"`
}

const opPackage = "internal.handlers.productHandler"

type CreateProductI interface {
	CreateProduct(createProductDto *storage.CreateProductDto) (int64, error)
}

func CreateProductHandler(log *slog.Logger, createProductI CreateProductI) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = opPackage + "CreateProductHandler"

		log := log.With(
			slog.String("op", op),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)

		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, response.Error("failed to decode JSON"))

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

		createProductDto := storage.CreateProductDto{Name: req.Name, TranslatedName: req.TranslatedName}

		id, err := createProductI.CreateProduct(&createProductDto)

		if errors.Is(err, productRepository.ErrProductExists) {
			const errMsg = "product already exists"

			log.Info(errMsg, slog.String("product", req.Name))

			render.JSON(w, r, response.Error(errMsg))

			return
		}

		if err != nil {
			const errMsg = "failed to add new product"

			log.Error(errMsg, slog.String("product", req.Name))

			render.JSON(w, r, response.Error(errMsg))

			return
		}

		log.Info("product created", slog.Int64("id", id))

		render.JSON(w, r, response.OK())
	}
}
