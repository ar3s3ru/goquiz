package auth

// basicAuthenticator definisce una struttura che mantiene il keyring degli utenti registrati
// e una coda di richieste da gestire.
type basicAuthenticator struct {
	users    map[string]string
	requests chan func(map[string]string)
}

// authenticator è l'agente di autenticazione che mantiene il keyring del processo.
// Esiste un solo keyring per processo, quindi viene usato come variabile globale
// con scope limitato a questa translation unit.
var authenticator = &basicAuthenticator{
	users:    make(map[string]string),
	requests: make(chan func(map[string]string)),
}

func init() {
	// Esegue una funzione asincrona che gestisce le richieste arrivate sull'authenticator
	go func(auth *basicAuthenticator) {
		for request := range auth.requests {
			request(auth.users)
		}
	}(authenticator)
}

// handle permette di aggiungere richieste da gestire alla coda dell'authenticator.
func (b *basicAuthenticator) handle(fn func(users map[string]string)) {
	if fn != nil {
		b.requests <- fn
	}
}
