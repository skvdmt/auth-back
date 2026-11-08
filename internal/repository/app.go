package repository

import (
	"context"
	"time"

	_ "github.com/lib/pq"
	"github.com/skvdmt/auth-back/internal/model"
)

const (
	PKG               = "repository"
	USER_NAME         = "dima"
	USER_PASSWORD     = "password"
	CLEANING_INTERVAL = time.Minute * 1
)

// App Репоситорный слой.
type App struct {
	*Posgres
	*Kafka
	*Redis
	*Files
}

// NewApp Конструктор.
func NewApp(_ context.Context) (*App, error) {
	model.Log.Info.Info("repository layer creating")
	p, err := NewPostgres()
	if err != nil {
		return nil, err
	}
	k, err := NewKafka()
	if err != nil {
		return nil, err
	}
	r, err := NewRedis()
	if err != nil {
		return nil, err
	}
	f, err := NewFiles()
	if err != nil {
		return nil, err
	}
	return &App{
		Posgres: p,
		Kafka:   k,
		Redis:   r,
		Files:   f,
	}, nil
}

// Stop Остановка.
func (a *App) Stop(_ context.Context) error {
	if err := a.Kafka.Stop(); err != nil {
		return err
	}
	if err := a.Posgres.Stop(); err != nil {
		return err
	}
	if err := a.Redis.Stop(); err != nil {
		return err
	}
	model.Log.Info.Info("repository layer stopped")
	return nil
}
