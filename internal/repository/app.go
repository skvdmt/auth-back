package repository

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	_ "github.com/lib/pq"
	"github.com/skvdmt/auth-back/internal/entities"
	"github.com/skvdmt/auth-back/internal/model"
	"github.com/skvdmt/erw"
	"github.com/skvdmt/jwt"
)

const (
	PKG               = "repository"
	USER_NAME         = "dima"
	USER_PASSWORD     = "password"
	CLEANING_INTERVAL = time.Minute * 1

	ACCESS_TOKEN_VALID       = time.Minute * 15
	REFRESH_TOKEN_VALID_DAYS = 7
	REFRESH_TOKEN_USER_LIMIT = 10
)

// App Репоситорный слой.
type App struct {
	// Соединение с базой данных.
	db            *sql.DB
	cleaningTiker *time.Ticker
	stopCleaning  chan struct{}
	wg            *sync.WaitGroup
	documentation []byte
}

// NewApp Конструктор.
func NewApp(ctx context.Context) (*App, error) {
	model.Log.Info.Info("repository layer creating")
	a := &App{
		cleaningTiker: time.NewTicker(CLEANING_INTERVAL),
		stopCleaning:  make(chan struct{}, 1),
		wg:            &sync.WaitGroup{},
	}
	// Соединение с базой данных.
	var err error
	if a.db, err = OpenPostgresDB(); err != nil {
		return nil, err
	}
	// Обновление документации
	if err := a.updateDocumentation(); err != nil {
		return nil, err
	}
	// Разовая очистка.
	if err := a.clearRefreshTokens(); err != nil {
		return nil, err
	}
	// Обработчик очистки.
	a.wg.Add(1)
	go a.cleaningHandler()
	return a, nil
}

// Stop Остановка.
func (a *App) Stop(ctx context.Context) error {
	// Остановка и закрытие тикера очистки.
	a.cleaningTiker.Stop()
	// Остановка обработчика тикера очистки.
	close(a.stopCleaning)

	// Ожидание завершения горутин.
	a.wg.Wait()

	// Закрытие соединения с базой данных.
	if err := a.closeDB(); err != nil {
		return err
	}

	model.Log.Info.Info("repository layer stopped")
	return nil
}

// Tokens Получение токенов.
func (a *App) Tokens(ctx context.Context, username, password string) (*entities.Tokens, error) {
	if username != USER_NAME || password != USER_PASSWORD {
		return nil, erw.New(
			erw.Status(http.StatusUnauthorized),
			erw.Response("incorrect login or password"),
		)
	}
	return a.createTokens(username)
}

// Refresh Обновление токенов.
func (a *App) Refresh(ctx context.Context, username, token string) (*entities.Tokens, error) {
	q := `DELETE FROM ONLY refresh_tokens WHERE id IN
	 (SELECT id FROM refresh_tokens
	  WHERE username = $1 AND refresh_token = $2
		 ORDER BY created_at LIMIT 1);`
	as := []any{username, token}
	r, err := a.db.Exec(q, as...)
	if err != nil {
		return nil, erw.New(
			erw.Status(http.StatusUnauthorized),
		)
	}
	c, err := r.RowsAffected()
	if err != nil {
		return nil, erw.New(
			erw.Status(http.StatusInternalServerError),
			erw.Logging(
				erw.Location(PKG, "RefreshTokens", "RowsAffected"),
				erw.Error(err),
				erw.SQL(q, as),
			),
		)
	}
	if c != 1 {
		return nil, erw.New(
			erw.Status(http.StatusUnauthorized),
		)
	}
	return a.createTokens(username)
}

// Documentation Документация.
func (a *App) Documentation() []byte {
	return a.documentation
}

// updateDocumentation Обновление документации.
func (a *App) updateDocumentation() error {
	const (
		name = "openapi.yaml"
		dev  = "."
		prod = "/usr/local/share/doc"
	)
	m, o := os.LookupEnv(model.MODE)
	dir := prod
	if o && m == model.MODE_DEV {
		dir = dev
	}
	doc, err := os.ReadFile(filepath.Join(dir, name))
	if err != nil {
		return err
	}
	a.documentation = doc
	return nil
}

