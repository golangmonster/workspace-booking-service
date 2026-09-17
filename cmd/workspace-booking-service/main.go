package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/golangmonster/workspace-booking-service/internal/app/booking"
	"github.com/golangmonster/workspace-booking-service/internal/app/user"
	"github.com/golangmonster/workspace-booking-service/internal/app/workspace"
	"github.com/golangmonster/workspace-booking-service/internal/controller"
	"github.com/golangmonster/workspace-booking-service/internal/kafka/producer"
	completeExpiredBooking "github.com/golangmonster/workspace-booking-service/internal/process/complete-expired-booking"
	"github.com/golangmonster/workspace-booking-service/internal/process/outbox"
	log "github.com/sirupsen/logrus"

	"github.com/golangmonster/pgxtransactor"
	"github.com/golangmonster/workspace-booking-service/internal/config"
	workspaceBookingProducer "github.com/golangmonster/workspace-booking-service/internal/kafka/producer/workspace-booking"
	bookingRepository "github.com/golangmonster/workspace-booking-service/internal/repository/booking"
	outboxRepository "github.com/golangmonster/workspace-booking-service/internal/repository/outbox"
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
	outboxRepo := outboxRepository.New(pgxTx)

	userSrv := userService.New(userRepo)
	workspaceSrv := workspaceService.New(workspaceRepo)
	bookingSrv := bookingService.New(
		bookingRepo,
		userRepo,
		outboxRepo,
		workspaceBookingProducer.MarshalWorkspaceBooking,
		cfg.KafkaWorkspaceBookingTopic,
	)

	workspaceBookingProducer, err := producer.New(cfg.KafkaWorkspaceBookingBrokers, cfg.KafkaWorkspaceBookingEnabled)
	if err != nil {
		log.Fatal("new workspace booking producer: ", err)
	}
	defer func() {
		err := workspaceBookingProducer.Close()
		if err != nil {
			log.Error("close workspace booking producer", err)
		}
	}()

	completeExpiredBookingProcess := completeExpiredBooking.NewProcess(bookingRepo)
	workspaceBookingOutboxProcess := outbox.NewProcess(workspaceBookingProducer, outboxRepo, cfg.KafkaWorkspaceBookingTopic)

	scheduler, err := gocron.NewScheduler()
	if err != nil {
		log.Fatal("new scheduler: ", err)
	}

	// Complete expired bookings
	if cfg.CompleteExpiredBookingEnabled {
		_, err = scheduler.NewJob(
			gocron.DurationJob(cfg.CompleteExpiredBookingDuration),
			gocron.NewTask(completeExpiredBookingProcess.Run, ctx),
			gocron.WithSingletonMode(gocron.LimitModeReschedule),
		)
		if err != nil {
			log.Fatal("new complete expired booking job: ", err)
		}
	}

	// Outbox for workspace-booking
	if cfg.WorkspaceBookingOutboxEnabled {
		_, err = scheduler.NewJob(
			gocron.DurationJob(cfg.WorkspaceBookingOutboxDuration),
			gocron.NewTask(workspaceBookingOutboxProcess.Run, ctx),
			gocron.WithSingletonMode(gocron.LimitModeReschedule),
		)
		if err != nil {
			log.Fatal("new workspace booking outbox job: ", err)
		}
	}

	scheduler.Start()

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

	// Shutting down

	shutdownCtx, cancel := context.WithTimeout(ctx, shutdownTime)
	defer cancel()

	ctrl.Stop(shutdownCtx)

	err = scheduler.ShutdownWithContext(shutdownCtx)
	if err != nil {
		log.Fatal("scheduler shutdown: ", err)
	}

	log.Info("service finished")
}
