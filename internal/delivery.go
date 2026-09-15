package internal

import (
	"context"
	"net/http"
)

// Delivery Интерфейс транспортного слоя.
type Delivery interface {
	// Остановка.
	Stop(context.Context) error
	// Получение токенов.
	Tokens(http.ResponseWriter, *http.Request)
	// Удаление токена обновления.
	DeleteRefresh(http.ResponseWriter, *http.Request)
	// Обновление токенов.
	Refresh(http.ResponseWriter, *http.Request)
	// Документация.
	Documentation(http.ResponseWriter, *http.Request)
}
