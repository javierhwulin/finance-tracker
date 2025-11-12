package user

import (
	"time"

	"github.com/google/uuid"
)

func NewUser(email, password string) (*User, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	user := &User{
		ID:        id,
		Email:     email,
		Password:  password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := user.Validate(); err != nil {
		return nil, err
	}
	return user, nil
}
