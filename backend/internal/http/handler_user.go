package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/javierhwulin/finance-tracker/internal/domain/user"
	"github.com/javierhwulin/finance-tracker/internal/http/dto"
)

// Reusable validator instance
var validate = validator.New()

// CreateUserHandler handles user creation requests
func CreateUserHandler(userRepository user.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.CreateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := validate.Struct(req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		newUser, err := user.NewUser(req.Email, req.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := userRepository.Create(newUser); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := dto.CreateUserResponse{
			ID:    newUser.ID,
			Email: newUser.Email,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}
}
