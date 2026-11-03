package internal

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/skvdmt/auth-back/internal/delivery"
	"github.com/skvdmt/auth-back/internal/model"
)

// App Основная структура приложения.
type App struct {
	// Канал сигналов остановки.
	interrupt chan os.Signal
	// Контекст приложения.
	ctx context.Context
	// Отмена контекста.
	cancel context.CancelFunc
	// Транспортный слой.
	delivery Delivery
	// Сервер авторизации.
	server *http.Server
	// HTTP Роутер.
	router *http.ServeMux
	// Ошибки приложения.
	errors []error
	// Приложение уже закрывается.
	stopping bool
	// Корректное завершение горутин.
	wg *sync.WaitGroup
	// Листнер готов.
	listenerReady bool
}

// NewApp Конструктор.
func NewApp() (*App, error) {
	var err error
	// Загрузка конфигурации.
	if model.Conf, err = model.NewConfig(); err != nil {
		return nil, err
	}
	model.Log.Info.Info(fmt.Sprintf("API url %s", model.Conf.Links.Api))
	model.Log.Info.Info(fmt.Sprintf("Documentation %s", model.Conf.Links.Documentation))

	model.Log.Info.Info(fmt.Sprintf("app %s creating", model.NAME))
	a := &App{
		interrupt: make(chan os.Signal),
		wg:        &sync.WaitGroup{},
	}
	// Создане глобальных каналов ошибок для всего приложения.
	model.Errors = make(chan error)
	model.ExternalErrors = make(chan error)
	// Создание контекста.
	a.ctx, a.cancel = context.WithCancel(context.Background())
	// Ключи.
	if err := model.CreateKeys(); err != nil {
		return nil, err
	}
	// Создание транспортного слоя.
	a.delivery, err = delivery.NewApp(a.ctx)
	if err != nil {
		return nil, err
	}
	// Создание сервера аутентификации.
	a.router = http.NewServeMux()
	a.routing()
	a.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", model.Conf.Server.Port),
		Handler: a.router,
	}
	return a, nil
}

// routing Настройка маршрутов.
func (a *App) routing() {
	a.router.HandleFunc("GET /tokens", a.delivery.Tokens)
	a.router.HandleFunc("GET /refresh", a.delivery.Refresh)
	a.router.HandleFunc("DELETE /refresh", a.delivery.DeleteRefresh)
	a.router.HandleFunc("GET /openapi.yaml", a.delivery.Documentation)
}

// Start Запуск приложения.
func (a *App) Start() []error {
	model.Log.Info.Info(fmt.Sprintf("app %s starting", model.NAME))

	// Запуск обработчика ошибок.
	a.wg.Add(1)
	go a.errorHandler()

	// Запуск сервера авторизации.
	a.wg.Add(1)
	go a.authServerStart()

	return a.interruptHandler()
}

// Ready Сервер готов обрабатывать запросы.
func (a *App) Ready() bool {
	return a.listenerReady
}

// HealthCheckLock Метод блокирует выполнение пока
// сервер не будет готов обрабатывать запросы.
func (a *App) HealthCheckLock() error {
	const interval = time.Millisecond * 10
	const timeout = time.Second * 3
	s := time.Now()
	for {
		if a.Ready() {
			return nil
		}
		time.Sleep(interval)
		if time.Since(s) > timeout {
			return fmt.Errorf("health check timeout")
		}
	}
}

// errorHandler Обработчик глобального канала ошибок.
func (a *App) errorHandler() {
	model.Log.Info.Info("error handle starting")
	defer func() {
		model.Log.Info.Info("error handle stopped")
		a.wg.Done()
	}()
	for {
		select {
		case err := <-model.Errors:
			if err == nil {
				return
			}
			if errors.Is(err, context.Canceled) {
				continue
			}
			a.errors = append(a.errors, err)
			if !a.stopping {
				a.interrupt <- syscall.SIGTERM
			}
		case err := <-model.ExternalErrors:
			if err == nil {
				return
			}
			if errors.Is(err, context.Canceled) {
				continue
			}
			model.Log.Error.Error(err.Error())
		}
	}
}

// authServerStart Запуск сервера авторизации.
func (a *App) authServerStart() {
	defer func() {
		model.Log.Info.Info("auth http server stopped")
		a.wg.Done()
	}()
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", model.Conf.Server.Port))
	if err != nil {
		model.Errors <- err
		return
	}
	model.Log.Info.Info(fmt.Sprintf("auth server starting on %d port",
		model.Conf.Server.Port))
	a.listenerReady = true
	if err := a.server.Serve(l); err != nil &&
		!errors.Is(err, http.ErrServerClosed) {
		model.Errors <- err
	}
}

// interruptHandler Обработчик сигналов остановки приложения.
func (a *App) interruptHandler() []error {
	model.Log.Info.Info("waiting for interrupt signals")
	signal.Notify(a.interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-a.interrupt
	a.stop()
	// Закрытие каналов ошибок.
	close(model.ExternalErrors)
	close(model.Errors)
	// Ожидание завершения горутин.
	a.wg.Wait()
	return a.errors
}

// stop Остановка.
func (a *App) stop() {
	defer model.Log.Info.Info(fmt.Sprintf("app %s stopped", model.NAME))
	a.stopping = true
	// Отмена контекста.
	a.cancel()
	model.Log.Info.Info("context canceled")
	// Остановка сервера авторизации.
	if err := a.server.Shutdown(context.Background()); err != nil {
		model.Errors <- err
	}
	// Остановка транспортного слоя.
	if err := a.delivery.Stop(a.ctx); err != nil {
		model.Errors <- err
	}
	// Закрытие канала сигналов прерывания.
	close(a.interrupt)
}
