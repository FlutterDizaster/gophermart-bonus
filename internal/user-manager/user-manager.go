package usermanager

import (
	"context"
	"encoding/hex"

	"github.com/FlutterDizaster/gophermart-bonus/internal/models"
	serr "github.com/FlutterDizaster/gophermart-bonus/internal/shared-errors"
	"github.com/FlutterDizaster/gophermart-bonus/pkg/validation"
)

type UserRepository interface {
	UserExist(ctx context.Context, username, passHash string) bool
	AddUser(ctx context.Context, username, passHash string) error
}

type TokenResolver interface {
	CreateToken(issuer, subject string) (string, error)
}

type Settings struct {
	Repo     UserRepository
	Resolver TokenResolver
}

type UserManager struct {
	repo     UserRepository
	resolver TokenResolver
}

func New(settings Settings) *UserManager {
	return &UserManager{
		repo:     settings.Repo,
		resolver: settings.Resolver,
	}
}

func (um *UserManager) Register(ctx context.Context, cred models.Credentials) (string, error) {
	// Подсчет хеша пароля
	passHash := validation.CalculateHashSHA256([]byte(cred.Password))

	// Добавление пользователя в репозиторий
	err := um.repo.AddUser(ctx, cred.Login, hex.EncodeToString(passHash))
	if err != nil {
		return "", err
	}

	// Возврат нового токена
	return um.resolver.CreateToken("gophermart-bonus", cred.Login)
}

func (um *UserManager) Login(ctx context.Context, cred models.Credentials) (string, error) {
	// Подсчет хеша пароля
	passHash := validation.CalculateHashSHA256([]byte(cred.Password))

	// Проверка есть ли в репозитории запись с подходящими данными
	if !um.repo.UserExist(ctx, cred.Login, hex.EncodeToString(passHash)) {
		return "", serr.ErrWrongLoginOrPassword
	}

	// Возврат нового токена
	return um.resolver.CreateToken("gophermart-bonus", cred.Login)
}
