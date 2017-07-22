package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"sync/atomic"

	"github.com/go-chi/chi"
)

// In questo sample creiamo un semplice server che esegue una semplice task:
// mantenere una lista di persone con nome, cognome e data di nascita!
//
// Le persone possono essere aggiunte tramite un particolare endpoint.
// La lista delle persone viene poi restituita dal server su una chiamata ad
// un endpoint diverso.

// Person indica una persona registrata su questo server.
// Ogni volta che una nuova persona viene aggiunta, gli viene assegnato un ID univoco.
type Person struct {
	ID        uint64 `json:"id"` // Questo field è accettato solo in output
	Name      string `json:"name"`
	Surname   string `json:"surname"`
	BirthDate string `json:"birthDate"`
}

// registrationsCounter è un intero che viene incrementato ad ogni registrazione.
// Viene usato per assegnare degli ID alle nuove registrazioni.
//
// P.S. è molto simile al sample su `context.Context`!
var registrationsCounter uint64

var (
	// people sarà la lista di utenti registrati sul server.
	people []Person

	// peopleMutex serve a garantire la consistenza della lista people.
	// Le chiamate al server possono essere concorrenti!
	peopleMutex sync.RWMutex
)

func main() {
	r := chi.NewRouter()
	r.Post("/new", registerNewPerson)
	r.Get("/list", listPeople)

	log.Fatal(http.ListenAndServe(":8080", r))
}

// handleRegistrationError stampa "500 Internal Server Error" sulla risposta HTTP
// e logga l'errore contenuto in `err`.
func handleRegistrationError(w http.ResponseWriter, err error) {
	log.Printf("Error registering new Person - Reading Body: %s\n", err)
	w.WriteHeader(http.StatusInternalServerError)
}

func registerNewPerson(w http.ResponseWriter, r *http.Request) {
	// Chiudiamo il body della richiesta alla fine dell'handler
	defer r.Body.Close()

	// Leggiamo il body e tentiamo di deserializzarlo nel tipo Person
	raw, err := ioutil.ReadAll(r.Body)
	if err != nil {
		// Argh! C'è stato qualche errore in lettura!
		// Ritorniamo un Internal Server Error e logghiamo sulla console l'errore
		handleRegistrationError(w, err)
		return
	}

	// Tentiamo di deserializzare il body
	var person Person
	if err := json.Unmarshal(raw, &person); err != nil {
		// Un altro errore...
		// Anche questa volta, logghiamo l'errore e ritorniamo...
		handleRegistrationError(w, err)
		return
	}

	// Assegniamo a questo nuovo Person un ID
	person.ID = atomic.AddUint64(&registrationsCounter, 1)

	// Ottimo! Abbiamo deserializzato il body, e non abbiamo incontrato errori!
	// Tralasciamo il check sul "birthDate" (che sia effettivamente una data) e aggiungiamo
	// `person` alla lista di persone sul server.
	peopleMutex.Lock()
	// L'aggiunta viene eseguita in una regione critica, in quanto `people` può essere acceduto
	// in maniera concorrente!
	people = append(people, person)
	peopleMutex.Unlock()

	// Abbiamo fatto!
	// Segnaliamo che la risorsa è stata creata e usciamo.
	w.WriteHeader(http.StatusCreated)
}

func listPeople(w http.ResponseWriter, r *http.Request) {
	// Acquisiamo una read lock sulla mutex della lista `people`
	peopleMutex.RLock()
	// La read lock viene rilasciata all'uscita della funzione!
	defer peopleMutex.RUnlock()

	// Proviamo a fare marshalling della lista...
	raw, err := json.Marshal(people)
	if err != nil {
		// Errore durante il marshalling: come prima, logghiamo l'errore e
		// usciamo con "500 Internal Server Error"
		log.Printf("Error while marshalling People list: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// È andato tutto a buon fine! Printiamo il JSON e status "200 OK"
	w.WriteHeader(http.StatusOK)
	w.Write(raw)
}
