package graph

import "homework-backend/graph/model"

func containsOptionId(options []*model.Option, id *string) bool {
	for _, option := range options {
		if id != nil {
			if option.ID == *id {
				return true
			}
		} else {
			return false
		}
	}

	return false
}
