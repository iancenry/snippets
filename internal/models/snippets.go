package models

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Snippet struct {
	ID uuid.UUID
	Title string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModelInterface interface {
	Insert(title, content string, expires int) (uuid.UUID, error)
	Get(id uuid.UUID) (*Snippet, error)
	Latest() ([]*Snippet, error)
}

// wraps a sql.DB connection pool and exposes methods for working with snippets data
type SnippetModel struct {
	DB *pgxpool.Pool
}

// inserts a new snippet into the database and returns the id of the newly inserted record
func (m *SnippetModel) Insert(title, content string, expires int) (uuid.UUID, error) {
	var id uuid.UUID
 
	stmt := `
		INSERT INTO snippets (title, content, created, expires)
		VALUES ($1, $2, current_timestamp, current_timestamp + interval '1 day' * $3)
		RETURNING id`
		
	err := m.DB.QueryRow(context.Background(), stmt, title, content, expires).Scan(&id)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

// returns a specific snippet based on its id. Returns a models.ErrNoRecord error if the id is invalid or there is no matching record in the database
func (m *SnippetModel) Get(id uuid.UUID) (*Snippet, error) {
	row := m.DB.QueryRow(context.Background(), `
		SELECT id, title, content, created, expires FROM snippets
		WHERE expires > current_timestamp AND id = $1`, id)

	s := &Snippet{}

	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return s, nil
}

// returns the 10 most recently created snippets that have not yet expired
func (m *SnippetModel) Latest() ([]*Snippet, error) {
	rows, err := m.DB.Query(context.Background(), `
		SELECT id, title, content, created, expires FROM snippets
		WHERE expires > current_timestamp ORDER BY created DESC LIMIT 10`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	snippets := []*Snippet{}

	for rows.Next() {
		s := &Snippet{}
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, ErrNoRecord
		}
		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}