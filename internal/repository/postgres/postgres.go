// This file is part of the gophermart-bonus project
//
// © 2024 Dmitriy Loginov
//
// Licensed under the MIT License.
//
// See the LICENSE.md file in the project root for more information.
//
// https://github.com/FlutterDizaster/gophermart-bonus
package postgres

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/FlutterDizaster/gophermart-bonus/internal/models"
	sharederrors "github.com/FlutterDizaster/gophermart-bonus/internal/shared-errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func New(conn string) (*Repository, error) {
	slog.Info("Creating Repository")
	repo := &Repository{}

	// Создание экземпляра DB
	poolConfig, err := pgxpool.ParseConfig(conn)
	if err != nil {
		return nil, err
	}

	db, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, err
	}

	repo.db = db

	return repo, nil
}

func (repo *Repository) Start(ctx context.Context) error {
	slog.Info("Starting Repository service")
	err := repo.checkAndCreateTables(ctx)
	if err != nil {
		slog.Error("error creating tables", slog.Any("error", err))
		return err
	}

	<-ctx.Done()
	return nil
}

func (repo *Repository) GetUserBalance(
	ctx context.Context,
	userID uint64,
) (models.Balance, error) {
	balance := models.Balance{}

	err := repo.db.QueryRow(ctx, userBalanceQuery, userID).
		Scan(&balance.Current, &balance.Withdrawn)
	if err != nil {
		return models.Balance{}, err
	}

	slog.Debug(
		"user balance",
		slog.Uint64("user id", userID),
		slog.Float64("balance", balance.Current),
		slog.Float64("withdrawn", balance.Withdrawn),
	)

	return balance, nil
}

func (repo *Repository) GetUserWithdrawals(
	ctx context.Context,
	userID uint64,
) (models.Withdrawals, error) {
	withdrawals := make([]models.Withdraw, 0)

	rows, err := repo.db.Query(ctx, userWithdrawalsQuery, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, sharederrors.ErrNoWithdrawalsFound
		}
		return nil, err
	}

	slogslice := make([]slog.Attr, 0)
	var counter int

	for rows.Next() {
		var w models.Withdraw
		var processedDate time.Time
		err = rows.Scan(&w.OrderID, &w.Sum, &processedDate)
		if err != nil {
			return nil, err
		}

		w.ProcessedAt = processedDate.Format(time.RFC3339)

		withdrawals = append(withdrawals, w)

		slogslice = append(
			slogslice,
			slog.Group(
				fmt.Sprintf("entry %d", counter),
				slog.Uint64("order id", w.OrderID),
				slog.Float64("sum", w.Sum),
				slog.String("processed at", w.ProcessedAt),
			),
		)
		counter++
	}

	slog.Debug(
		"user withdrawals",
		slog.Uint64(
			"user id",
			userID,
		),
		slog.Any(
			"withdrawals",
			slogslice,
		),
	)

	return withdrawals, nil
}

func (repo *Repository) ProcessWithdraw(
	ctx context.Context,
	userID uint64,
	withdraw models.Withdraw,
) error {
	// Получение баланса пользователя
	balance, err := repo.GetUserBalance(ctx, userID)
	if err != nil {
		return err
	}

	if balance.Current < withdraw.Sum {
		return sharederrors.ErrNotEnoughFunds
	}

	// Проверка кем был добавлен заказ и есть ли он
	// var existingUserID uint64
	// err = repo.db.QueryRow(
	// 	ctx,
	// 	checkOrderQuery,
	// 	withdraw.OrderID,
	// ).Scan(&existingUserID)
	// if err == nil {
	// 	if existingUserID != userID {
	// 		// Если заказ уже существует, но принадлежит другому пользователю
	// 		return sharederrors.ErrWithdrawNotAllowed
	// 	}
	// } else {
	// 	// Возвращаем необробатываемую ошибку
	// 	return err
	// }

	// Проверка не было ли списаний по этому заказу
	var status bool
	err = repo.db.QueryRow(ctx, checkWithdrawQuery, withdraw.OrderID).Scan(&status)
	if err != nil {
		return err
	}
	if status {
		return sharederrors.ErrWithdrawNotAllowed
	}

	// Списание средств
	_, err = repo.db.Exec(
		ctx,
		processWithdrawQuery,
		withdraw.OrderID,
		userID,
		withdraw.Sum,
	)

	slog.Debug(
		"process withdraw",
		slog.Uint64("user id", userID),
		slog.Group(
			"withdraw",
			slog.Uint64("order id", withdraw.OrderID),
			slog.Float64("sum", withdraw.Sum),
		),
	)

	return err
}

