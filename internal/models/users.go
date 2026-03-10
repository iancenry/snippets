package models

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID uuid.UUID
	Name string
	Email string
	HashedPassword []byte
	Created time.Time
}

type UserModelInterface interface {
	Insert(name, email, password string) (uuid.UUID, error)
	Authenticate(email, password string) (uuid.UUID, error)
	Exists(id uuid.UUID) (bool, error)
}

type UserModel struct {
	DB *pgxpool.Pool
}

func (m *UserModel) Insert(name, email, password string) (uuid.UUID, error) {
	var id uuid.UUID

	stmt := `
		INSERT INTO users (name, email, hashed_password)
		VALUES ($1, $2, $3)
		RETURNING id`

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return uuid.Nil, err
	}

	err = m.DB.QueryRow(context.Background(), stmt, name, email, hashedPassword).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return uuid.Nil, ErrDuplicateEmail
		}
		return uuid.Nil, err
	}
	return id, nil
}

func (m *UserModel) Authenticate(email, password string) (uuid.UUID, error) {
	var id uuid.UUID
	var hashedPassword []byte

	row := m.DB.QueryRow(context.Background(), `
		SELECT id, hashed_password FROM users WHERE email = $1`, email)

	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, ErrInvalidCredentials
		} else {
			return uuid.Nil, err
		}
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return uuid.Nil, ErrInvalidCredentials
		} else {
			return uuid.Nil, err
		}
	}

	return id, nil
}

func (m *UserModel) Exists(id uuid.UUID) (bool, error) {
	var exists bool

	row := m.DB.QueryRow(context.Background(), `
		SELECT EXISTS(SELECT true FROM users WHERE id = $1)`, id)

	err := row.Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}