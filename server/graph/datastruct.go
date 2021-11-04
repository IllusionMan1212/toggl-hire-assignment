package graph

import "homework-backend/graph/model"

type Data struct {
	Answers []*model.AnswerInput `json:"answers"`
}