// cleaningHandler Обработчик очистки.
func (a *App) cleaningHandler() {
	model.Log.Info.Info("cleaning handler starting")
	defer func() {
		model.Log.Info.Info("cleaning handler stopped")
		a.wg.Done()
	}()
	for {
		select {
		case <-a.cleaningTiker.C:
			if err := a.clearRefreshTokens(); err != nil {
				model.Errors <- err
			}
		case <-a.stopCleaning:
			return
		}
	}
}

// clearRefreshTokens Очистка токенов обновления.
func (a *App) clearRefreshTokens() error {
	q := fmt.Sprintf(`DELETE FROM refresh_tokens
	 WHERE created_at < (now() AT TIME ZONE 'Europe/Moscow') - interval '%d days';`,
		REFRESH_TOKEN_VALID_DAYS)
	_, err := a.db.Exec(q)
	if err != nil {
		return err
	}
	model.Log.Info.Info("clear refresh tokens complete")
	return nil
}

// createTokens Созжание токенов.
func (a *App) createTokens(username string) (*entities.Tokens, error) {
	// Оставить не более лимита самых свежих токенов обновления пользователя.
	q := fmt.Sprintf(`DELETE FROM ONLY refresh_tokens WHERE id IN
	 (SELECT id FROM refresh_tokens WHERE username = $1 ORDER BY created_at DESC OFFSET %d)`,
		REFRESH_TOKEN_USER_LIMIT-1)
	as := []any{username}
	_, err := a.db.Exec(q, as...)
	if err != nil {
		return nil, erw.New(
			erw.Status(http.StatusInternalServerError),
			erw.Logging(
				erw.Location(PKG, "createTokens", "delete old refresh tokens"),
				erw.SQL(q, as...),
			),
		)
	}
	rt := jwt.NewToken(
		jwt.Headers(
			jwt.Header("alg", "HS256"),
			jwt.Header("typ", "JWT"),
		),
		jwt.Claims(
			jwt.Claim("nbf", time.Now().Format(time.RFC3339)),
			jwt.Claim("iat", time.Now().Format(time.RFC3339)),
			jwt.Claim("exp", time.Now().AddDate(0, 0, REFRESH_TOKEN_VALID_DAYS).Format(time.RFC3339)),
			jwt.Claim("usr", username),
		),
		model.RefreshKey,
	)
	if err := a.saveRefreshToken(username, rt); err != nil {
		return nil, err
	}
	return entities.NewTokens(jwt.NewToken(
		jwt.Headers(
			jwt.Header("alg", "HS256"),
			jwt.Header("typ", "JWT"),
		),
		jwt.Claims(
			jwt.Claim("nbf", time.Now().Format(time.RFC3339)),
			jwt.Claim("iat", time.Now().Format(time.RFC3339)),
			jwt.Claim("exp", time.Now().Add(ACCESS_TOKEN_VALID).Format(time.RFC3339)),
			jwt.Claim("usr", username),
		),
		model.AccessKey,
	), rt), nil
}

// saveRefreshToken Сохранение токена обновления.
func (a *App) saveRefreshToken(u, t string) error {
	q := `INSERT INTO refresh_tokens (username, refresh_token) values($1, $2);`
	as := []any{u, t}
	_, err := a.db.Exec(q, as...)
	if err != nil {
		return erw.New(
			erw.Status(http.StatusInternalServerError),
			erw.Logging(
				erw.Location(PKG, "saveRefreshToken", "Exec"),
				erw.Error(err),
				erw.SQL(q, as),
			),
		)
	}
	return nil
}

// closeDB Закрывает соединение с базой данных.
func (a *App) closeDB() error {
	defer model.Log.Info.Info("database connection closed")
	return a.db.Close()
}
