package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: no matching record found")

// ErrInvalidCredentials is used when a user tries to login with an incorrect email or password
var ErrInvalidCredentials = errors.New("models: invalid credentials")

// ErrDuplicateEmail is used when a user tries to sign up with an email that os already in use
var ErrDuplicateEmail = errors.New("models: duplicate email")

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}
