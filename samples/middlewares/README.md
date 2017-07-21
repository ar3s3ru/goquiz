# *Working sample*: Middlewares

Se siete arrivati qui dal sample sui `context.Context`, avrete visto delle funzioni abbastanza strane all'interno del `main.go` di quel sample.

Cosa sono i *__middleware__*?

Dall'articolo [*Making and Using HTTP Middleware*](http://www.alexedwards.net/blog/making-and-using-middleware):  
> [...] self-contained code which independently acts on a request before or after your normal application handlers.

In altre parole, sono *__decoratori__* che vengono applicati sugli `http.Handler` che vanno a definire la logica di un endpoint.

## Yes, but... how?

All'interno dell'articolo citato in precedenza ci sono diversi esempi di middlewares.  

Noi ne vedremo uno abbastanza semplice nel nostro `main.go`: *loggare una richiesta HTTP ricevuta*!

È molto utile farsi stampare le varie richieste HTTP ricevute sulla console dove il processo è stato avviato.

```bash
$ ./myServer
[dd-MM-AAAA hh:mm] Server started!
[dd-MM-AAAA hh:mm] 192.168.1.127 | GET -- "/hello" | 59ms
```

Questo potrebbe essere una tipica linea di log su console, dove viene stampato:
* la data del log entry
    ```
    [dd-MM-AAAA hh:mm] ...
    ```
* il *Remote Address* del client
    ```
    [dd-MM-AAAA hh:mm] 192.168.1.127 ...
    ```
* il *metodo HTTP* usato e l'*endpoint* richiesto
    ```
    [dd-MM-AAAA hh:mm] 192.168.1.127 | GET -- "/hello" ...
    ```
* il tempo di servizio per processare la richiesta
    ```
    [dd-MM-AAAA hh:mm] 192.168.1.127 | GET -- "/hello" | 59ms
    ```

## Oltre gli esempi...

I middleware vengono tipicamente usati per aggiungere particolari valori sul `context.Context` di una richiesta HTTP.

```go
// addValueOnContext prende in input un http.Handler iniziale e ritorna
// un http.Handler con logica estesa (aggiunge un valore sulla richiesta HTTP 
// originale ed esegue l'handler `h`)
func addValueOnContext(h http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Esegui h con un context.Context diverso
        h(w, r.WithContext(r.Context(), contextKey("new-value"), 42))
    })
}
```

Riprendete il working sample sui `context.Context` per avere un insight maggiore sul funzionamento dei middleware + `context.Context`!