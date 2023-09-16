package mysql

import (
	"database/sql"

	models "dhiren.brahmbhatt/snippetbox/pkg"
)

type UserModel struct {
	DB *sql.DB
}

// Insert method will add a new record to the users table.
func (um *UserModel) Insert(name, email, password string) error {
	return nil
}

// Authenticate will verify whether a user exists with the provided email and password.
// This will return the relevant user ID if they do.
func (um *UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

// Get will fetch specific details for a user based on their ID
func (um *UserModel) Get(id int) (models.User, error) {
	return models.User{}, nil
}
