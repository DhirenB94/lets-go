package mock

import (
	"time"

	models "dhiren.brahmbhatt/snippetbox/pkg/models"
)

type MockUserModel struct{}

var mockUser = &models.User{
	ID:      1,
	Name:    "mock name",
	Email:   "mock@email.com",
	Created: time.Now(),
}

func (mu *MockUserModel) Insert(name, email, password string) error {
	switch email {
	case "dupe@example.com":
		return models.ErrDuplicateEmail
	default:
		return nil
	}
}

func (mu *MockUserModel) Authenticate(email, password string) (int, error) {
	switch email {
	case "mock@email.com":
		return 1, nil
	default:
		return 0, models.ErrInvalidCredentials
	}
}

func (mu *MockUserModel) Get(id int) (*models.User, error) {
	switch id {
	case 1:
		return mockUser, nil
	default:
		return nil, models.ErrNoRecord
	}
}
