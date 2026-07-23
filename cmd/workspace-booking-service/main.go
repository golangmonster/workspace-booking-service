package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golangmonster/workspace-booking-service/internal/app/booking"
	"github.com/golangmonster/workspace-booking-service/internal/app/user"
	"github.com/golangmonster/workspace-booking-service/internal/app/workspace"
	"github.com/golangmonster/workspace-booking-service/internal/controller"
	log "github.com/sirupsen/logrus"

	"github.com/golangmonster/pgxtransactor"
	"github.com/golangmonster/workspace-booking-service/internal/config"
	bookingRepository "github.com/golangmonster/workspace-booking-service/internal/repository/booking"
	userRepository "github.com/golangmonster/workspace-booking-service/internal/repository/user"
	workspaceRepository "github.com/golangmonster/workspace-booking-service/internal/repository/workspace"
	bookingService "github.com/golangmonster/workspace-booking-service/internal/service/booking"
	userService "github.com/golangmonster/workspace-booking-service/internal/service/user"
	workspaceService "github.com/golangmonster/workspace-booking-service/internal/service/workspace"
	"github.com/kelseyhightower/envconfig"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	shutdownTime = 10 * time.Second
)

func main() {
	var cfg config.Config

	log.SetFormatter(&log.JSONFormatter{})

	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal("failed to load config ", err)
	}

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, cfg.PostgresDSN)
	if err != nil {
		log.Fatal("failed to connect to postgres ", err)
	}
	defer pool.Close()

	pgxTx := pgxtransactor.New(pool)

	userRepo := userRepository.New(pgxTx)
	workspaceRepo := workspaceRepository.New(pgxTx)
	bookingRepo := bookingRepository.New(pgxTx)

	userSrv := userService.New(userRepo)
	workspaceSrv := workspaceService.New(workspaceRepo)
	bookingSrv := bookingService.New(bookingRepo)

	ctrl := controller.New(&cfg,
		user.New(userSrv),
		workspace.New(workspaceSrv),
		booking.New(bookingSrv),
	)

	ctrl.Run(ctx)

	sigch := make(chan os.Signal, 1)

	signal.Notify(sigch, syscall.SIGTERM)
	signal.Notify(sigch, syscall.SIGINT)

	<-sigch

	shutdownCtx, cancel := context.WithTimeout(ctx, shutdownTime)
	defer cancel()

	ctrl.Stop(shutdownCtx)

	log.Info("service finished")
}
