package mock

import (
	"time"

	models "dhiren.brahmbhatt/snippetbox/pkg"
)

type MockSnippetModel struct{}

var mockSnippet = &models.Snippet{
	ID:      1,
	Title:   "mock title",
	Content: "mock content",
	Created: time.Now(),
	Expires: time.Now(),
}

func (m *MockSnippetModel) Get(id int) (*models.Snippet, error) {
	return mockSnippet, nil
}

func (m *MockSnippetModel) Insert(title, content, expires string) (int, error) {
	return 0, nil
}

func (m *MockSnippetModel) Latest() ([]*models.Snippet, error) {
	return nil, nil
}
