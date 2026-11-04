package repository

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/skvdmt/auth-back/internal/model"
)

const (
	KAFKA_MAX_ATTEMPS   = 3
	KAFKA_BATCH_SIZE    = 100
	KAFKA_BATCH_BYTES   = 10_000_000
	KAFKA_BATCH_TIMEOUT = 1_000 * time.Millisecond
	KAFKA_REQUIRED_ACKS = kafka.RequireAll
	KAFKA_ASYNC         = false
	KAFKA_WRITE_TIMEOUT = 3_000 * time.Millisecond
)

// KafkaWriter писатель данных в kafka
type KafkaWriter struct {
	writer  *kafka.Writer
	wg      sync.WaitGroup
	produce chan kafka.Message
	ctx     context.Context
	errors  chan error
}

// KafkaWriterOption Функциональная опция.
type KafkaWriterOption func(k *KafkaWriter)

// KafkaWriterWithTopic Указывает topic в который пишет Writer.
// По-умолчанию он отсутствует и требуется указывать его в сообщениях.
func KafkaWriterWithTopic(topic string) KafkaWriterOption {
	return func(k *KafkaWriter) {
		k.writer.Topic = topic
	}
}

// NewKafkaWriter Контруктор.
func NewKafkaWriter(produce chan kafka.Message, errors chan error, options ...KafkaWriterOption) *KafkaWriter {
	return NewKafkaWriterWithContext(context.Background(), produce, errors, options...)
}

// NewKafkaWriterWithContext Контруктор с контекстом.
func NewKafkaWriterWithContext(ctx context.Context, produce chan kafka.Message, errors chan error, options ...KafkaWriterOption) *KafkaWriter {
	model.Log.Info.Info("kafka writer creating")
	k := &KafkaWriter{
		writer: &kafka.Writer{
			Addr: kafka.TCP(fmt.Sprintf("%s:%d",
				model.Conf.Kafka.Host,
				model.Conf.Kafka.Port)),
			Topic:        "",
			MaxAttempts:  KAFKA_MAX_ATTEMPS,
			BatchSize:    KAFKA_BATCH_SIZE,
			BatchBytes:   KAFKA_BATCH_BYTES,
			BatchTimeout: KAFKA_BATCH_TIMEOUT,
			RequiredAcks: KAFKA_REQUIRED_ACKS,
			Async:        KAFKA_ASYNC,
		},
		wg:      sync.WaitGroup{},
		produce: produce,
		ctx:     ctx,
		errors:  errors,
	}
	for _, o := range options {
		o(k)
	}
	return k
}

// Start Запуск записи сообщений.
func (k *KafkaWriter) Start() {
	k.wg.Add(1)
	go func() {
		defer func() {
			model.Log.Info.Info("kafka writer handler stopped")
			k.wg.Done()
		}()
		model.Log.Info.Info("kafka writer handler starting")
		ms := []kafka.Message{}
		for {
			select {
			case <-k.ctx.Done():
				return
			case m, o := <-k.produce:
				if !o {
					if len(ms) > 0 {
						if err := k.write(ms); err != nil {
							k.errors <- err
						}
					}
					return
				}
				ms = append(ms, m)
			default:
				if len(ms) == 0 {
					continue
				}
				if err := k.write(ms); err != nil {
					k.errors <- err
				}
				ms = ms[:0]
			}
		}
	}()
}

// write Запись бетча сообщений в kafka.
func (k *KafkaWriter) write(messages []kafka.Message) error {
	ctx, cancel := context.WithTimeout(k.ctx, KAFKA_WRITE_TIMEOUT)
	defer cancel()
	if err := k.writer.WriteMessages(ctx, messages...); err != nil {
		return fmt.Errorf("kafka write messages error: %v", err)
	}
	model.Log.Info.Info("kafka write messages complete")
	return nil
}

// Wait Ожидание завершения обработчика записи в kafka.
func (k *KafkaWriter) Wait() error {
	k.wg.Wait()
	if err := k.writer.Close(); err != nil {
		return fmt.Errorf("kafwa writer close error: %v", err)
	}
	model.Log.Info.Info("kafka writer stopped")
	return nil
}
