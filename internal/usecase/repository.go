package usecase

import (
	"context"

	"github.com/skvdmt/auth-back/internal/entities"
)

// Repository Интерфейс репозиторного слоя.
type Repository interface {
	// Остановка.
	Stop(context.Context) error
	// Авторизация пользователя.
	Login(c context.Context, username, password string) (*entities.Tokens, error)
	// Обновление токенов.
	Refresh(c context.Context, username, token string) (*entities.Tokens, error)
}
