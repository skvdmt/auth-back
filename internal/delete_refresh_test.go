//go:build integration

package internal

// import (
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"testing"

// 	"github.com/skvdmt/auth-back/internal/model"
// 	_ "github.com/skvdmt/auth-back/testing_init"
// )

// // TestDeleteRefresh Тестирование удаления токена обновления.
// func TestDeleteRefresh(t *testing.T) {
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

// 	t.Run("delete refresh token", func(t *testing.T) {
// 		c := http.Client{}
// 		// Получение токена обновления
// 		req, err := http.NewRequest(
// 			"DELETE",
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

// 		req.SetBasicAuth("dima", "password")
// 		resp, err := c.Do(req)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		const SET_COOKIE = "Set-Cookie"
// 		if len(resp.Header.Get(SET_COOKIE)) == 0 {
// 			t.Errorf("\n\nerror: set-cookie not exists\n\n")
// 		}
// 	})
// }
