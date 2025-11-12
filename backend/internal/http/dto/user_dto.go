package dto

import "github.com/google/uuid"

type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=255"`
}

type CreateUserResponse struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
}
