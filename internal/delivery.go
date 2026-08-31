package internal

import (
	"context"
	"net/http"
)

// Delivery Интерфейс транспортного слоя.
type Delivery interface {
	// Остановка.
	Stop(context.Context) error
	// Авторизация пользователя.
	Login(http.ResponseWriter, *http.Request)
	// Разлогиниться.
	Logout(http.ResponseWriter, *http.Request)
	// Обновление токенов.
	Refresh(http.ResponseWriter, *http.Request)
}
