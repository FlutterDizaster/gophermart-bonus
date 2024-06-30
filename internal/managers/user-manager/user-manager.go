// This file is part of the gophermart-bonus project
//
// © 2024 Dmitriy Loginov
//
// Licensed under the MIT License. See the LICENSE.md file in the project root for more information.
//
// https://github.com/FlutterDizaster/gophermart-bonus
package usermanager

import (
	"context"
	"encoding/hex"

	"github.com/FlutterDizaster/gophermart-bonus/internal/models"
	"github.com/FlutterDizaster/gophermart-bonus/pkg/validation"
)

type UserRepository interface {
	CheckUser(ctx context.Context, username, passHash string) (uint64, error)
	AddUser(ctx context.Context, username, passHash string) (uint64, error)
}

type TokenResolver interface {
	CreateToken(issuer, subject string, userID uint64) (string, error)
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
	userID, err := um.repo.AddUser(ctx, cred.Login, hex.EncodeToString(passHash))
	if err != nil {
		return "", err
	}

	// Возврат нового токена
	return um.resolver.CreateToken("gophermart-bonus", cred.Login, userID)
}

func (um *UserManager) Login(ctx context.Context, cred models.Credentials) (string, error) {
	// Подсчет хеша пароля
	passHash := validation.CalculateHashSHA256([]byte(cred.Password))

	// Проверка есть ли в репозитории запись с подходящими данными
	userID, err := um.repo.CheckUser(ctx, cred.Login, hex.EncodeToString(passHash))
	if err != nil {
		return "", err
	}

	// Возврат нового токена
	return um.resolver.CreateToken("gophermart-bonus", cred.Login, userID)
}
