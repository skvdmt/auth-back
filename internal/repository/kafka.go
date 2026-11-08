package repository

import (
	"context"

	"github.com/segmentio/kafka-go"
	"github.com/skvdmt/auth-back/internal/model"
	"github.com/skvdmt/secret-contracts/pkg/go/gen/auth_error_v1"
	"google.golang.org/protobuf/proto"
)

const (
	AUTH_ERROR_TOPIC = "auth-error"
)

// Kafka Репозиторий кафка.
type Kafka struct {
	AuthErrorMessages chan kafka.Message
	AuthErrorWriter   *KafkaWriter
}

// NewKafka Конструктор.
func NewKafka() (*Kafka, error) {
	m := make(chan kafka.Message)
	aew := NewKafkaWriter(
		m,
		model.ExternalErrors,
		KafkaWriterWithTopic(AUTH_ERROR_TOPIC),
	)
	aew.Start()
	return &Kafka{
		AuthErrorMessages: m,
		AuthErrorWriter:   aew,
	}, nil
}

// AuthErrorSending Отправка ошибки аутентификации.
func (k *Kafka) AuthErrorSending(ctx context.Context, event *auth_error_v1.Event) error {
	m, err := proto.Marshal(event)
	if err != nil {
		return err
	}
	k.AuthErrorMessages <- kafka.Message{
		Key:   event.Id,
		Value: m,
	}
	return nil
}

// Stop Отстановка.
func (k *Kafka) Stop() error {
	// Закрытие канала сообщений с ошибками авторизации.
	close(k.AuthErrorMessages)
	// Ожидание сброса оставшихся сообщений с ошибками авторизации.
	k.AuthErrorWriter.Wait()
	return nil
}
