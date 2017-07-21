# goquiz

Quiz a risposta multipla online, in cui le domande sono a tema Go :-)

## Overview

Il topic di questo codelab è di sviluppare assieme il server del gioco.  
All'interno di questo repository troverete:
- il file `.json` delle domande del quiz
- dei *working samples* che vi spiegheranno come implementare alcune feature del server
- implementazione di alcune *componenti* del server opzionali (i.e. *__bonus points__*)

## First steps

Presupponendo che il vostro [`GOPATH` sia settato](https://golang.org/doc/code.html#GOPATH) correttamente, dobbiamo creare ora il *package* che ospiterà il vostro server.

Usando il vostro terminale, eseguite il comando  
```bash
$ mkdir -p $GOPATH/src/goquiz && cd $GOPATH/src/goquiz
```

Avete creato il package `goquiz` all'interno del vostro `GOPATH`!  
Ora possiamo iniziare a scrivere codice, ma...

...parliamo della *gestione delle dipendenze*.

### Dependency management

Per installare nuovi package nel vosto `GOPATH` locale, affinchè altri package possano usarli, è tramite il comando `go get`.  
Un esempio:  
```bash
$ go get -u github.com/go-chi/chi
```

Questo comando eseguirà una `git clone` automatica del package `github.com/go-chi/chi`, compilerà il package e renderà disponibile
i file compilati all'interno di `$GOPATH/pkg/<platform>/github.com/go-chi/chi`.

Nel corso del codelab faremo uso di diverse librerie esterne (tra cui proprio `chi`), dunque avrete bisogno di eseguire
`go get` più volte.

In alternativa, potete usare [`glide`](https://github.com/Masterminds/glide).

`glide` risolve automaticamente gli import necessari a partire da un file `glide.yaml`.  
Il `glide.yaml` necessario al codice è fornito all'interno di questo repository.

Vi basterà copiare `glide.yaml`, `glide.lock` ed eseguire il comando  
```bash
$ glide install
```

Nel caso in cui voleste utilizzare una libreria esterna non contenuta nel `glide.yaml` all'interno del vostro codice,
vi basterà eseguire il comando  
```bash
$ glide get <nome libreria>
```
in maniera analoga al comando `go get`.

## Specifications

Le specifiche del server sono date dall'*__API REST__ che espone*.  

Nel nostro caso, il server presenterà *__3 endpoints__*:

| **Metodo HTTP** | **Endpoint**      | **Dettagli**                                     |
|-----------------|-------------------|--------------------------------------------------|
| GET             | `/new-quiz`       | Richiedi un nuovo quiz al server                 |
| POST            | `/submit/:quizId` | Invia le risposte del quiz al server             |
| GET             | `/leaderboard`    | Richiedi la lista degli score migliori al server |

Il *JSON schema* delle risposte e delle richieste del/al server è il seguente

### `/new-quiz`

Dopo la richiesta `GET` al server su questo endpoint, il server ci manderà un JSON così formato:  
```json
{
    "id": string,
    "questions": [
        {
            "question": string,
            "answers": [
                string,
                string,
                string,
                ...
            ]
        }
    ]
}
```

### `/submit/:quizId`

Dopo una `POST` al server su questo endpoint
- se `quizId` è un ID invalido, il server risponde con `400 Bad Request`
- se `quizId` è un ID valido, ma non assegnato alla propria utenza, il server risponde con `403 Forbidden`
- se `quizId` è un ID valido ed assegnato alla propria utenza, il server risponde con `200 OK`

Il JSON che il server si aspetta è:  
```json
{
    "answers": [
        0,  // Indice della risposta alla domanda 1
        1,  // Indice della risposta alla domanda 2
        2,
        -1, // Se una domanda è stata skippata, usare -1
        ...
    ]
}
```

Il server invece risponderà così:  
```json
{
    "score": <valore intero dello score totale>,
    "results": [
        {
            "given": <intero della risposta data>,
            "correct": <intero della risposta corretta>,
        },
        {
            ...
        }
    ]
}
```

### `/leaderboard`

Il server risponderà sempre `200 OK` su questo endpoint, usando il JSON:  
```json
{
    "scores": [
        {
            "user": string,
            "score": int
        },
        {
            ...
        }
    ]
}
```

---

Created by @ar3s3ru and @CDimonaco