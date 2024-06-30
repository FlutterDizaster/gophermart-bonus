// This file is part of the gophermart-bonus project
//
// © 2024 Dmitriy Loginov
//
// Licensed under the MIT License.
//
// See the LICENSE.md file in the project root for more information.
//
// https://github.com/FlutterDizaster/gophermart-bonus
package balancemanager

import (
	"context"

	"github.com/FlutterDizaster/gophermart-bonus/internal/models"
)

type BalanceRepository interface {
	GetUserBalance(ctx context.Context, userID uint64) (models.Balance, error)
	GetUserWithdrawals(ctx context.Context, userID uint64) (models.Withdrawals, error)
	ProcessWithdraw(ctx context.Context, userID uint64, withdraw models.Withdraw) error
}

type Settings struct {
	BalanceRepo BalanceRepository
}

type BalanceManager struct {
	balanceRepo BalanceRepository
}

func New(settings Settings) *BalanceManager {
	return &BalanceManager{
		balanceRepo: settings.BalanceRepo,
	}
}

func (m *BalanceManager) Get(ctx context.Context, userID uint64) (models.Balance, error) {
	return m.balanceRepo.GetUserBalance(ctx, userID)
}

func (m *BalanceManager) ProcessWithdraw(
	ctx context.Context,
	userID uint64,
	withdraw models.Withdraw,
) error {
	return m.balanceRepo.ProcessWithdraw(ctx, userID, withdraw)
}

func (m *BalanceManager) GetWithdrawals(
	ctx context.Context,
	userID uint64,
) (models.Withdrawals, error) {
	return m.balanceRepo.GetUserWithdrawals(ctx, userID)
}
