package usecase

import (
	"context"

	"github.com/skvdmt/secret-contracts/pkg/go/gen/auth_error_v1"
)

// Repository Интерфейс репозиторного слоя.
type Repository interface {
	// Остановка.
	Stop(context.Context) error
	// Документация API.
	Documentation() string
	// Отправка ошибки аутентификации.
	AuthErrorSending(context.Context, *auth_error_v1.Event) error
	// Проверка при входе.
	LoginVerification(ctx context.Context, username, password string) (bool, error)
	// Сохранение токена обновления.
	RefreshTokenSaving(ctx context.Context, username, token string) error
	// Валидация токена обновления.
	RefreshTokenValidation(ctx context.Context, username, token string) (bool, error)
}
