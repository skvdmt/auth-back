package inserts

import (
	"database/sql"
	"errors"

	"github.com/lib/pq"
	"github.com/lib/pq/pqerror"
)

// InsertToPostgres Вставка в postgres
func InsertToPostgres(db *sql.DB) error {
	// Данные для тестов.
	_, err := db.Exec(`INSERT INTO users (username, password_hash) VALUES('dima', '$2a$12$pmbuRJbrzDQvpd/.zfGrjejTuhWpQufTXPI1WHA6Ozvg1NPz9H476')`)
	if err != nil && !IsPostgresErrorCode(err, pqerror.UniqueViolation) {
		return err
	}
	return nil
}

// IsPostgresErrorCode Поиск кода в ошибках postgres.
func IsPostgresErrorCode(err error, errcode pqerror.Code) bool {
	if pgErr, ok := errors.AsType[*pq.Error](err); ok {
		return pgErr.Code == errcode
	}
	return false
}
