package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"homework-backend/db"
	"homework-backend/graph/generated"
	"homework-backend/graph/model"
	"strconv"
	"strings"
)

func (r *mutationResolver) SubmitAnswers(ctx context.Context, answers []*model.AnswerInput) (*model.Message, error) {
	conn := db.Pool.Get(ctx)
	if conn == nil {
		return nil, errors.New("An error has occurred, please try again later")
	}
	defer db.Pool.Put(conn)

	// not sure if this is how i'm supposed to do the validation ??
	questions, err := r.Query().Questions(ctx)
	if err != nil {
		return nil, err
	}

	if len(questions) != len(answers) {
		return nil, errors.New("Submitted answers are either less or more than the provided questions")
	}

	for index, answer := range answers {
		if question, ok := questions[index].(*model.ChoiceQuestion); ok {
			if question.ID != answer.QuestionID {
				return nil, errors.New("Invalid question id on submitted answer")
			}

			if answer.QuestionType != "ChoiceQuestion" {
				return nil, errors.New("Invalid question type on submitted answer")
			}

			if !containsOptionId(question.Options, answer.SelectedOptionID) {
				return nil, errors.New("Invalid option id on submitted answer")
			}
		} else if question, ok := questions[index].(*model.TextQuestion); ok {
			if question.ID != answer.QuestionID {
				return nil, errors.New("Invalid question id on submitted answer")
			}

			if answer.QuestionType != "TextQuestion" {
				return nil, errors.New("Invalid question type on submitted answer")
			}
		} else {
			return nil, errors.New("Unknown question type")
		}
	}

	insertAnswer := `INSERT INTO answers(
		questionId, questionWeight, questionType, selectedOptionId, selectedOptionWeight, enteredText) 
	VALUES($qId, $qWeight, $qType, $optionId, $optionWeight, $enteredText);`

	for _, answer := range answers {
		stmt := conn.Prep(insertAnswer)

		parsedQId, _ := strconv.Atoi(answer.QuestionID)
		stmt.SetInt64("$qId", int64(parsedQId))
		stmt.SetFloat("$qWeight", answer.QuestionWeight)
		stmt.SetText("$qType", answer.QuestionType)

		if answer.QuestionType == "ChoiceQuestion" {
			parsedId, _ := strconv.Atoi(*answer.SelectedOptionID)
			stmt.SetInt64("$optionId", int64(parsedId))
			stmt.SetFloat("$optionWeight", *answer.SelectedOptionWeight)
			stmt.SetNull("$enteredText")
		} else if answer.QuestionType == "TextQuestion" {
			stmt.SetNull("$optionId")
			stmt.SetNull("$optionWeight")
			stmt.SetText("$enteredText", *answer.EnteredText)
		}
		stmt.Step()
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
	conn := db.Pool.Get(ctx)
	if conn == nil {
		return nil, errors.New("Failed to establish db connection")
	}
	defer db.Pool.Put(conn)

	stmt := conn.Prep(`SELECT questions.*,
	group_concat(options.id) as options_ids, group_concat(options.body, '-|-') as options_bodies, group_concat(options.weight) as options_weights
	FROM questions
	LEFT JOIN options
	on options.questionId = questions.id
	GROUP BY questions.id;`)

	questions := make([]model.Question, 0)

	for {
		if hasRow, err := stmt.Step(); err != nil {
			return nil, errors.New("An error has occurred while fetching questions, please try again later")
		} else if !hasRow {
			break
		}

		questionId := stmt.GetInt64("id")
		questionBody := stmt.GetText("body")
		questionWeight := stmt.GetFloat("weight")
		questionType := stmt.GetText("type")
		optionsIds := stmt.GetText("options_ids")
		optionsBodies := stmt.GetText("options_bodies")
		optionsWeights := stmt.GetText("options_weights")

		if questionType == "ChoiceQuestion" {
			options := make([]*model.Option, 0)

			ids := strings.Split(optionsIds, ",")
			bodies := strings.Split(optionsBodies, "-|-")
			weights := strings.Split(optionsWeights, ",")

			for index := range ids {
				weight, _ := strconv.ParseFloat(weights[index], 32)

				option := &model.Option{
					ID:     fmt.Sprintf("%v", ids[index]),
					Body:   bodies[index],
					Weight: weight,
				}

				options = append(options, option)
			}

			question := &model.ChoiceQuestion{
				ID:      fmt.Sprintf("%v", questionId),
				Body:    questionBody,
				Weight:  questionWeight,
				Options: options,
			}

			questions = append(questions, question)
		} else if questionType == "TextQuestion" {
			question := &model.TextQuestion{
				ID:     fmt.Sprintf("%v", questionId),
				Body:   questionBody,
				Weight: questionWeight,
			}

			questions = append(questions, question)
		}
	}

	return questions, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
