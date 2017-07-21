package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

// Key per ricavare lo userId dal contesto
const userIDKey = "userId"

func main() {
	r := chi.NewRouter()

	// `userId` è un parametro che viene dall'URL della richiesta HTTP.
	// Una richiesta del tipo
	//
	//     http://localhost:8080/greatNickname
	//
	// assegna "greatNickname" ad `userId` nel `context.Context``
	// della `http.Request`.
	//
	// Vedi https://github.com/go-chi/chi#url-parameters
	//
	r.With(simplifyContextValueRetrieval).Get("/{userId}", namedHelloWorld)

	r.With(
		simplifyContextValueRetrieval,
		// Questo è un middleware, verrà illustrato nel working sample relativo.
		// In breve è un *decoratore*, ovvero una funzione che estende
		// il comportamento di un'altra funzione.
		//
		// La modifica del contesto tramite `modifyUserID` è fatta qui.
		func(h http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				h.ServeHTTP(w, r.WithContext(
					// Utiliza la stessa `http.Request` ma con `context.Context`
					// ricavato dalla chiamata a `modifyUserID`.
					modifyUserID(r.Context()),
				))
			})
		},
	).Get("/modify/{userId}", namedHelloWorld)

	// Avvia il server
	log.Fatal(http.ListenAndServe(":8080", r))
}

// getUserID ricava lo userId (definito come stringa) dal context.Context
// fornito.
//
// `ok` è una guardia booleana che indica se il contenuto di `id` è valido.
// In altre parole, stabilisce se la funzione è andata a buon fine.
//
// Usate la funzione così (è, tra l'altro, "idiomatic go")
//
//     if id, ok := getUserID(ctx); ok {
//	       // Usa `id` qui
//     }
//
// Oppure
//
//     id, ok := getUserID(ctx)
//     if !ok {
//         // Gestisci il caso sfavorevole...
//     }
//     // Gestisci il caso favorevole...
//
func getUserID(ctx context.Context) (id string, ok bool) {
	// Prendi lo userId dal contesto
	v := ctx.Value(userIDKey)
	if v == nil { // Se v è nil, allora nessuno userId è stato inserito
		return
	}
	// Fai type assertion per stabilire se il valore dello userId
	// trovato sia effettivamente un intero
	id, ok = v.(string)
	return
}

// modifyUserID prende il `context.Context`, cerca lo userId e lo modifica,
// restituendo il contesto con l'id modificato.
func modifyUserID(ctx context.Context) context.Context {
	id, ok := getUserID(ctx)
	if !ok {
		return ctx
	}
	// Modifica il valore dello userId aggiungendo "-MODIFIED" in coda
	return context.WithValue(ctx, userIDKey, fmt.Sprintf("%s-MODIFIED", id))
}

func namedHelloWorld(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// Prendi lo userId dal contesto della richiesta
	id, ok := getUserID(r.Context())
	if !ok {
		// Nessun id è stato specificato, very bad
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Printa un "Hello World" usando l'id ricavato
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hello World, %s!", id)))
}

// Questa funzione è un middleware, serve per semplificare l'utilizzo di `chi`
// in modo da farvi capire come funziona il context :-)
func simplifyContextValueRetrieval(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, userIDKey)
		h.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), userIDKey, id)))
	})
}
