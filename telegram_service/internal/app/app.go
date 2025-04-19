package app

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	_ "github.com/lib/pq"
	"golang.org/x/sync/errgroup"

	"tg_service/internal/config"
	"tg_service/internal/service"
	"tg_service/store"
)

type App struct {
	cfg         config.Config
	grpcService *service.Service
	db          *sql.DB
}

func NewApp(ctx context.Context, cfg config.Config) (*App, error) {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.DB,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	repo := store.NewPostgresRepository(db)

	grpcService := service.NewService(cfg.GRPC, repo)

	return &App{
		cfg:         cfg,
		grpcService: grpcService,
		db:          db,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	defer a.db.Close()
	slog.Info("running app")
	wg, ctx := errgroup.WithContext(ctx)

	wg.Go(a.grpcService.Run)
	<-ctx.Done()
	a.grpcService.Close()
	return ctx.Err()
}
