package main

import (
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"qr_code_scanner/internal/config"
	productHandler "qr_code_scanner/internal/http-server/handlers/url/productHandler"
	receipthandler "qr_code_scanner/internal/http-server/handlers/url/receiptHandler"
	transactionHandler "qr_code_scanner/internal/http-server/handlers/url/transactionHandler"
	urlhandler "qr_code_scanner/internal/http-server/handlers/url/urlHandler"
	kafkaconsumer "qr_code_scanner/internal/kafka/kafka-consumer"
	"qr_code_scanner/internal/lib/sl"
	receipt_service "qr_code_scanner/internal/service/receiptService"
	"qr_code_scanner/internal/storage"
	receiptrepository "qr_code_scanner/internal/storage/repository/receiptRepository"
	"syscall"

	"github.com/go-chi/chi"
	_ "github.com/lib/pq"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)

	log.Info("Starting logger", slog.String("env", cfg.Env))
	log.Debug("Debug info")

	storage, err := storage.DbInit(cfg)

	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	go func() {
		if err := kafkaconsumer.CreateConsumer(log, []string{"localhost:9092"}, "transactions", storage); err != nil {
			log.Error("failed to start kafka server")
		}
	}()

	receiptService := initService(storage, log)

	router := chi.NewRouter()
	router.Post("/products", productHandler.CreateProductHandler(log, storage))
	router.Post("/transaction", transactionHandler.CreateTransactionHandler(log, storage))
	router.Post("/scan_url", urlhandler.ScanUrlHandler(log, storage))

	receipthandler.CreateReceiptHandler(router, receiptService, log)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error("failed to start server")
		}
	}()

	log.Info("server started")

	<-done
	log.Info("stopping server")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

func initService(strg *storage.Storage, log *slog.Logger) receipt_service.ReceiptService {
	receiptRepository := receiptrepository.GetReceiptRepository(strg)

	receiptService := receipt_service.GetReceiptService(receiptRepository, log)

	return receiptService
}
