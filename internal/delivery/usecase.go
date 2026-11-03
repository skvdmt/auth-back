package delivery

import (
	"context"

	"github.com/skvdmt/auth-back/internal/entities"
	"github.com/skvdmt/auth-back/internal/go/gen/auth_error_v1"
)

// Usecase Интерфейс сервисного слоя.
type Usecase interface {
	// Остановка.
	Stop(context.Context) error
	// Получение токенов.
	Tokens(c context.Context, username, password string) (*entities.Tokens, error)
	// Обновление токенов.
	Refresh(c context.Context, token string) (*entities.Tokens, error)
	// Документация.
	Documentation(c context.Context) []byte
	// Отправка ошибки аутентификации.
	SendAuthError(c context.Context, e *auth_error_v1.Event) error
}
