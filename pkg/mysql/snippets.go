package mysql

import (
	"database/sql"

	models "dhiren.brahmbhatt/snippetbox/pkg"
)

type SnippetModel struct {
	DB *sql.DB
}

func (sm *SnippetModel) Insert(title, content, expires string) (int, error) {
	query := `INSERT INTO snippets (title, content, created, expires)
	VALUES (?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := sm.DB.Exec(query, title, content, expires)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (sm *SnippetModel) Get(id int) (*models.Snippet, error) {
	query := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP AND id = ?`

	//this will return a pointer to a sql.Row database object
	row := sm.DB.QueryRow(query, id)

	//Initialise a new models.Snippet
	s := &models.Snippet{}

	//copy the values from each field in the sql.Row object into the corresponding fields in the Snippet Struct.
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}

	return s, nil
}

func (sm *SnippetModel) Latest() ([]*models.Snippet, error) {
	return nil, nil
}
