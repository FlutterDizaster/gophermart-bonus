package usermanager

import (
	"context"

	"github.com/FlutterDizaster/gophermart-bonus/internal/models"
)

type UserRepository interface {
	UserExist(ctx context.Context, username, passHash string) bool
	AddUser(ctx context.Context, username, passHash string) error
}

type TokenResolver interface {
	// DecryptToken(tokenString string) (*models.Claims, error)
	CreateToken(claims models.Claims) (string, error)
}

type Settings struct {
	UsrRepo     UserRepository
	TknResolver TokenResolver
}

type UserManager struct {
	usrRepo     UserRepository
	tknResolver TokenResolver
}

func New(settings Settings) *UserManager {
	return &UserManager{
		usrRepo:     settings.UsrRepo,
		tknResolver: settings.TknResolver,
	}
}

func (um *UserManager) Register(context.Context, models.Credentials) (string, error) {
	// Добавление пользователя в репозиторий

	return "", nil
}

func (um *UserManager) Login(context.Context, models.Credentials) (string, error) {
	return "", nil
}
