# *Working sample*: JSON (un)marshalling

La realtà dei fatti è che, tipicamente, un server non risponderà ad un client usando una semplice stringa come *"Hello world!"* o *"Hi, visitor #123!"*.

Server e client comunicano, *tipicamente*, usando *__dati strutturati__*.  
Questi dati strutturati sono definiti dall'*__API__ di un server*.

Cosa c'entra allora *__JSON__*?

> In informatica, nell'ambito della programmazione web, JSON, acronimo di JavaScript Object Notation, è un formato adatto all'interscambio di dati fra applicazioni client-server.  
> (Da Wikipedia)

JSON può essere considerato, al giorno d'oggi, lo *standard* per lo scambio dei messaggi in applicazioni client-server.

Quindi, vediamo come utilizzarlo in Go! :-)

## `struct`ured data

Il modo più semplice per utilizzare JSON in Go è *definire un tipo concreto* (a.k.a. una `struct`) e utilizzare un
particolare costrutto di Go chiamato *__tag__*.

```go
type Person struct {
    Name    string `json:"name"`    // Questo è un tag!
    Surname string `json:"surname"` // Questo è un altro tag!
}
```

Un *tag* serve per dare informazioni sui *field* di una particolare *struttura concreta*.  
Questi tag vengono utilizzati dalla libreria che gestirà l'encoding e il decoding in JSON.

In questo caso, quei tag stanno ad indicare il *nome JSON per i field di quella struct*.

```json
{
    "name": "Mario",
    "surname": "Rossi"
}
```

si traduce, grazie ai tag, in

```go
person := Person{Name: "Mario", Surname: "Rossi"}
```

Semplice, no?

## JSON (Un)Marshalling

> In computer science, marshalling or marshaling is the process of transforming the memory representation of an object to a data format suitable for storage or transmission, and it is typically used when data must be moved between different parts of a computer program or from one program to another.  
> (Da Wikpedia)

*__Marshalling__* e *__Unmarshalling__* equivale a *serializzazione* e *deserializzazione* di un oggetto.  
In questo caso, parliamo di *(de)serializzazione in JSON*.

In Go, questo è possibile grazie al package `encoding/json` della stdlib

```go
import "encoding/json"

// ...

func main() {
    person := Person{Name: "Mario", Surname: "Rossi"}
    // `Marshal` prende un oggetto qualunque e lo serializza in JSON.
    // Attenzione! Può tornare un errore, e l'errore va gestito a modino :)
    // personJSON è un []byte, può essere scritto direttamente nel body di una risposta HTTP!
    personJSON, err := json.Marshal(person)

    // ...

    // Consideriamo il caso di aver ricevuto in input un JSON da un client,
    // e questo JSON viene rappresentato come []byte (raw data, per intenderci).
    // Per deserializzare il JSON in input, facciamo così...

    // newPerson conterrà il risultato della deserializzazione
    var newPerson Person

    // `Unmarshal` prende in input il JSON in formato raw e il **puntatore**
    // al risultato della deserializzazione, in modo che possa modificarne il contenuto
    // durante il decoding.
    // Attenzione! Può tornare un errore, e come sempre, va gestito!
    err := json.Unmarshal(rawData, &newPerson)

    // ...

    // Stamperà
    //     Person{Name: "Mario", Surname: "Rossi"}
    // se in `rawData` abbiamo messo il valore ricevuto prima in `personJSON` :-)
    fmt.Printf("%+v\n", newPerson)
}
```

### NOTA BENE!

* *Marshalling* e *unmarshalling* funzionano solo sui field *__public__* di una `struct`
    - field *public* (esportato) se la *__prima lettera è maiuscola__*  
        ```go
        type myType struct {
            Name string `json:"name"` // Questo field verrà (de)serializzato correttamente!
        }
        ```
    - field *private* (non-esportato) se la *__prima lettera è minuscola, o non maiuscola__*
        ```go
        type myType struct {
            name string `json:"name"` // Questo field NON verrà (de)serializzato!
        }
        ```
* __NON__ utilizzare interfacce come tipi dei field di una struct
    - il deserializzatore non ha conoscenza del tipo concreto dell'interfaccia, quindi produrrà un errore  
        ```go
        type myStruct struct {
            Data myInterface `json:"data"` // `myInterface` è un'interfaccia, può essere serializzato ma non deserializzato!
        }
        ``` 
    - *(P.S. in verità si può utilizzare, ma bisogna fare dei giochini strani quindi evitate nel corso del codelab)*
        - *(P.P.S. se siete __davvero__ curiosi di sapere come risolvere questa situazione, chiedete pure e vi sarà rivelato :-) )*

## Ok, now show me the __code__!

Il `main.go` di questo sample contiene un esempio di comunicazione client-server tipica, utilizzando una chiamata __POST__.

Per fare chiamate al server di questo sample, non è più possibile usare il browser... :-(

Usate `curl` o qualunque altro tool per fare chiamate HTTP con un body JSON (noi di GeoUniq siamo fan di [Postman](https://www.getpostman.com/apps)!)