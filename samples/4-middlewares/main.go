package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
)

func main() {
	r := chi.NewRouter()

	// `Use()` ci permette di specificare dei middleware che devono essere
	// eseguiti per tutti gli handler registrati con questo router.
	// È un modo conveniente per evitare di ripetere
	//
	//     r.With(middlewares...).Get(...)
	//
	// per ogni richiesta che utilizza un particolare middleware.
	// In questo caso, vuol dire che `logRequest` viene eseguito per ogni richiesta
	// che arriva a questo router (quindi `helloWorld` e `greeter`).
	r.Use(logRequest)

	r.Get("/hello", helloWorld)
	r.Get("/greet", greeter)

	log.Fatal(http.ListenAndServe(":8080", r))
}

const (
	// logTimeFormat specifica il layout della data che vogliamo utilizzare.
	// In Go, a differenza di Java ed altri linguaggi che utilizzano
	// layout del tipo "AAAA-MM-dd hh:mm:ss", si utilizza una data vera e propria
	// come layout di partenza, che è
	//
	//     Mon Jan 2 15:04:05 MST 2006        (dalla godoc del package `time`)
	//
	logTimeFormat = "2006/01/02 15:04:05"
)

// logRequest è il middleware che viene usato per loggare i dettagli della richiesta
// sulla console.
func logRequest(h http.Handler) http.Handler {
	// Istanziamo un nuovo logger per stampare i dettagli della richiesta.
	//
	// N.B. `log.New()` viene chiamato una sola volta, non tutte le volte
	//      che una richiesta con middleware `logRequest` viene eseguita!
	//
	// Infatti, è vero che `logger` è una variabile locale all'ambiente di definizione (questa funzione),
	// però `logRequest` viene effettivamente eseguito *una sola volta* (quando si usa `r.Use()` del router di `chi`):
	// la funzione che viene richiamata tutte le volte che arriva una richiesta è la closure che `logRequest` ritorna!
	//
	// Se vi ho confuso le idee, chiamatemi per chiarimenti :-)
	logger := log.New(os.Stdout, "", 0)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Riferimento al time prima di eseguire l'handler vero e proprio
		start := time.Now()
		// Esegue la logica dell'endpoint chiamando l'handler in input
		// N.B. l'handler in input sarà uno tra le funzioni `helloWorld` e `greeter`!
		h.ServeHTTP(w, r)
		// Riferimento alla fine della gestione della richiesta vera e propria
		end := time.Now()
		// Decommentando la riga qui sotto potete vedere effettivamente che `logger` viene istanziato una sola volta,
		// poichè il valore del puntatore è sempre lo stesso!
		// fmt.Printf("[!] Logger pointer is %p\n", logger)

		// Stampa vera e propria sulla console (a.k.a. stdout)
		logger.Printf(
			"[%s] %s -- %s  %s -- %s\n",         // Questo è il formato che la stampa di log dovrebbe avere
			start.Format("2006/01/02 15:04:05"), // Timestamp di quando la richiesta è "arrivata"
			r.RemoteAddr,                        // Indirizzo IPv(4/6) del client
			r.Method,                            // Metodo usato per la richiesta
			r.RequestURI,                        // URI della richiesta
			end.Sub(start),                      // Tempo di servizio della richiesta
		)
	})
}

func helloWorld(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello world!"))
}

func greeter(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hello, %s!", r.RemoteAddr)))
}
