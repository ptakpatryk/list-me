package mocks

import (
	"time"

	"github.com/ptakpatryk/list-me/internals/models"
)

var MockList = models.List{
	ID:          1,
	Title:       "Shopping list",
	Description: "shopping list description",
	Created:     time.Now(),
	Expires:     time.Now(),
}

type ListModel struct{}

func (m *ListModel) Insert(title, description string, expires int) (int, error) {
	return 2, nil
}

func (m *ListModel) Get(id int) (models.List, error) {
	switch id {
	case 1:
		return MockList, nil
	default:
		return models.List{}, models.ErrNoRecord
	}
}

func (m *ListModel) Latest() ([]models.List, error) {
	return []models.List{MockList}, nil
}
