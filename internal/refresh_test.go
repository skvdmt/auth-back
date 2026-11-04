//go:build integration

package internal

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"regexp"
// 	"testing"

// 	"github.com/skvdmt/auth-back/internal/entities"
// 	"github.com/skvdmt/auth-back/internal/model"
// 	_ "github.com/skvdmt/auth-back/testing_init"
// )

// // TestRefresh Тестирование обновления токенов.
// func TestRefresh(t *testing.T) {
// 	p, err := NewPrepareApp()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	if err := p.Start(); err != nil {
// 		log.Fatal(err)
// 	}
// 	defer func() {
// 		if err := p.Stop(); err != nil {
// 			log.Fatal(err)
// 		}
// 	}()

// 	t.Run("refresh tokens", func(t *testing.T) {
// 		c := http.Client{}
// 		// Получение токена обновления
// 		req, err := http.NewRequest(
// 			"GET",
// 			fmt.Sprintf(
// 				"http://localhost:%d/tokens",
// 				model.Conf.Server.Port,
// 			),
// 			nil,
// 		)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		req.SetBasicAuth("dima", "password")
// 		resp, err := c.Do(req)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		const SET_COOKIE = "Set-Cookie"
// 		re := regexp.MustCompile(`^refresh_token=[^;]+`)
// 		rt := re.FindString(resp.Header.Get(SET_COOKIE))
// 		if len(rt) == 0 {
// 			t.Errorf("\n\nerror: refresh token are empty\n\n")
// 		}

// 		// Запрос
// 		req, err = http.NewRequest(
// 			"GET",
// 			fmt.Sprintf(
// 				"http://localhost:%d/refresh",
// 				model.Conf.Server.Port,
// 			),
// 			nil,
// 		)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		req.Header.Set("X-Forwarded-For", "192.168.1.1")

// 		const COOKIE = "Cookie"
// 		req.Header.Set(COOKIE, rt)

// 		resp, err = c.Do(req)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		// Проверка статуса
// 		if http.StatusOK != resp.StatusCode {
// 			t.Errorf(
// 				"error: status code not equal\n expected: %d\n actual: %d\n",
// 				http.StatusOK, resp.StatusCode)
// 		}
// 		// Получение тела ответа
// 		buf := new(bytes.Buffer)
// 		_, err = buf.ReadFrom(resp.Body)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		if resp.StatusCode == http.StatusOK {
// 			var tks entities.Tokens
// 			if err := json.Unmarshal(buf.Bytes(), &tks); err != nil {
// 				log.Fatal(err)
// 			}
// 			// Проверка длинны токена доступа.
// 			if len(tks.AccessToken) == 0 {
// 				t.Errorf("\n\nerror: access token are empty\n\n")
// 			}
// 		}
// 	})
// }
