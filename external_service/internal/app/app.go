package app

import (
	"context"
	"log/slog"
	"os"

	"ext_service/internal/config"
	tgservice "ext_service/internal/interconnect/tg_service"
	"ext_service/internal/service"

	"golang.org/x/sync/errgroup"
)

type App struct {
	cfg         config.Config
	restService *service.RestService
}

func NewApp(ctx context.Context, cfg config.Config) (*App, error) {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	tgServiceClient, err := tgservice.NewClient(ctx, cfg.TGService)
	if err != nil {
		return nil, err
	}

	restService, err := service.NewRestService(ctx, cfg.Rest, tgServiceClient)
	if err != nil {
		return nil, err
	}

	return &App{
		cfg:         cfg,
		restService: restService,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	slog.Info("running app")
	wg, ctx := errgroup.WithContext(ctx)

	wg.Go(func() error {
		return a.restService.Run(ctx)
	})

	<-ctx.Done()
	return ctx.Err()
}
