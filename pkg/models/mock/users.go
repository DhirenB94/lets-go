package mock

import (
	"time"

	models "dhiren.brahmbhatt/snippetbox/pkg"
)

type MockUserModel struct{}

var mockUser = models.User{
	ID:      1,
	Name:    "mock name",
	Email:   "mock@email.com",
	Created: time.Now(),
}

func (mu *MockUserModel) Insert(name, email, password string) error {
	return nil
}

func (mu *MockUserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

func (mu *MockUserModel) Get(id int) (*models.User, error) {
	return &mockUser, nil
}
