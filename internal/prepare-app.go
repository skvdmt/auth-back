package internal

import (
	"os"
	"sync"
	"syscall"

	"github.com/skvdmt/auth-back/internal/model"
)

// PrepareApp Подготовка приложения к тестированию.
type PrepareApp struct {
	app *App
	wg  *sync.WaitGroup
}

// NewPrepareApp Конструктор.
func NewPrepareApp() (*PrepareApp, error) {
	// Создание логгера.
	var err error
	model.Log, err = model.NewLogger()
	if err != nil {
		return nil, err
	}
	// Создание приложения.
	app, err := NewApp()
	if err != nil {
		return nil, err
	}
	return &PrepareApp{
		app: app,
		wg:  &sync.WaitGroup{},
	}, nil
}

// Start Запуск
func (p *PrepareApp) Start() error {
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		if errs := p.app.Start(); len(errs) > 0 {
			for _, err := range errs {
				model.Log.Error.Error(err.Error())
			}
			os.Exit(1)
		}
	}()
	return p.app.HelthCheckLock()
}

// Stop Остановка.
func (p *PrepareApp) Stop() error {
	// Остановка сервера.
	p.app.interrupt <- syscall.SIGTERM
	p.wg.Wait()
	// Закрытие логгера.
	if err := model.Log.Close(); err != nil {
		return err
	}
	return nil
}
