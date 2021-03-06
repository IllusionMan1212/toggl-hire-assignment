// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type Question interface {
	IsQuestion()
}

type AnswerInput struct {
	QuestionID           string   `json:"questionId"`
	QuestionType         string   `json:"questionType"`
	QuestionWeight       float64  `json:"questionWeight"`
	EnteredText          *string  `json:"enteredText"`
	SelectedOptionID     *string  `json:"selectedOptionId"`
	SelectedOptionWeight *float64 `json:"selectedOptionWeight"`
}

type ChoiceQuestion struct {
	ID      string    `json:"id"`
	Body    string    `json:"body"`
	Weight  float64   `json:"weight"`
	Options []*Option `json:"options"`
}

func (ChoiceQuestion) IsQuestion() {}

type Message struct {
	Content string `json:"content"`
}

type Option struct {
	ID     string  `json:"id"`
	Body   string  `json:"body"`
	Weight float64 `json:"weight"`
}

type Result struct {
	TotalQuestions int     `json:"totalQuestions"`
	CorrectAnswers int     `json:"correctAnswers"`
	Score          float64 `json:"score"`
}

type TextQuestion struct {
	ID     string  `json:"id"`
	Body   string  `json:"body"`
	Weight float64 `json:"weight"`
}

func (TextQuestion) IsQuestion() {}
