package kafka

import (
	"fmt"
	"log/slog"
	"qr_code_scanner/internal/lib/sl"
	"strconv"

	"github.com/IBM/sarama"
)

const opPackage = "internal.kafka.kafka-producer"

type TopicsState int

const (
	Transactions TopicsState = iota
	Links
)

var TopicsEnum = map[TopicsState]string{
	Transactions: "transactions",
	Links:        "links",
}

func ConnectProducer(brokers []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	return sarama.NewSyncProducer(brokers, config)
}

func PushToKafkaProducer(log *slog.Logger, topic string, message []byte, key int64) error {
	const op = opPackage + ".PushToKafkaProducer"
	brokers := []string{"localhost:9092"}

	// Create connection

	log = log.With(
		slog.String("op", op),
	)

	producer, err := ConnectProducer(brokers)
	if err != nil {
		const errMsg = "failed to connect producer"
		log.Error(errMsg, sl.Err(err))
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
		Key:   sarama.StringEncoder(strconv.Itoa(int(key))),
	}

	partition, offset, err := producer.SendMessage(msg)

	if err != nil {
		const errMsg = "failed to send message"
		log.Error(errMsg, sl.Err(err))
		return err
	}

	log.Info(fmt.Sprintf("order is stored in topic(%s)/partition(%d)/offset(%d)\n", topic, partition, offset))

	return nil
}
