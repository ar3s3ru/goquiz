package repo

import (
	"errors"

	"math/rand"
	"time"

	"sync"

	"net/http"

	"os"

	"encoding/json"
	"io/ioutil"

	"context"

	"github.com/satori/go.uuid"
)

type question struct {
	Question string   `json:"question"`
	Answers  []string `json:"answers"`
	Score    uint     `json:"score"`
	Correct  uint     `json:"correct"`
}

type quiz struct {
	user      string
	questions []int
}

type questionsImpl struct {
	questions []question
	quizzes   map[string]quiz // Non so qual è il plurale di "quiz", lol
	requests  chan func([]question, map[string]quiz)
}

var questions *questionsImpl
var once sync.Once

func (q *questionsImpl) init() {
	go func(q *questionsImpl) {
		for request := range q.requests {
			request(q.questions, q.quizzes)
		}
	}(questions)
}

func (q *questionsImpl) handle(ok chan<- bool, fn func([]question, map[string]quiz)) {
	if fn != nil {
		q.requests <- func(s []question, qz map[string]quiz) {
			fn(s, qz)
			ok <- true
		}
	}
}

func (q *questionsImpl) New(user string) (quizz Quiz, err error) {
	ok := make(chan bool, 1)
	defer close(ok)
	q.handle(ok, func(s []question, qz map[string]quiz) {
		id := uuid.NewV4().String()
		qq := quiz{user: user, questions: randIndexes(6, uint(len(s)))}

		qz[id] = qq
		quizz.id = id
		for i := range qq.questions {
			qq := s[i]
			quizz.questions = append(quizz.questions, Question{
				question: qq.Question,
				correct:  qq.Correct,
				answers:  qq.Answers,
			})
		}
	})
	<-ok
	return
}

func (q *questionsImpl) Get(user, id string) (quiz Quiz, err error) {
	ok := make(chan bool, 1)
	defer close(ok)
	q.handle(ok, func(s []question, qz map[string]quiz) {
		v, ok := qz[id]
		if !ok {
			err = errors.New("no quiz found")
			return
		}
		if v.user != user {
			err = errors.New("user mismatch")
			return
		}
		for i := range v.questions {
			qq := s[i]
			quiz.questions = append(quiz.questions, Question{
				question: qq.Question,
				correct:  qq.Correct,
				answers:  qq.Answers,
			})
		}
		quiz.id = id
	})
	<-ok
	return
}

// Good ol' functions -------------------------------------------------------------------

func randIndexes(len, limit uint) (res []int) {
	res = make([]int, len, len)
	used := make(map[int]interface{}, len)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := uint(0); i < limit; i++ {
		v := r.Int() % int(limit)
		for _, ok := used[v]; !ok; v = r.Int() % int(limit) {
		}
		used[v] = nil
		res[i] = v
	}
	return
}

// Context stuff ------------------------------------------------------------------------

var questionsKey = &contextKey{name: "questions"}

func panicEventually(err error) {
	if err != nil {
		panic(err)
	}
}

func InmemQuestionsRepository(jsonFile string) func(http.Handler) http.Handler {
	// Inizializza per la prima volta il repository.
	once.Do(func() {
		file, err := os.Open(jsonFile)
		panicEventually(err)

		raw, err := ioutil.ReadAll(file)
		panicEventually(err)

		var slice []question
		panicEventually(json.Unmarshal(raw, &slice))

		questions = &questionsImpl{
			questions: slice,
			quizzes:   make(map[string]quiz),
			requests:  make(chan func([]question, map[string]quiz)),
		}
		questions.init()
	})
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r.WithContext(
				context.WithValue(r.Context(), questionsKey, questions),
			))
		})
	}
}

func NewQuizRepositoryFromContext(ctx context.Context) (NewQuizRepository, error) {
	v := ctx.Value(questionsKey)
	if v == nil {
		goto err
	}
	if q, ok := v.(NewQuizRepository); ok {
		return q, nil
	}
err:
	return nil, ErrNoLeaderboardRepositoryInContext
}

func GetQuizRepositoryFromContext(ctx context.Context) (GetQuizRepository, error) {
	v := ctx.Value(questionsKey)
	if v == nil {
		goto err
	}
	if q, ok := v.(GetQuizRepository); ok {
		return q, nil
	}
err:
	return nil, ErrNoLeaderboardRepositoryInContext
}

var (
	ErrNoQuizRepositoryInContext = errors.New("no quiz repository in the context")
)
