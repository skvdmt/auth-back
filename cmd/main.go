package main

import (
	"fmt"
	"os"

	"github.com/skvdmt/auth-back/internal"
	"github.com/skvdmt/auth-back/internal/model"
)

// Точка входа в приложение.
func main() {
	var err error
	// Создание логгера.
	model.Log, err = model.NewLogger()
	if err != nil {
		panic(err)
	}
	defer func() {
		// Закрытие логгера.
		if err := model.Log.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "logger close error: %v\n", err)
		}
	}()
	// Создание приложения.
	app, err := internal.NewApp()
	if err != nil {
		model.Log.Error.Error(err.Error())
		return
	}
	// Запуск приложения.
	if errs := app.Start(); len(errs) > 0 {
		for _, err := range errs {
			model.Log.Error.Error(err.Error())
		}
	}
}
