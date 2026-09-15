package delivery

import (
	"context"

	"github.com/skvdmt/auth-back/internal/entities"
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
	Documentation() []byte
}
