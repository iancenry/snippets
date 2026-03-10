package mocks

import (
	"time"

	"github.com/google/uuid"
	"github.com/iancenry/snippetbox/internal/models"
)


var mockSnippet = &models.Snippet{
	ID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
	Title: "Test Snippet 1",
	Content: "This is the first test snippet.",
	Created: time.Date(2024, time.January, 1, 12, 0, 0, 0, time.UTC),
	Expires: time.Date(2024, time.January, 2, 12, 0, 0, 0, time.UTC),
}

type SnippetModel struct {}

func (m *SnippetModel) Insert(title, content, expires string) (uuid.UUID, error){
	return uuid.MustParse("00000000-0000-0000-0000-000000000001"), nil
}

func (m *SnippetModel) Get(id uuid.UUID) (*models.Snippet, error){
	switch id {
	case uuid.MustParse("00000000-0000-0000-0000-000000000001"):
		return mockSnippet, nil
	case uuid.MustParse("00000000-0000-0000-0000-000000000002"):
		return nil, nil
	default:
		return nil, models.ErrNoRecord
	}
	
}

func (m *SnippetModel) Latest() ([]*models.Snippet, error){
	return []*models.Snippet{mockSnippet}, nil
}