package delivery

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/skvdmt/auth-back/internal/model"
	"github.com/skvdmt/auth-back/internal/usecase"
	"github.com/skvdmt/erw"
)

// App Транспортный слой.
type App struct {
	// Сервисный слой.
	usecase Usecase
}

// NewApp Конструктор.
func NewApp(ctx context.Context) (*App, error) {
	model.Log.Info.Info("delivery layer creating")
	a := &App{}
	var err error
	// Создание сервисного слоя.
	a.usecase, err = usecase.NewApp(ctx)
	if err != nil {
		return nil, err
	}
	return a, nil
}

// Stop Остановка.
func (a *App) Stop(ctx context.Context) error {
	if err := a.usecase.Stop(ctx); err != nil {
		return err
	}
	model.Log.Info.Info("delivery layer stopped")
	return nil
}

// Login Аутентификация пользователя.
func (a *App) Login(w http.ResponseWriter, r *http.Request) {
	model.Log.Info.Info("Auth")
	u := r.URL.Query().Get("username")
	p := r.URL.Query().Get("password")
	t, err := a.usecase.Login(r.Context(), u, p)
	if err != nil {
		a.errorHandle(w, err)
		return
	}
	a.setRefreshTokenCookie(w, t.RefreshToken())
	a.response(w, http.StatusOK, t)
}

// Refresh Обновление токенок.
func (a *App) Refresh(w http.ResponseWriter, r *http.Request) {
	model.Log.Info.Info("RefreshTokens")
	c, err := r.Cookie("refresh_token")
	if err != nil {
		fmt.Println(err)
		a.errorHandle(w, erw.New(
			erw.Status(http.StatusUnauthorized),
		))
		return
	}
	n, err := a.usecase.Refresh(r.Context(), c.Value)
	if err != nil {
		a.errorHandle(w, err)
		return
	}
	a.setRefreshTokenCookie(w, n.RefreshToken())
	a.response(w, http.StatusOK, n)
}

// errorHandle Обработка ошибок.
func (a *App) errorHandle(w http.ResponseWriter, e error) {
	err, ok := e.(*erw.ErrorWrapper)
	if !ok {
		model.Errors <- fmt.Errorf("error convertion faled; error: %v", e)
		return
	}
	// Логирование ошибки.
	if err.LoggingError() != nil {
		model.Log.Error.Error(err.LoggingErrorText())
	}
	// Отправка ответа клиенту.
	a.response(w, err.StatusCode(), err.ResponseJSON())
}

// response Отправить ответ.
func (a *App) response(w http.ResponseWriter, status int, resp any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Vary", "Origin")
	w.WriteHeader(status)
	b, ok := resp.([]byte)
	if ok {
		w.Write(b)
		return
	}
	json.NewEncoder(w).Encode(resp)
}

// setRefreshTokenCookie Запись в cookie токена обновления.
func (a *App) setRefreshTokenCookie(w http.ResponseWriter, t string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    t,
		Domain:   "localhost",
		Path:     "/",
		Expires:  time.Now().AddDate(0, 0, 7),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	})
}
