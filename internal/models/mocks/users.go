package mocks

import (
	"github.com/google/uuid"
	"github.com/iancenry/snippetbox/internal/models"
)

var mockUser = &models.User{
	ID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
	Name: "Test User",
	Email: "testuser@example.com",
}

type UserModel struct {}

func (m *UserModel) Insert(name, email, password string) (uuid.UUID, error){
	switch email {
	case "dupe@example.com":
		return uuid.Nil, models.ErrDuplicateEmail
	default:
		return uuid.New(), nil
	}
}

func (m *UserModel) Authenticate(email, password string) (uuid.UUID, error){
	switch {
	case email == "testuser@example.com" && password == "correctpassword":
		return mockUser.ID, nil
	default:
		return uuid.Nil, models.ErrInvalidCredentials
	}
}

func (m *UserModel) Exists(id uuid.UUID) (bool, error){
	switch id {
	case uuid.MustParse("00000000-0000-0000-0000-000000000001"):
		return true, nil
	default:
		return false, models.ErrNoRecord
	}
}