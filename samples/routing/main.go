package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

func main() {
	// Istanzia un nuovo router
	r := chi.NewRouter()
	// Routing molto semplice, nell'endpoint "/" esegui una Hello World col metodo GET
	r.Get("/", helloWorld)
	// Routing più complesso: la funzione `Route()` crea un "sub-router" che verrà utilizzato
	// se il pattern "/sub" viene matchato; tutte le altre funzioni di routing all'interno di `Route()`
	// sono composte con il pattern iniziale
	//
	// e.g. "/sub/hello", "/sub/ip"
	r.Route("/sub", func(r chi.Router) {
		// `r` parametro attuale della funzione è il sub-router.
		// Chiamarlo `r` "oscura" (i.e. shadowing), ovviamente, il router iniziale.
		r.Get("/hello", helloWorldSub)
		r.Get("/ip", getIP)
	})

	// Avvia il server
	log.Fatal(http.ListenAndServe(":8080", r))
}

func helloWorld(w http.ResponseWriter, r *http.Request) {
	// Chiudi il body all'uscita della funzione
	defer r.Body.Close()
	// Scrivi prima lo status code
	w.WriteHeader(http.StatusOK)
	// Scrivi la risposta
	w.Write([]byte("Hello world!"))
}

func helloWorldSub(w http.ResponseWriter, r *http.Request) {
	// Chiudi il body all'uscita della funzione
	defer r.Body.Close()
	// Scrivi prima lo status code
	w.WriteHeader(http.StatusOK)
	// Scrivi la risposta
	w.Write([]byte("Hello SUB-world!"))
}

func getIP(w http.ResponseWriter, r *http.Request) {
	// Chiudi il body
	defer r.Body.Close()
	// Scrivi status code
	w.WriteHeader(http.StatusOK)
	// Scrivi risposta
	w.Write([]byte(fmt.Sprintf("Your ip is %s! Hello!", r.RemoteAddr)))
}
