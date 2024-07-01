// This file is part of the gophermart-bonus project
//
// © 2024 Dmitriy Loginov
//
// Licensed under the MIT License. See the LICENSE.md file in the project root for more information.
//
// https://github.com/FlutterDizaster/gophermart-bonus
package server

import (
	"time"

	"github.com/FlutterDizaster/gophermart-bonus/internal/api"
	"github.com/FlutterDizaster/gophermart-bonus/internal/application"
	jwtresolver "github.com/FlutterDizaster/gophermart-bonus/internal/jwt-resolver"
	balancemanager "github.com/FlutterDizaster/gophermart-bonus/internal/managers/balance-manager"
	ordermanager "github.com/FlutterDizaster/gophermart-bonus/internal/managers/order-manager"
	usermanager "github.com/FlutterDizaster/gophermart-bonus/internal/managers/user-manager"
	"github.com/FlutterDizaster/gophermart-bonus/internal/repository/postgres"
)

type Settings struct {
	Addr        string
	DBConn      string
	AccrualAddr string
	JWTSecret   string
	SHASecret   string
}

type Server struct {
	application.Application
}

func New(settings Settings) (*Server, error) {
	// созданеи репозитория
	repo, err := postgres.New(settings.DBConn)
	if err != nil {
		return nil, err
	}

	// создание jwt resolver
	resolverSettings := &jwtresolver.Settings{
		Secret:   settings.JWTSecret,
		TokenTTL: 24 * time.Hour,
	}

	resolver := jwtresolver.New(*resolverSettings)

	// Создание user manager
	userSettings := usermanager.Settings{
		Repo:     repo,
		Resolver: resolver,
	}

	userMgr := usermanager.New(userSettings)

	// создание order manager
	balanceSettings := balancemanager.Settings{
		BalanceRepo: repo,
	}

	balanceMgr := balancemanager.New(balanceSettings)

	// созданеи order manager
	orderSettings := ordermanager.Settings{
		OrderRepo:   repo,
		AccrualAddr: settings.AccrualAddr,
	}
	orderMgr := ordermanager.New(orderSettings)

	apiSettings := api.Settings{
		OrderMgr:      orderMgr,
		BalanceMgr:    balanceMgr,
		UserMgr:       userMgr,
		Addr:          settings.Addr,
		TokenResolver: resolver,
		HashSumSecret: settings.SHASecret,
		CookieTTL:     24 * time.Hour,
	}

	apiService := api.New(apiSettings)

	// создание сервера
	srv := &Server{}

	srv.RegisterService(repo)
	srv.RegisterService(orderMgr)
	srv.RegisterService(apiService)

	return srv, nil
}
