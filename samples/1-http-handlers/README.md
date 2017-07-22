# *Working samples*: HTTP Handlers

Quando una richiesta HTTP su un particolare endpoint viene eseguita, la logica che il server deve eseguire è definita
da un'*entità* chiamata *__HTTP handler__*.

In Go, questo *handler* altro non è se non un'*__interfaccia__*.

```go
// From `net/http` package
type Handler interface {
    // ...
    ServeHTTP(http.ResponseWriter, *http.Request)
}
```

Tutti i tipi di dato che soddisfano quest'interfaccia possono essere usati come *HTTP handlers*.

Nel caso medio, piuttosto che quest'interfaccia viene usata una *__funzione__* con una *signature particolare*
(ricordiamo che le funzioni in Go sono *first-class citizens*, i.e. tipi primitivi). 

```go
type HandlerFunc func(http.ResponseWriter, *http.Request)
```

Notate che la segnatura di `http.HandlerFunc` è la stessa del metodo `ServeHTTP` di `http.Handler`?

Infatti, `http.HandlerFunc` implementa `http.Handler` in questo modo:  
```go
func (fn HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    fn(w, r)
}
```

Semplice, no?

## `chi` e gli `http.Handler`

`chi` permette di eseguire il *__routing__* delle richieste HTTP (ne parleremo in un altro working sample).  

I metodi che utilizza `chi` impiegano tipicamente sia `http.HandlerFunc` (per le chiamate *GET*, *POST*, *DELETE*, ...)
che `http.Handler` (utile in casi più complessi, vedi il working sample sui *middlewares*).

Come abbiamo detto prima, tipicamente utilizziamo una funzione di tipo `http.HandlerFunc`.

```go
r := chi.NewRouter()
r.Get("/", func(w http.ResponseWriter, r *http.Request) {
    // Funzione che il server eseguirà ad una chiamata HTTP all'indirizzo http://<indirizzo-server>:<porta-server>/
})
```

Per avere un'idea dell'utilizzo, consultare il `main.go` contenuto in questa cartella!