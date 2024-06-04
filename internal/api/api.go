package api

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/FlutterDizaster/gophermart-bonus/internal/api/middleware"
	"github.com/FlutterDizaster/gophermart-bonus/internal/models"
	"github.com/go-chi/chi/v5"
	"golang.org/x/sync/errgroup"
)

type BalanceManager interface {
	Get(ctx context.Context, username string) (models.Balance, error)
	ProcessWithdraw(ctx context.Context, username string, withdraw models.Withdraw) error
	GetWithdrawals(ctx context.Context, username string) (models.Withdrawals, error)
}

type OrderManager interface {
	Register(ctx context.Context, username string, orderID uint64) error

	//TODO: Должен отдавать отсортированный слайс
	Get(ctx context.Context, username string) (models.Orders, error)
}

type UserManager interface {
	Register(context.Context, models.Credentials) (string, error)
	Login(context.Context, models.Credentials) (string, error)
}

type Settings struct {
	OrderMgr      OrderManager
	BalanceMgr    BalanceManager
	UserMgr       UserManager
	Addr          string
	TokenResolver middleware.TokenResolver
	HashSumSecret string
}

type API struct {
	orderMgr   OrderManager
	BalanceMgr BalanceManager
	userMgr    UserManager
	server     *http.Server
}

func New(settings Settings) *API {
	slog.Debug("Creating API service")
	// Создание экземпляра API
	api := &API{
		orderMgr:   settings.OrderMgr,
		BalanceMgr: settings.BalanceMgr,
		userMgr:    settings.UserMgr,
	}

	// Создание роутера
	r := chi.NewRouter()

	// Установка middlewares
	compressorMiddleware := middleware.Compressor{
		MinDataLength: 1024,
	}
	r.Use(compressorMiddleware.Handle)

	decompressorMiddleware := middleware.Decompressor{}
	r.Use(decompressorMiddleware.Handle)

	validatorMiddleware := middleware.Validator{
		Key: []byte(settings.HashSumSecret),
	}
	r.Use(validatorMiddleware.Handle)

	authMiddleware := middleware.AuthMiddleware{
		Resolver: settings.TokenResolver,
	}
	r.Use(authMiddleware.Handle)

	// Настройка роутинга
	r.Route("/api/user", func(r chi.Router) {
		r.Post("/register/", api.registerHandler)
		r.Post("/login/", api.loginHandler)
		r.Post("/orders/", api.ordersPOSTHandler)
		r.Get("/orders/", api.ordersGETHandler)
		r.Route("/balance", func(r chi.Router) {
			r.Get("/", api.balanceHandler)
			r.Post("/withdraw/", api.withdrawHandler)
		})
		r.Get("/withdrawals/", api.withdrawalsHandler)
	})

	// настройка ответов на не обрабатываемые сервером запросы
	r.NotFound(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	r.MethodNotAllowed(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	})

	// Создание http сервера
	api.server = &http.Server{
		Addr:    settings.Addr,
		Handler: r,
	}

	slog.Debug("API service created")

	return api
}

// Функция запуска сервиса.
func (api *API) Start(ctx context.Context) error {
	slog.Info("Starting API service")
	defer slog.Info("API server succesfully stopped")
	eg, egCtx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		slog.Info("Listening...")
		err := api.server.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	})

	<-egCtx.Done()
	eg.Go(func() error {
		slog.Info("Shutingdown API service")
		return api.server.Shutdown(context.TODO())
	})

	return eg.Wait()
}
