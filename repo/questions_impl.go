package repo

type question struct {
	Question string   `json:"question"`
	Answers  []string `json:"answers"`
	Score    uint     `json:"score"`
	Correct  uint     `json:"correct"`
}

type questions struct {
	questions []question
	request   chan func(s []question)
}

func init() {

}
