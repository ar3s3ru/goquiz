package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

func main() {
	// Istanziamo un nuovo router: si occuperà principalmente di eseguire le funzioni
	// registrate sugli endpoint specificati.
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		// La logica definita in "handler" viene eseguita sull'endpoint http://localhost:8080/hello
		r.Get("/hello", handler)
	})

	// Fai partire il server sulla porta 8080
	log.Fatal(http.ListenAndServe(":8080", r))
}

func handler(rw http.ResponseWriter, r *http.Request) {
	// Chiudi il body della Request all'uscita
	defer r.Body.Close()

	// Semplice hello world!
	w := "Hello world!"
	// Possiamo usare diversi campi della richiesta HTTP, per esempio...
	w = fmt.Sprintf("%s\nfrom %s\n", w, r.RemoteAddr)

	// Settiamo come status "200 OK"
	rw.WriteHeader(http.StatusOK)
	// Scriviamo la risposta
	rw.Write([]byte(w))
}
