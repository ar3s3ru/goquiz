package auth

import (
	"context"
	"net/http"
)

var (
	// unauthorizedErrorResponse viene ritornato dal middleware se l'header Authorization non è stato fornito
	unauthorizedErrorResponse = []byte(`{"message":"no Authorization header provided"}`)

	// emptyPasswordResponse viene ritornato dal middleware se è stata usata una password vuota
	emptyPasswordResponse = []byte(`{"message":"empty passwords are not allowed!"}`)

	// wrongPasswordResponse viene ritornato dal middleware se la password non è corretta
	wrongPasswordResponse = []byte(`{"message":"wrong password for user"}`)
)

type contextKey struct {
	key string
}

var myContextKey = &contextKey{"auth"}

type contextAuthentication struct {
	user, password string
}

func handleMiddlewareError(w http.ResponseWriter, code int, err []byte) {
	w.WriteHeader(code)
	w.Header().Add("Content-Type", "application/json")
	w.Write(err)
}

// BasicAuthMiddleware può essere usato per gestire l'autenticazione tramite `Authorization` header.
// Se nell'header viene specificato un utente non presente nel keyring, viene automaticamente aggiunto.
// Se viene specificato un utente già presente ma vengono usate credenziali diverse da quelle iniziali,
// viene ritornato un errore.
//
// Nel caso favorevole, vengono aggiunti i dettagli di autenticazione al contesto della richiesta HTTP
// che possono essere ricavati dalla funzione `UserID` di questo package.
func BasicAuthMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, p, ok := r.BasicAuth()
		if !ok {
			// Nessun header Authorization è stato fornito, errore!
			handleMiddlewareError(w, http.StatusUnauthorized, unauthorizedErrorResponse)
			return
		}
		if p == "" {
			// La password fornita è vuota, è una Bad Request!
			handleMiddlewareError(w, http.StatusBadRequest, emptyPasswordResponse)
			return
		}

		// Apriamo un nuovo blocco per le dichiarazioni locali di variabili
		// che non devono provocare errori con l'utilizzo di `goto` nello scope esterno
		{
			// Il pattern usato qui è di creare un canale di comunicazione tra una closure
			// e il chiamante della funzione.
			// Prima o poi, qualcuno inserirà un valore su quel canale, e noi continueremo l'esecuzione
			ch := make(chan string, 1)
			// N.B. i canali vanno chiusi! Altrimenti, leaks!! :-O
			defer close(ch)
			authenticator.handle(func(users map[string]string) {
				// Cerca l'utente u e restituisce la password
				if p, ok := users[u]; ok {
					ch <- p
					return
				}
				ch <- ""
			})

			v := <-ch
			if v == "" {
				// L'utente non è stato trovato nella mappa, aggiungilo!
				ok := make(chan bool, 1)
				defer close(ok)
				authenticator.handle(func(users map[string]string) {
					users[u] = p
					ok <- true
				})
				<-ok // Attendi che la funzione venga gestita
				goto handleOk
			}
			// L'utente è stato trovato, testiamo che le password abbiano match
			if p == v {
				goto handleOk
			}
			// Le password non matchano, errore!!
			handleMiddlewareError(w, http.StatusUnauthorized, wrongPasswordResponse)
		}
		return

	handleOk:
		// Esegui l'handler con un contesto che include i dettagli di autenticazione
		h.ServeHTTP(w, r.WithContext(
			context.WithValue(
				r.Context(), myContextKey, contextAuthentication{user: u, password: p},
			),
		))
		return
	})

}

// UserID è un helper che restituisce l'id dell'utente autorizzato e una guardia booleana
// che specifica se il risultato in `id` ha senso oppure no.
func UserID(r *http.Request) (id string, ok bool) {
	v := r.Context().Value(myContextKey)
	if v == nil {
		return
	}
	var ctx contextAuthentication
	if ctx, ok = v.(contextAuthentication); ok {
		id = ctx.user
	}
	return
}
