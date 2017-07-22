# *Working samples*: `context.Context`

Il `context.Context` serve a definire un *contesto* vero e proprio che può essere passato in giro per il proprio codice.

In parole povere, `context.Context` è una *mappa key-value estensibile*.

La documentazione riguardante `context.Context` si trova [qui](https://godoc.org/context)!

Un articolo interessante riguardante le keys di `context.Context` si trova [qui](https://medium.com/@matryer/context-keys-in-go-5312346a868d)!

## `context.Context` nelle HTTP requests

`context.Context` trova particolare utilizzo all'interno degli `http.Handler`.

Difatti, `http.Request` possiede un metodo `.Context()` e `.WithContext()` che permettono, rispettivamente, di ottenere il *__contesto relativo alla richiesta HTTP__* e *__modificare la richiesta HTTP con un nuovo contesto__*.

Poichè é normale prassi utilizzare una serie di *__middleware__* prima dell'effettivo `http.Handler` che contiene la logica dell'endpoint, il `context.Context` della richiesta HTTP viene tipicamente *arricchito* con valori computati nei vari *middleware*, in modo da rendere l'`http.Handler` finale indipendente da oggetti esterni (variabili globali, ecc...) -- è, in un certo senso, una forma di *__dependency injection__*.

Il `main.go` di questo sample contiene un esempio abbastanza esplicativo che illustra l'utilizzo di `context.Context` all'interno di un `http.Handler`.

La cosa avrà più senso dopo il sample sui *middleware*, in particolare dopo il sample sull'*authentication* e *data-access-layer*!