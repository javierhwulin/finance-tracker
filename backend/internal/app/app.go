package app

import "github.com/javierhwulin/finance-tracker/internal/domain/user"

type App struct {
	UserRepository user.UserRepository
}

func NewApp(userRepository user.UserRepository) *App {
	return &App{UserRepository: userRepository}
}
