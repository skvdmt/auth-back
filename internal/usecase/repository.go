package usecase

import (
	"context"

	"github.com/skvdmt/auth-back/internal/entities"
	"github.com/skvdmt/secret-contracts/pkg/go/gen/auth_error_v1"
)

// Repository Интерфейс репозиторного слоя.
type Repository interface {
	// Остановка.
	Stop(context.Context) error
	// Получение токенов.
	Tokens(c context.Context, username, password string) (*entities.Tokens, error)
	// Обновление токенов.
	Refresh(c context.Context, username, token string) (*entities.Tokens, error)
	// Документация.
	Documentation(c context.Context) []byte
	// Отправка ошибки аутентификации.
	SendAuthError(c context.Context, e *auth_error_v1.Event) error
}
