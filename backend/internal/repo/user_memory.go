package repo

import (
	"errors"
	"sync"

	"github.com/google/uuid"
	"github.com/javierhwulin/finance-tracker/internal/domain/user"
)

// UserMemoryRepository implements user.UserRepository with in-memory storage
type UserMemoryRepository struct {
	mu         sync.RWMutex
	users      map[uuid.UUID]*user.User
	emailIndex map[string]uuid.UUID
}

// NewUserMemoryRepository creates a new in-memory user repository
func NewUserMemoryRepository() *UserMemoryRepository {
	return &UserMemoryRepository{
		users:      make(map[uuid.UUID]*user.User),
		emailIndex: make(map[string]uuid.UUID),
	}
}

// Create stores a new user
func (r *UserMemoryRepository) Create(u *user.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if email already exists
	if _, exists := r.emailIndex[u.Email]; exists {
		return errors.New("email already exists")
	}

	r.users[u.ID] = u
	r.emailIndex[u.Email] = u.ID
	return nil
}

// GetByID retrieves a user by ID
func (r *UserMemoryRepository) GetByID(id uuid.UUID) (*user.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	u, ok := r.users[id]
	if !ok {
		return nil, errors.New("user not found")
	}
	return u, nil
}

// GetByEmail retrieves a user by email
func (r *UserMemoryRepository) GetByEmail(email string) (*user.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	id, ok := r.emailIndex[email]
	if !ok {
		return nil, errors.New("user not found")
	}
	return r.users[id], nil
}

// Update updates an existing user
func (r *UserMemoryRepository) Update(u *user.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.users[u.ID]; !ok {
		return errors.New("user not found")
	}

	r.users[u.ID] = u
	r.emailIndex[u.Email] = u.ID
	return nil
}

// Delete removes a user
func (r *UserMemoryRepository) Delete(id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	u, ok := r.users[id]
	if !ok {
		return errors.New("user not found")
	}

	delete(r.emailIndex, u.Email)
	delete(r.users, id)
	return nil
}
