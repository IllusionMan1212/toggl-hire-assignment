package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"homework-backend/graph/generated"
	"homework-backend/graph/model"
)

func (r *mutationResolver) SubmitAnswers(ctx context.Context, answers []*model.AnswerInput) (*model.Message, error) {
	// not sure if this is how i'm supposed to do the validation ??
	questions, err := r.Query().Questions(ctx)
	if err != nil {
		return nil, err
	}

	if len(questions) != len(answers) {
		return nil, errors.New("Submitted answers are either less or more than the provided questions")
	}

	for index, answer := range answers {
		if question, ok := questions[index].(model.ChoiceQuestion); ok {
			if question.ID != answer.QuestionID {
				return nil, errors.New("Invalid question id on submitted answer")
			}

			if answer.QuestionType != "ChoiceQuestion" {
				return nil, errors.New("Invalid question type on submitted answer")
			}

			if !containsOptionId(question.Options, answer.SelectedOptionID) {
				return nil, errors.New("Invalid option id on submitted answer")
			}
		} else if question, ok := questions[index].(model.TextQuestion); ok {
			if question.ID != answer.QuestionID {
				return nil, errors.New("Invalid question id on submitted answer")
			}

			if answer.QuestionType != "TextQuestion" {
				return nil, errors.New("Invalid question type on submitted answer")
			}
		}
	}

	data := &Data{}
	data.Answers = answers

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	fmt.Printf("%s\n", jsonData)

	return &model.Message{Content: "Successfully submitted your answers"}, nil
}

func (r *queryResolver) Questions(ctx context.Context) ([]model.Question, error) {
	return []model.Question{
		model.ChoiceQuestion{
			ID:     "100",
			Body:   "Where does the sun set?",
			Weight: 0.5,
			Options: []*model.Option{
				{ID: "200", Body: "East", Weight: 0},
				{ID: "201", Body: "West", Weight: 1},
			},
		},
		model.TextQuestion{
			ID:     "101",
			Body:   "What is your favourite food?",
			Weight: 1,
		},
	}, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
