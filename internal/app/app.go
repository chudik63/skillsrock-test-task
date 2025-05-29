package app

import (
	"context"
	"os"
	"os/signal"
	"skillsrock-test-task/internal/config"
	"skillsrock-test-task/internal/database/postgres"
	"skillsrock-test-task/internal/delivery/http/v1/handler"
	"skillsrock-test-task/internal/delivery/routes"
	"skillsrock-test-task/internal/repository"
	"skillsrock-test-task/internal/service"
	"skillsrock-test-task/pkg/logger"
	"skillsrock-test-task/pkg/migrator"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

const shutdownTimeout = 5 * time.Second

func Run() {
	log, err := logger.New()
	if err != nil {
		panic(err)
	}

	ctx := logger.SetToCtx(context.Background(), log)

	cfg, err := config.New()
	if err != nil {
		log.Fatal(ctx, "Failed to load configuration", zap.Error(err))
	}

	db, err := postgres.NewDatabase(cfg.Postgres)
	if err != nil {
		log.Fatal(ctx, "Failed to connect to the database", zap.Error(err))
	}
	defer db.Close()

	err = migrator.Start(cfg)
	if err != nil {
		log.Fatal(ctx, "Failed to run migrations", zap.Error(err))
	}

	repo := repository.NewTaskRepository(db)
	serv := service.NewTaskService(repo)

	app := fiber.New()

	routes.RegistrateRoutes(app, handler.NewHandler(serv, log))

	go func() {
		if err := app.Listen(":" + cfg.HTTP.Port); err != nil {
			log.Fatal(ctx, "Failed running the server", zap.Error(err))
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	<-c

	ctx, shutdown := context.WithTimeout(ctx, shutdownTimeout)
	defer shutdown()

	if err = app.ShutdownWithContext(ctx); err != nil {
		log.Error(ctx, "Failed shutting down the server", zap.Error(err))
	}

	log.Info(ctx, "Server gracefully stopped")
}
