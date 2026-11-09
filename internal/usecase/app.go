package usecase

import (
	"context"
	"crypto/rand"
	"math/big"
	"net/http"
	"time"
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
		utf8.RuneCountInString(username) > 15 ||
		utf8.RuneCountInString(password) < 8 ||
		utf8.RuneCountInString(password) > 20 {
		// timing attack securety
		time.Sleep(150*time.Millisecond + time.Duration(a.rand(50, 100))*time.Millisecond)
		return nil, erw.New(
			erw.CodeHTTP(http.StatusUnauthorized),
			erw.Message("incorrect login or password"),
		)
	}
	v, err := a.repository.LoginVerification(ctx, username, password)
	if err != nil {
		return nil, erw.New(
			erw.CodeHTTP(http.StatusInternalServerError),
			erw.Internal(
				erw.Location("repository.LoginVerification"),
				erw.Error(err),
			),
		)
	}
	if !v {
		return nil, erw.New(
			erw.CodeHTTP(http.StatusUnauthorized),
			erw.Message("incorrect login or password"),
		)
	}
	return a.newTokens(ctx, username)
}

// Refresh Обновление токенов.
func (a *App) Refresh(ctx context.Context, token string) (*entities.Tokens, error) {
	t, err := jwt.Parse(token)
	if err != nil {
		return nil, erw.New(
			erw.CodeHTTP(http.StatusUnauthorized),
			erw.Message("invalid refresh token"),
		)
	}
	if v, _ := t.Valid(model.RefreshKey); !v {
		return nil, erw.New(
			erw.CodeHTTP(http.StatusUnauthorized),
			erw.Message("invalid refresh token"),
		)
	}
	u, err := t.Username()
	if err != nil {
		return nil, erw.New(
			erw.CodeHTTP(http.StatusUnauthorized),
			erw.Message("invalid refresh token"),
		)
	}
	v, err := a.repository.RefreshTokenValidation(ctx, u, token)
	if err != nil {
		return nil, erw.New(
			erw.CodeHTTP(http.StatusInternalServerError),
			erw.Internal(
				erw.Location("repository.RefreshTokenValidation"),
				erw.Error(err),
			),
		)
	}
	if !v {
		return nil, erw.New(
			erw.CodeHTTP(http.StatusUnauthorized),
			erw.Message("invalid refresh token"),
		)
	}
	return a.newTokens(ctx, u)
}

// Documentation Документация.
func (a *App) Documentation() string {
	return a.repository.Documentation()
}

// SendAuthError Отправка ошибки аутентификации.
func (a *App) SendAuthError(c context.Context, e *auth_error_v1.Event) error {
	return a.repository.AuthErrorSending(c, e)
}

// newTokens Создание пары токенов и сохранение токена обновления.
func (a *App) newTokens(ctx context.Context, username string) (*entities.Tokens, error) {
	t := entities.NewTokens(
		a.newToken(model.Conf.AccessTokenTTL, username, model.AccessKey),
		a.newToken(model.Conf.RefreshTokenTTL, username, model.RefreshKey),
	)
	if err := a.repository.RefreshTokenSaving(ctx, username, t.RefreshToken()); err != nil {
		return nil, erw.New(
			erw.CodeHTTP(http.StatusInternalServerError),
			erw.Internal(
				erw.Location("repository.RefreshTokenSaving"),
				erw.Error(err),
			),
		)
	}
	return t, nil
}

// newToken Новый токен.
func (a *App) newToken(exp time.Duration, username, key string) string {
	return jwt.NewToken(
		jwt.Headers(
			jwt.Header("alg", "HS256"),
			jwt.Header("typ", "JWT"),
		),
		jwt.Claims(
			jwt.Claim("nbf", time.Now().Format(time.RFC3339)),
			jwt.Claim("iat", time.Now().Format(time.RFC3339)),
			jwt.Claim("exp", time.Now().Add(exp).Format(time.RFC3339)),
			jwt.Claim("sub", username),
		),
		key,
	)
}

// rand Случайное число.
func (a *App) rand(bottom, upper int64) int64 {
	if bottom == upper {
		return bottom
	}
	if bottom > upper {
		bottom, upper = upper, bottom
	}
	diff := upper - bottom + 1
	rnd := big.NewInt(diff)
	num, err := rand.Int(rand.Reader, rnd)
	if err != nil {
		return bottom + (upper-bottom)/2
	}
	return num.Int64() + bottom
}
