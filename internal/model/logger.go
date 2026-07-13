package model

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

const (
	// Путь к директории журналов в произвозстве.
	LOG_DIRECTORY_PROD = "/var/log"

	// Путь к директории журналов в разработке.
	LOG_DIRECTORY_DEV = "./logs"

	// Имя файла журнала ошибок.
	LOG_ERROR_FILENAME = "error.log"

	// Флаги файла журнала ошибок.
	LOG_ERROR_FLAGS = os.O_CREATE | os.O_APPEND | os.O_RDWR
)

// Log Глобальный экземпляр журнала.
var Log *Logger

// Logger Журнал.
type Logger struct {
	// Журнал информирования.
	Info *slog.Logger
	// Журнал ошибок.
	Error *slog.Logger
	// Файл журнала ошибок.
	errorFile *os.File
}

// NewLogger Конструктор.
func NewLogger() (*Logger, error) {
	e, err := os.OpenFile(
		filepath.Join(logDir(), LOG_ERROR_FILENAME),
		LOG_ERROR_FLAGS,
		os.ModePerm,
	)
	if err != nil {
		return nil, err
	}
	return &Logger{
		Info:      slog.New(slog.NewTextHandler(os.Stdout, nil)),
		Error:     slog.New(slog.NewJSONHandler(io.MultiWriter(e, os.Stderr), nil)),
		errorFile: e,
	}, nil
}

// Close Закрытие журнала.
func (l *Logger) Close() error {
	if l.errorFile != nil {
		return l.errorFile.Close()
	}
	return nil
}

// logDir Директоия журналов.
func logDir() string {
	m, o := os.LookupEnv(MODE)
	if o && m == MODE_DEV {
		return LOG_DIRECTORY_DEV
	}
	return filepath.Join(LOG_DIRECTORY_PROD, NAME)
}
