package main

import (
	"fmt"
	"os"
	"unicode/utf8"

	"github.com/skvdmt/auth-back/internal"
	"github.com/skvdmt/auth-back/internal/model"
)

// main Точка входа в приложение.
func main() {

	fmt.Println(utf8.RuneCountInString("GHCLADSDD5uTq4UgO3ffTNR4i4r6hBt4"))
	os.Exit(1)

	// Создание логгера.
	var err error
	if model.Log, err = model.NewLogger(); err != nil {
		panic(err)
	}
	// Создание приложения.
	app, err := internal.NewApp()
	if err != nil {
		model.Log.Error.Error(err.Error())
		os.Exit(1)
	}
	// Запуск приложения.
	if errs := app.Start(); len(errs) > 0 {
		for _, err := range errs {
			model.Log.Error.Error(err.Error())
		}
	}
	// Закрытие логгера.
	if err := model.Log.Close(); err != nil {
		panic(err)
	}
}
