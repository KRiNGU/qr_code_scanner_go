package productHandler

import (
	"errors"
	"log/slog"
	"net/http"
	response "qr_code_scanner/internal/lib/api"
	"qr_code_scanner/internal/lib/sl"
	"qr_code_scanner/internal/storage"
	productRepository "qr_code_scanner/internal/storage/repository/productRepository"

	"github.com/go-chi/render"
	"github.com/go-playground/validator"
)

const opPackage = "internal.handlers.productHandler"

type CreateProductRequest struct {
	Name           string `json:"name" validate:"required"`
	TranslatedName string `json:"translatedName" validate:"required"`
}

type CreateProductResponse struct {
	response.Response
	Error string `json:"error,omitempty"`
}

func CreateProductHandler(log *slog.Logger, strg *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = opPackage + ".CreateProductHandler"

		log := log.With(
			slog.String("op", op),
		)

		var req CreateProductRequest

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

		createProductDto := productRepository.CreateProductDto{Name: req.Name, TranslatedName: req.TranslatedName}

		id, err := productRepository.CreateProduct(&createProductDto, strg)

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

		log.Info("product added", slog.Int64("id", id))

		render.JSON(w, r, response.OK())
	}
}
