package mysql

import (
	"database/sql"
	"strings"

	models "dhiren.brahmbhatt/snippetbox/pkg"
	"github.com/go-sql-driver/mysql"

	"golang.org/x/crypto/bcrypt" // New import
)

type UserModel struct {
	DB *sql.DB
}

// Insert method will add a new record to the users table.
func (um *UserModel) Insert(name, email, password string) error {
	// Create a bcrypt hash of the plain-text password.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	query := `INSERT INTO users (name, email, hashed_password, created)
    VALUES(?, ?, ?, UTC_TIMESTAMP())`

	// Use the Exec() method to insert the user details and hashed password into the users table.
	// If this returns an error, we try to type assert it to a *mysql.MySQLError object, so we can check if the error number is 1062 (ER_DUP_ENTRY)
	// If it is, we also check whether or not the error relates to our users_uc_email key by checking the contents of the message string.
	// If it does, we return an ErrDuplicateEmail error.
	// Otherwise, we just return the original error (or nil if everything worked)
	_, err = um.DB.Exec(query, name, email, string(hashedPassword))
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if mysqlErr.Number == 1062 && strings.Contains(mysqlErr.Message, "'users.users_uc_email") {
				return models.ErrDuplicateEmail
			}
		}
	}
	return err
}

// Authenticate will verify whether a user exists with the provided email and password.
// This will return the relevant user ID if they do.
func (um *UserModel) Authenticate(email, password string) (int, error) {
	//Retrieve the id and hashed password for an emai.
	//If the email does not exist, retrun ErrInvalidCredentials
	var id int
	var hashedPassword []byte
	query := `SELECT id, hashed_password FROM users
	WHERE email = ?`

	row := um.DB.QueryRow(query, email)
	err := row.Scan(&id, &hashedPassword)
	if err == sql.ErrNoRows {
		return 0, models.ErrInvalidCredentials
	} else if err != nil {
		return 0, err
	}

	//Check whether hashed password from the db is the same as the user entered password
	//If they do not match return custom Err
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, models.ErrInvalidCredentials
	} else if err != nil {
		return 0, err
	}

	return id, nil
}

// Get will fetch specific details for a user based on their ID
func (um *UserModel) Get(id int) (models.User, error) {
	return models.User{}, nil
}
