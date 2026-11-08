package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"os"

	"github.com/lib/pq"
	"github.com/skvdmt/auth-back/init/inserts"
	"github.com/skvdmt/auth-back/internal/model"
	"golang.org/x/crypto/bcrypt"
)

const (
	POSTGRES_PASSWORD_ENV = "POSTGRES_PASSWORD"
)

// Репозиторий Postgres.
type Posgres struct {
	db *sql.DB
}

// NewPostgres Конструктор.
func NewPostgres() (*Posgres, error) {
	model.Log.Info.Info("postgres database connection opening")
	p, o := os.LookupEnv(POSTGRES_PASSWORD_ENV)
	if !o {
		return nil, fmt.Errorf("env %s is not set", POSTGRES_PASSWORD_ENV)
	}
	pwd, err := url.QueryUnescape(p)
	if err != nil {
		pwd = p
	}
	c, err := pq.NewConnectorConfig(pq.Config{
		User:     model.Conf.Postgres.User,
		Password: pwd,
		Host:     model.Conf.Postgres.Host,
		Port:     model.Conf.Postgres.Port,
		Database: model.Conf.Postgres.Database,
		SSLMode:  pq.SSLModeDisable,
	})
	if err != nil {
		return nil, err
	}
	db := sql.OpenDB(c)
	if err := db.Ping(); err != nil {
		return nil, err
	}
	if err := inserts.InsertToPostgres(db); err != nil {
		return nil, err
	}
	model.Log.Info.Info("postgres database connection success")
	return &Posgres{
		db: db,
	}, nil
}

// LoginVerification Проверка при входе.
func (p *Posgres) LoginVerification(ctx context.Context, username, password string) (bool, error) {
	var hash string
	if err := p.db.QueryRowContext(ctx, `SELECT password_hash FROM users WHERE username = $1`, username).Scan(&hash); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil, nil
}

// Stop Отстановка.
func (p *Posgres) Stop() error {
	return p.db.Close()
}
