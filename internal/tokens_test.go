//go:build integration

package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/skvdmt/auth-back/internal/entities"
	"github.com/skvdmt/auth-back/internal/model"
	_ "github.com/skvdmt/auth-back/testing_init"
)

// TestTokens Тестирование получения токенов.
func TestTokens(t *testing.T) {
	p, err := NewPrepareApp()
	if err != nil {
		log.Fatal(err)
	}
	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := p.Stop(); err != nil {
			log.Fatal(err)
		}
	}()
	textTests := []struct {
		name           string
		login          string
		password       string
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "correct login and pasword",
			login:          "dima",
			password:       "password",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "short login",
			login:          "qq",
			password:       "password",
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "incorrect login or password",
		},
		{
			name:           "short password",
			login:          "dima",
			password:       "qqq",
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "incorrect login or password",
		},
		{
			name:           "incorrect login or password",
			login:          "dima2",
			password:       "password2",
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "incorrect login or password",
		},
	}
	c := http.Client{}
	for _, test := range textTests {
		t.Run(test.name, func(t *testing.T) {
			// Запрос
			req, err := http.NewRequest(
				"GET",
				fmt.Sprintf(
					"http://localhost:%d/tokens",
					model.Conf.Server.Port,
				),
				nil,
			)
			if err != nil {
				log.Fatal(err)
			}
			req.Header.Set("X-Forwarded-For", "192.168.1.1")
			req.SetBasicAuth(test.login, test.password)
			resp, err := c.Do(req)
			if err != nil {
				log.Fatal(err)
			}
			// Проверка статуса
			if test.expectedStatus != resp.StatusCode {
				t.Errorf(
					"error: status code not equal\n expected: %d\n actual: %d\n",
					test.expectedStatus, resp.StatusCode)
			}
			// Получение тела ответа
			buf := new(bytes.Buffer)
			_, err = buf.ReadFrom(resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			if resp.StatusCode == http.StatusOK {
				const SET_COOKIE = "Set-Cookie"
				rt := resp.Header.Get(SET_COOKIE)
				if len(rt) == 0 {
					t.Errorf("\n\nerror: refresh token not set\n\n")
				}
				var tks entities.Tokens
				if err := json.Unmarshal(buf.Bytes(), &tks); err != nil {
					log.Fatal(err)
				}
				// Проверка длинны токена доступа.
				if len(tks.AccessToken) == 0 {
					t.Errorf("\n\nerror: access token are empty\n\n")
				}
			}
			if resp.StatusCode == http.StatusUnauthorized {
				var s struct {
					Error string `json:"error"`
				}
				if err := json.Unmarshal(buf.Bytes(), &s); err != nil {
					log.Fatal(err)
				}
				if test.expectedError != s.Error {
					t.Errorf(
						"error: error message not equal\n expected: %s\n actual: %s\n",
						test.expectedError, s.Error)
				}
			}
		})
	}
}
