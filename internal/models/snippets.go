package models

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Snippet struct {
	ID uuid.UUID
	Title string
	Content string
	Created time.Time
	Expires time.Time
}

// wraps a sql.DB connection pool and exposes methods for working with snippets data
type SnippetModel struct {
	DB *pgxpool.Pool
}

func (m *SnippetModel) Insert(title, content string, expires int) (uuid.UUID, error) {
	var id uuid.UUID
	err := m.DB.QueryRow(context.Background(), `
		INSERT INTO snippets (title, content, created, expires)
		VALUES ($1, $2, current_timestamp, current_timestamp + interval '1 day' * $3)
		RETURNING id`, title, content, expires).Scan(&id)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func (m *SnippetModel) Get(id uuid.UUID) (*Snippet, error) {
	row := m.DB.QueryRow(context.Background(), `
		SELECT id, title, content, created, expires FROM snippets
		WHERE expires > current_timestamp AND id = $1`, id)

	s := &Snippet{}

	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		return nil, err
	}
	return s, nil
}

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
			return nil, err
		}
		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}