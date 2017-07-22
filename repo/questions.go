package repo

// Question rappresenta una domanda del quiz.
type Question struct {
	question string
	answers  []string
}

func (q Question) Question() string {
	return q.question
}

func (q Question) Answers() []string {
	return q.answers
}

// Quiz rappresenta un'insieme di domande. Ogni quiz è "univoco", poichè possiede
// un id univoco.
type Quiz struct {
	id        string
	questions []Question
}

func (q Quiz) ID() string {
	return q.id
}

func (q Quiz) Questions() []Question {
	return q.questions
}

// NewQuizRepository è un repository che permette di creare un nuovo Quiz per un certo utente.
type NewQuizRepository interface {
	New(user string) (Quiz, error)
}

// GetQuizRepository è un repository che permette di recuperare un Quiz precedentemente creato
// dal data provider.
type GetQuizRepository interface {
	Get(user, id string) (Quiz, error)
}

// // Result indica il risultato di una domanda.
// type Result struct {
// 	correct int8
// 	actual  int8
// }

// func (r Result) IsCorrect() bool {
// 	return r.correct == r.actual
// }

// func (r Result) Correct() int8 {
// 	return r.correct
// }

// func (r Result) Actual() int8 {
// 	return r.actual
// }

// // Results indica il risultato di un Quiz. Ogni Results ha uno score del Quiz
// // e la lista di tutti i risultati.
// type Results struct {
// 	score   uint
// 	results []Result
// }

// func (r Results) Score() uint {
// 	return r.score
// }

// func (r Results) Results() []Result {
// 	return r.results
// }

// // SubmitQuizRepository è il repository che permette di inviare
// type SubmitQuizRepository interface {
// 	Submit(id string, answers []int8) (Results, error)
// }
