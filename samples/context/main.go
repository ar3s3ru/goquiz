package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"

	"github.com/go-chi/chi"
)

// contextKey identifica la stringa della chiave da usare nel contesto.
// È pratica comune ridefinire un tipo *unexported* per le chiavi
// in modo da non rilevare informazioni riguardante valori contenuti nel contesto.
//
// È, dunque, una pratica di sicurezza (ed una best-practice, il linter darà un warning :-) )
type contextKey string

const counterKey = contextKey("requestCounter")

// counter è il contatore atomico che utilizzeremo per sapere il numero di richieste effettuate
var counter uint64

func main() {
	r := chi.NewRouter()
	// Questa funzione è un middleware; qui andiamo a modificare il contesto per tutte le chiamate
	// HTTP che arrivano. Consultate il working sample sui middleware per saperne di più!
	r.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// setRequestNumber modifica il contesto della richiesta in arrivo,
			// e il middleware esegue l'Handler vero e proprio con il contesto modificato.
			h.ServeHTTP(w, r.WithContext(setRequestNumber(r.Context())))
		})
	})

	// `userId` è un parametro che viene dall'URL della richiesta HTTP.
	// Una richiesta del tipo
	//
	//     http://localhost:8080/greatNickname
	//
	// assegna "greatNickname" ad `userId` nel `context.Context``
	// della `http.Request`.
	// Il contenuto del parametro URL può essere ricavato da
	//
	//     chi.URLParam(r *http.Request, key string)
	//
	// Vedi https://github.com/go-chi/chi#url-parameters
	r.Get("/{userId}", namedHelloWorld)

	// Avvia il server
	log.Fatal(http.ListenAndServe(":8080", r))
}

// getRequestNumber ricava il numero della richiesta (definito come uint64)
// dal context.Context fornito.
//
// `ok` è una guardia booleana che indica se il contenuto di `id` è valido.
// In altre parole, stabilisce se la funzione è andata a buon fine.
//
// Usate la funzione così (è, tra l'altro, "idiomatic go")
//
//     if id, ok := getRequestNumber(ctx); ok {
//	       // Usa `id` qui
//     }
//
// Oppure
//
//     id, ok := getRequestNumber(ctx)
//     if !ok {
//         // Gestisci il caso sfavorevole...
//     }
//     // Gestisci il caso favorevole...
//
func getRequestNumber(ctx context.Context) (id uint64, ok bool) {
	// Prendi lo userId dal contesto
	v := ctx.Value(counterKey)
	if v == nil {
		// Se v è nil, allora nessun counter è stato aggiunto al contesto,
		// il che vuol dire che avete sbagliato qualcosa in fase di routing :)
		return
	}
	// Fai type assertion per stabilire se il valore dell'id
	// trovato sia effettivamente un uin64
	id, ok = v.(uint64)
	return
}

// setRequestNumber prende il `context.Context` e inserisce il numero
// della richiesta attuale al suo interno.
func setRequestNumber(ctx context.Context) context.Context {
	// `WithValue` ritorna un contesto "arricchito" da un nuovo valore.
	// In questo caso, poichè le richieste possono essere concorrenti,
	// utilizziamo `atomic` per assicurare la consistenza del counter generale.
	return context.WithValue(ctx, counterKey, atomic.AddUint64(&counter, 1))
}

func namedHelloWorld(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// Prendi lo userId dal contesto della richiesta
	id, ok := getRequestNumber(r.Context())
	if !ok {
		// Nessun id è stato specificato, very bad
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Printa un "Hello World" usando l'id ricavato
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(
		"Hello World, %s!\nRequest number: %d",
		chi.URLParam(r, "userId"), // Prendi il parametro URL relativo a `userId`
		id, // ID della richiesta trovato in precedenza
	)))
}
