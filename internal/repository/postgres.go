package repository

import (
	"database/sql"
	"fmt"
	"net/url"
	"os"

	"github.com/lib/pq"
	"github.com/lib/pq/pqerror"
	"github.com/skvdmt/auth-back/internal/model"
)

const (
	POSTGRES_PASSWORD = "POSTGRES_PASSWORD"
)

// OpenPostgresDB Открывает соединение с базой данных postgres.
func OpenPostgresDB() (*sql.DB, error) {
	model.Log.Info.Info("postgres database connection opening")
	p, o := os.LookupEnv(POSTGRES_PASSWORD)
	if !o {
		return nil, fmt.Errorf("env %s is not set", POSTGRES_PASSWORD)
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
	model.Log.Info.Info("postgres database connection success")
	return db, nil
}

// IsPostgresErrorCode Поиск кода в ошибках postgres.
func IsPostgresErrorCode(err error, errcode pqerror.Code) bool {
	if pgerr, ok := err.(*pq.Error); ok {
		return pgerr.Code == errcode
	}
	return false
}
