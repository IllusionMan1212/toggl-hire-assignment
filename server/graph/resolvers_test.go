package graph

import (
	"homework-backend/db"
	"homework-backend/graph/generated"
	"homework-backend/graph/model"
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
)

func TestResolvers(t *testing.T) {
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &Resolver{}}))
	c := client.New(srv)

	err := db.InitializeDB()
	if err != nil {
		t.Error(err)
	}

	t.Run("testing valid answers", func(t *testing.T) {
		var resp struct {
			SubmitAnswers struct {
				Content string
			}
		}

		selectedOptionId := "200"
		selectedOptionWeight := 1.0
		enteredText := "Spicy food"

		c.MustPost(`
			mutation($answers: [AnswerInput!]!) {
				submitAnswers(answers: $answers) {
					content
				}
			}
		`,
			&resp,
			client.Var("answers", []model.AnswerInput{
				{QuestionID: "100",
					QuestionType:         "ChoiceQuestion",
					QuestionWeight:       0.5,
					EnteredText:          nil,
					SelectedOptionID:     &selectedOptionId,
					SelectedOptionWeight: &selectedOptionWeight},

				{QuestionID: "101",
					QuestionType:         "TextQuestion",
					QuestionWeight:       1.0,
					EnteredText:          &enteredText,
					SelectedOptionID:     nil,
					SelectedOptionWeight: nil},
			}))

		if resp.SubmitAnswers.Content != "Successfully submitted your answers" {
			t.Error("Failed to submit valid answers")
		}
	})

	t.Run("testing choice question with no selected option", func(t *testing.T) {
		var resp struct {
			SubmitAnswers struct {
				Content string
			}
		}

		enteredText := "Spicy food"

		err := c.Post(`
			mutation($answers: [AnswerInput!]!) {
				submitAnswers(answers: $answers) {
					content
				}
			}
		`,
			&resp,
			client.Var("answers", []model.AnswerInput{
				{QuestionID: "100",
					QuestionType:         "ChoiceQuestion",
					QuestionWeight:       0.5,
					EnteredText:          nil,
					SelectedOptionID:     nil,
					SelectedOptionWeight: nil},

				{QuestionID: "101",
					QuestionType:         "TextQuestion",
					QuestionWeight:       1.0,
					EnteredText:          &enteredText,
					SelectedOptionID:     nil,
					SelectedOptionWeight: nil},
			}))

		if err == nil {
			t.Error("Failed to validate option id on answers")
		}
	})

	t.Run("testing choice question with no selected weight", func(t *testing.T) {
		var resp struct {
			SubmitAnswers struct {
				Content string
			}
		}

		selectedOptionId := "200"
		enteredText := "Spicy food"

		err := c.Post(`
			mutation($answers: [AnswerInput!]!) {
				submitAnswers(answers: $answers) {
					content
				}
			}
		`,
			&resp,
			client.Var("answers", []model.AnswerInput{
				{QuestionID: "100",
					QuestionType:         "ChoiceQuestion",
					QuestionWeight:       0.5,
					EnteredText:          nil,
					SelectedOptionID:     &selectedOptionId,
					SelectedOptionWeight: nil},

				{QuestionID: "101",
					QuestionType:         "TextQuestion",
					QuestionWeight:       1.0,
					EnteredText:          &enteredText,
					SelectedOptionID:     nil,
					SelectedOptionWeight: nil},
			}))

		if err == nil {
			t.Error("Failed to validate option weight on answers")
		}
	})

	t.Run("testing mismatched questions and answers", func(t *testing.T) {
		var resp struct {
			SubmitAnswers struct {
				Content string
			}
		}

		enteredText := "Spicy food"

		err := c.Post(`
			mutation($answers: [AnswerInput!]!) {
				submitAnswers(answers: $answers) {
					content
				}
			}
		`,
			&resp,
			client.Var("answers", []model.AnswerInput{
				{QuestionID: "101",
					QuestionType:         "TextQuestion",
					QuestionWeight:       1.0,
					EnteredText:          &enteredText,
					SelectedOptionID:     nil,
					SelectedOptionWeight: nil},
			}))

		if err == nil {
			t.Error("Failed to validate questions and answers length")
		}
	})

	t.Run("testing unknown question type", func(t *testing.T) {
		var resp struct {
			SubmitAnswers struct {
				Content string
			}
		}

		enteredText := "Spicy food"

		err := c.Post(`
			mutation($answers: [AnswerInput!]!) {
				submitAnswers(answers: $answers) {
					content
				}
			}
		`,
			&resp,
			client.Var("answers", []model.AnswerInput{
				{QuestionID: "101",
					QuestionType:         "EssayQuestion",
					QuestionWeight:       1.0,
					EnteredText:          &enteredText,
					SelectedOptionID:     nil,
					SelectedOptionWeight: nil},
			}))

		if err == nil {
			t.Error("Failed to validate question type")
		}
	})

	t.Run("Get results", func(t *testing.T) {
		var resp struct {
			Results struct {
				TotalQuestions int
				CorrectAnswers int
				Score          float64
			}
		}

		var mutationResp struct {
			SubmitAnswers struct {
				Content string
			}
		}

		// send some data to populate the db
		selectedOptionId := "200"
		selectedOptionWeight := 1.0
		enteredText := "Spicy food"

		c.MustPost(`
			mutation($answers: [AnswerInput!]!) {
				submitAnswers(answers: $answers) {
					content
				}
			}
		`,
			&mutationResp,
			client.Var("answers", []model.AnswerInput{
				{QuestionID: "100",
					QuestionType:         "ChoiceQuestion",
					QuestionWeight:       0.5,
					EnteredText:          nil,
					SelectedOptionID:     &selectedOptionId,
					SelectedOptionWeight: &selectedOptionWeight},

				{QuestionID: "101",
					QuestionType:         "TextQuestion",
					QuestionWeight:       1.0,
					EnteredText:          &enteredText,
					SelectedOptionID:     nil,
					SelectedOptionWeight: nil},
			}))

		c.MustPost(`
			query {
				results {
					totalQuestions
					correctAnswers
					score
				}
			}`, &resp)

		if resp.Results.TotalQuestions <= 0 {
			t.Error("Total questions was 0 or less")
		}

		if resp.Results.Score < 0 || resp.Results.Score > 100 {
			t.Error("Score out of bounds")
		}
	})

	t.Run("Get questions", func(t *testing.T) {
		var resp struct {
			Questions []struct {
				ID      string
				Weight  float64
				Type    string
				Body    string
				Options []model.Option
			}
		}

		c.MustPost(`
		query {
			questions {
				type: __typename
				... on ChoiceQuestion {
					id
					body
					weight
					options {
						id
						body
						weight
					}
				}
				... on TextQuestion {
					id
					body
					weight
				}
			}
		}`, &resp)

		for _, question := range resp.Questions {
			if question.Weight <= 0 {
				t.Error("Received question with weight 0 or less")
			}

			if question.Body == "" {
				t.Error("Received question with empty body")
			}

			if question.ID == "" {
				t.Error("Received question with empty ID")
			}

			if question.Type == "ChoiceQuestion" {
				if len(question.Options) <= 1 {
					t.Error("Received choice question with 1 or less options")
				}
			}
		}
	})
}
