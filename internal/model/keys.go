package model

import (
	"fmt"
	"os"
	"unicode/utf8"
)

const (
	ACCESS_KEY  = "ACCESS_KEY"
	REFRESH_KEY = "REFRESH_KEY"
	KEY_LEN     = 32
)

var (
	AccessKey  = ""
	RefreshKey = ""
)

// CreateKeys Создание ключей.
func CreateKeys() error {
	var a, r string
	var err error
	if a, err = parseKeyEnv(ACCESS_KEY); err != nil {
		return err
	}
	if r, err = parseKeyEnv(REFRESH_KEY); err != nil {
		return err
	}
	AccessKey = a
	RefreshKey = r
	return nil
}

// parseKeyEnv Парсинг ключа из переменной окружения.
func parseKeyEnv(key string) (string, error) {
	v, o := os.LookupEnv(key)
	if !o {
		return "", fmt.Errorf("env %s is not set", key)
	}
	if utf8.RuneCountInString(v) != KEY_LEN {
		return "", fmt.Errorf("env %s must be %d chars length", key, KEY_LEN)
	}
	return v, nil
}