func (repo *Repository) AddOrder(ctx context.Context, userID uint64, order models.Order) error {
	// Проверка существует ли заказ
	var existingUserID uint64
	err := repo.db.QueryRow(
		ctx,
		checkOrderQuery,
		order.ID,
	).Scan(&existingUserID)

	if err == nil {
		// Если заказ уже существует и принадлежит текущему пользователю
		if existingUserID == userID {
			return sharederrors.ErrOrderAlreadyLoaded
		}
		// Если заказ уже существует, но принадлежит другому пользователю
		return sharederrors.ErrOrderLoadedByAnotherUsr
	}

	// Если произошла ошибка, и она не связана с отсутствием записи
	if !errors.Is(err, pgx.ErrNoRows) {
		return err
	}

	// Добавление нового заказа
	_, err = repo.db.Exec(
		ctx,
		addOrderQuery,
		order.ID,
		userID,
		order.Status,
	)
	if err != nil {
		return err
	}

	slog.Debug(
		"add order",
		slog.Uint64("user id", userID),
		slog.Group(
			"order",
			slog.Uint64("id", order.ID),
			slog.String("status", string(order.Status)),
		),
	)

	return nil
}

func (repo *Repository) UpdateOrder(ctx context.Context, order models.Order) error {
	_, err := repo.db.Exec(
		ctx,
		updateOrderQuery,
		order.Status,
		order.Accrual,
		order.ID,
	)
	if err != nil {
		return err
	}
	var slogAccr float64
	if order.Accrual != nil {
		slogAccr = *order.Accrual
	}
	slog.Info(
		"update order",
		slog.Group(
			"order",
			slog.Uint64("id", order.ID),
			slog.String("status", string(order.Status)),
			slog.Float64("accrual", slogAccr),
		),
	)
	return nil
}

func (repo *Repository) GetAllOrders(ctx context.Context, userID uint64) (models.Orders, error) {
	rows, err := repo.db.Query(
		ctx,
		getUserOrdersQuery,
		userID,
	)
	if err != nil {
		return nil, err
	}

	orders := make([]models.Order, 0)

	for rows.Next() {
		var order models.Order

		var uploadDate time.Time

		err = rows.Scan(&order.ID, &order.Status, &order.Accrual, &uploadDate)
		if err != nil {
			return nil, err
		}

		order.UploadedAt = uploadDate.Format(time.RFC3339)

		orders = append(orders, order)
	}
	if len(orders) == 0 {
		return nil, sharederrors.ErrNoOrdersFound
	}

	return orders, nil
}

func (repo *Repository) GetNotUpdatedOrders(ctx context.Context) ([]uint64, error) {
	rows, err := repo.db.Query(ctx, getNotUpdatedOrdersQuery)
	if err != nil {
		return nil, err
	}

	orders := make([]uint64, 0)

	for rows.Next() {
		var orderID uint64

		err = rows.Scan(&orderID)
		if err != nil {
			return nil, err
		}

		orders = append(orders, orderID)
	}

	return orders, nil
}

func (repo *Repository) CheckUser(ctx context.Context, username, passHash string) (uint64, error) {
	var userID uint64
	err := repo.db.QueryRow(ctx, checkUserQuery, username, passHash).Scan(&userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, sharederrors.ErrWrongLoginOrPassword
		}

		return 0, err
	}
	return userID, nil
}

func (repo *Repository) AddUser(ctx context.Context, username, passHash string) (uint64, error) {
	var userID uint64
	err := repo.db.QueryRow(ctx, addUserQuery, username, passHash).Scan(&userID)

	if err != nil {
		pgErr := &pgconn.PgError{}
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return 0, sharederrors.ErrUserAlreadyExist
			}
		}

		return 0, err
	}

	return userID, nil
}

func (repo *Repository) checkAndCreateTables(ctx context.Context) error {
	slog.Debug("Creating tables")
	_, err := repo.db.Exec(ctx, createUsersTable)
	if err != nil {
		return err
	}
	_, err = repo.db.Exec(ctx, createOrdersTable)
	if err != nil {
		return err
	}
	_, err = repo.db.Exec(ctx, createWithdrawlsTable)
	if err != nil {
		return err
	}

	return nil
}
