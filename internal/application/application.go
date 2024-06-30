// This file is part of the gophermart-bonus project
//
// © 2024 Dmitriy Loginov
//
// Licensed under the MIT License. See the LICENSE.md file in the project root for more information.
//
// https://github.com/FlutterDizaster/gophermart-bonus
package application

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"time"

	"golang.org/x/sync/errgroup"
)

type Service interface {
	Start(ctx context.Context) error
}

type Application struct {
	services []Service
}

func (a *Application) RegisterService(service Service) {
	a.services = append(a.services, service)
}

func (a *Application) Start(ctx context.Context) error {
	slog.Debug("Starting application instance")

	if len(a.services) == 0 {
		return errors.New("error no registered services")
	}

	eg, egCtx := errgroup.WithContext(ctx)

	// запуск сервисов
	for i := range a.services {
		service := a.services[i]
		eg.Go(func() error {
			return service.Start(egCtx)
		})
	}

	// ожидание завершения контекста
	<-egCtx.Done()

	// Принудительное завершение работы через 30 секунд
	forceCtx, forceCancleCtx := context.WithTimeout(context.Background(), 30*time.Second)
	defer forceCancleCtx()
	go func() {
		<-forceCtx.Done()
		if forceCtx.Err() == context.DeadlineExceeded {
			slog.Error("shutdown timed out... forcing exit.")
			os.Exit(1)
		}
	}()

	return eg.Wait()
}
