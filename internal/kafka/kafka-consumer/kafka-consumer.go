package kafkaconsumer

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	transactionhandler "qr_code_scanner/internal/http-server/handlers/url/transactionHandler"
	"qr_code_scanner/internal/storage"
	"strconv"
	"syscall"

	"github.com/IBM/sarama"
)

func connectConsumer(brokers []string) (sarama.Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer(brokers, config)

	return consumer, err
}

func CreateConsumer(log *slog.Logger, brokers []string, topic string, strg *storage.Storage) error {
	const op = "internal.kafka.kafka-consumer.CreateConsumer"
	worker, err := connectConsumer(brokers)
	msgCnt := 0

	log = log.With(
		slog.String("op", op),
	)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	consumer, err := worker.ConsumePartition(topic, 0, sarama.OffsetNewest)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	doneCh := make(chan struct{})

	go func() {
		for {
			log.Info("Waiting...")
			select {
			case err := <-consumer.Errors():
				log.Info(err.Error())
			case msg := <-consumer.Messages():
				var transactions []transactionhandler.SingleTransactionRequest
				msgCnt++
				log.Info(fmt.Sprintf("received message for topic(%d),count(%s),message(%s)\n", msgCnt, string(msg.Topic), string(msg.Value)))
				err := json.Unmarshal(msg.Value, &transactions)
				if err != nil {
					log.Info("failed to deserialize")
				}
				key, _ := strconv.ParseInt(string(msg.Key), 10, 64)
				transactionhandler.CreateTransactionsBatch(log, strg, transactions, key)
			case <-sigchan:
				log.Info("kafka interrupt is detected")
				doneCh <- struct{}{}
			}
		}
	}()

	<-doneCh
	log.Info("processed " + string(msgCnt) + " messages")

	if err := worker.Close(); err != nil {
		return err
	}

	return nil
}
