package repo

type LeaderboardResult struct {
	User  string `json:"user"`
	Score uint   `json:"score"`
}

type LeaderboardResults struct {
	Results []LeaderboardResult `json:"results"`
}

type LeaderboardResultsRepository interface {
	// Results ritorna tutti i risultati della leaderboard del server.
	// Ritorna un errore se il data provider fallisce in qualche modo.
	Results() (LeaderboardResults, error)
}

type LeaderboardAddScoreRepository interface {
	// AddScore aggiorna la leaderboard aggiungendo uno score per l'utente specificato.
	// Ritorna un errore se il data provider fallisce in qualche modo.
	AddScore(user string, score uint) error
}
