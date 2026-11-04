package usecase

import (
	"context"
	"net/http"
	"unicode/utf8"

	"github.com/skvdmt/auth-back/internal/entities"
	"github.com/skvdmt/auth-back/internal/model"
	"github.com/skvdmt/auth-back/internal/repository"
	"github.com/skvdmt/jwt"
	"github.com/skvdmt/secret-contracts/pkg/go/gen/auth_error_v1"
	erw "github.com/skvdmt/skvdmt-back/pkg/errwrap"
)

// App Сервисный слой.
type App struct {
	// Репозиторный слой.
	repository Repository
}

// NewApp Конструктор.
func NewApp(ctx context.Context) (*App, error) {
	model.Log.Info.Info("usecase layer creating")
	a := &App{}
	var err error
	// Создание репозиторного слоя.
	a.repository, err = repository.NewApp(ctx)
	if err != nil {
		return nil, err
	}
	return a, nil
}

// Stop Остановка.
func (a *App) Stop(ctx context.Context) error {
	if err := a.repository.Stop(ctx); err != nil {
		return err
	}
	model.Log.Info.Info("usecase layer stopped")
	return nil
}

// Tokens Получение токенов.
func (a *App) Tokens(ctx context.Context, username, password string) (*entities.Tokens, error) {
	// Валидация длины имени пользователя.
	if utf8.RuneCountInString(username) < 3 ||
		utf8.RuneCountInString(username) > 15 {
		return nil, erw.New(
			erw.CodeHTTP(http.StatusUnauthorized),
			erw.Message("username must be from 3 to 15 characters"),
		)
	}
	// Валидация длины пароля.
	if utf8.RuneCountInString(password) < 8 ||
		utf8.RuneCountInString(password) > 20 {
		return nil, erw.New(
			erw.CodeHTTP(http.StatusUnauthorized),
			erw.Message("password must be from 8 to 20 characters"),
		)
	}
	return a.repository.Tokens(ctx, username, password)
}

// Refresh Обновление токенов.
func (a *App) Refresh(ctx context.Context, token string) (*entities.Tokens, error) {
	p, err := jwt.Parse(token)
	if err != nil {
		return nil, erw.New(
			erw.CodeHTTP(http.StatusUnauthorized),
		)
	}
	if v, _ := p.Valid(model.RefreshKey); !v {
		return nil, erw.New(
			erw.CodeHTTP(http.StatusUnauthorized),
		)
	}
	u, ok := p.Claim("usr")
	if !ok {
		return nil, erw.New(
			erw.CodeHTTP(http.StatusUnauthorized),
		)
	}
	return a.repository.Refresh(ctx, u, token)
}

// Documentation Документация.
func (a *App) Documentation(c context.Context) []byte {
	return a.repository.Documentation(c)
}

func (a *App) SendAuthError(c context.Context, e *auth_error_v1.Event) error {
	return a.repository.SendAuthError(c, e)
}
