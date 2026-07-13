package usecase

import (
	"context"
	"net/http"
	"unicode/utf8"

	"github.com/skvdmt/auth-back/internal/entities"
	"github.com/skvdmt/auth-back/internal/model"
	"github.com/skvdmt/auth-back/internal/repository"
	"github.com/skvdmt/erw"
	"github.com/skvdmt/jwt"
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

// Login Авторизация пользователя.
func (a *App) Login(ctx context.Context, username, password string) (*entities.Tokens, error) {
	// Валидация длинны логина.
	if utf8.RuneCountInString(username) < 3 ||
		utf8.RuneCountInString(username) > 15 {
		return nil, erw.New(
			erw.Status(http.StatusBadRequest),
			erw.Response("username must be from 3 to 15 characters"),
		)
	}
	// Валидация длинны пароля.
	if utf8.RuneCountInString(password) < 8 ||
		utf8.RuneCountInString(password) > 20 {
		return nil, erw.New(
			erw.Status(http.StatusBadRequest),
			erw.Response("password must be from 8 to 20 characters"),
		)
	}
	return a.repository.Login(ctx, username, password)
}

// Refresh Обновление токенов.
func (a *App) Refresh(ctx context.Context, token string) (*entities.Tokens, error) {
	p, err := jwt.Parse(token)
	if err != nil {
		return nil, erw.New(
			erw.Status(http.StatusUnauthorized),
		)
	}
	if v, _ := p.Valid(model.RefreshKey); !v {
		return nil, erw.New(
			erw.Status(http.StatusUnauthorized),
		)
	}
	u, ok := p.Claim("usr")
	if !ok {
		return nil, erw.New(
			erw.Status(http.StatusUnauthorized),
		)
	}
	return a.repository.Refresh(ctx, u, token)
}
