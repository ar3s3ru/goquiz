# *Working sample*: Routing

L'*__HTTP routing__* indica l'*indirizzamento della richiesta HTTP* verso un particolare *handler*, se è stato registrato.

Ciò vuol dire che, avendo registrato gli endpoint `/hello` e `/world`, due chiamate HTTP diverse a questi endpoint
produrranno l'esecuzione dei propri handler.

Nel `main.go` fornito in questo sample potete vedere in che modo il routing viene eseguito.

Ci sono routing più o meno complessi che possono essere eseguiti grazie a `chi`.  
I routing più semplici son quelli che richiedono un'unica chiamata col metodo che si intende utilizzare per quell'endpoint:  
```go
r := chi.NewRouter()
r.Get("/", myHandler)
r.Post("/", myHandlerPOST)
// ...
``` 

Routing più complessi possono essere eseguiti tramite la funzione `Route()` o tramite le funzioni `Mount()` e `Handle()`:
- `Route()` viene illustrata nel `main.go` di questo working sample
- `Mount/Handle()` vengono illustrati nel working sample dei *middleware*