package repo

import (
	"errors"

	"math/rand"
	"time"

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

var questions = &questionsImpl{
	quizzes:  make(map[string]quiz),
	requests: make(chan func([]question, map[string]quiz), 1),
}

func init() {
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
