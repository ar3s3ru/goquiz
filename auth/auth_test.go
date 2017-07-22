package auth_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/ar3s3ru/goquiz/auth"
	"github.com/go-chi/chi"
	"github.com/h2non/baloo"
)

var client *baloo.Client

func TestMain(m *testing.M) {
	// Imposta il server
	r := chi.NewRouter()
	r.Use(auth.BasicAuthMiddleware)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		id, _ := auth.UserID(r)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{"user":"%s"}`, id)))
	})
	// Fai partire il server
	go func(r chi.Router) {
		http.ListenAndServe(":8080", r)
	}(r)
	// Attendi 1 secondo per l'avvio del server
	<-time.Tick(time.Second)
	// Crea il client baloo
	client = baloo.New("http://localhost:8080/")
	// Fai partire i test
	m.Run()
}

func TestNoAuthorizationHeader(t *testing.T) {
	client.Get("/").Expect(t).
		JSON([]byte(`{"message":"no Authorization header provided"}`)).
		Status(http.StatusUnauthorized).
		Done()
}

func TestNewUser(t *testing.T) {
	client.Get("/").
		SetHeader("Authorization", "Basic ZGFuaWxvOlg="). // danilo:X
		Expect(t).
		JSON([]byte(`{"user":"danilo"}`)).
		Status(http.StatusOK).
		Done()
}

func TestExistingUserWrongPassword(t *testing.T) {
	client.Get("/").
		SetHeader("Authorization", "Basic ZGFuaWxvOlhY"). // danilo:x
		Expect(t).
		JSON([]byte(`{"message":"wrong password for user"}`)).
		Status(http.StatusUnauthorized).
		Done()
}
func TestExistingUserGoodPassword(t *testing.T) {
	client.Get("/").
		SetHeader("Authorization", "Basic ZGFuaWxvOlg="). // danilo:X
		Expect(t).
		JSON([]byte(`{"user":"danilo"}`)).
		Status(http.StatusOK).
		Done()
}

func TestEmptyPassword(t *testing.T) {
	client.Get("/").
		SetHeader("Authorization", "Basic ZGFuaWxvOg=="). // danilo:X
		Expect(t).
		JSON([]byte(`{"message":"empty passwords are not allowed!"}`)).
		Status(http.StatusBadRequest).
		Done()
}
