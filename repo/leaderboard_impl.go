package repo

import (
	"context"
	"errors"
	"net/http"
)

type leaderboardImpl struct {
	results  map[string]uint
	requests chan func(results map[string]uint)
}

var leaderboard = &leaderboardImpl{
	results:  make(map[string]uint),
	requests: make(chan func(results map[string]uint), 1), // Bufferizziamo per le performance
}

func init() {
	go func(lb *leaderboardImpl) {
		for fn := range lb.requests {
			fn(lb.results)
		}
	}(leaderboard)
}

func (lb *leaderboardImpl) handle(ok chan<- bool, fn func(map[string]uint)) {
	if fn != nil {
		lb.requests <- func(results map[string]uint) {
			fn(results)
			ok <- true
		}
	}
}

func (lb *leaderboardImpl) Results() (r LeaderboardResults, _ error) {
	ok := make(chan bool, 1)
	defer close(ok)
	lb.handle(ok, func(results map[string]uint) {
		for u, s := range results {
			r.Results = append(r.Results, LeaderboardResult{User: u, Score: s})
		}
	})
	<-ok
	return
}

func (lb *leaderboardImpl) AddScore(user string, score uint) error {
	ok := make(chan bool, 1)
	defer close(ok)
	lb.handle(ok, func(results map[string]uint) {
		if s, ok := results[user]; ok {
			// Incrementa lo score attuale
			score += s
		}
		results[user] = score
	})
	<-ok
	return nil
}

// Middlewares e Context helpers ----------------------------------------------------------------

var leaderboardKey = &contextKey{name: "leaderboardKey"}

func AddLeaderboardAccess(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r.WithContext(
			context.WithValue(
				r.Context(),
				leaderboardKey,
				leaderboard,
			),
		))
	})
}

func ResultsRepositoryFromContext(ctx context.Context) (LeaderboardResultsRepository, error) {
	v := ctx.Value(leaderboardKey)
	if v == nil {
		goto err
	}
	if lb, ok := v.(LeaderboardResultsRepository); ok {
		return lb, nil
	}
err:
	return nil, ErrNoLeaderboardRepositoryInContext
}

func AddScoreRepositoryFromContext(ctx context.Context) (LeaderboardAddScoreRepository, error) {
	v := ctx.Value(leaderboardKey)
	if v == nil {
		goto err
	}
	if lb, ok := v.(LeaderboardAddScoreRepository); ok {
		return lb, nil
	}
err:
	return nil, ErrNoLeaderboardRepositoryInContext
}

var ( // Gli errori sono definiti qui
	ErrNoLeaderboardRepositoryInContext = errors.New("no leaderboard repository in the context")
)
