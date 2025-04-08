package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"order-service/config"

	httpservice "order-service/internal/adapter/http/service"
	postgresrepo "order-service/internal/adapter/postgres"
	"order-service/internal/usecase"
	"order-service/pkg/postgres"
)

const serviceName = "Order"

type App struct {
	httpServer *httpservice.API
	postgresDB *postgres.PostgreDB
	// grpcServer *grpc.Server // Example
}

func New(ctx context.Context, cfg *config.Config) (*App, error) {
	log.Printf("starting %v service\n", serviceName)

	log.Println("connecting to postgres")
	postgresDB, err := postgres.New(ctx, cfg.Postgres)
	if err != nil {
		return nil, fmt.Errorf("mongo: %w", err)
	}

	log.Println("connection established")

	orderRepo := postgresrepo.NewOrderRepository(postgresDB.Pool)

	// UseCase
	orderUsecase := usecase.NewOrder(orderRepo)

	// http service
	httpServer := httpservice.New(cfg.Server, orderUsecase)

	app := &App{
		httpServer: httpServer,
		postgresDB: postgresDB,
	}

	return app, nil
}

// TODO: close postgres connection
func (a *App) Close() {
	// Closing http server
	err := a.httpServer.Stop()

	// Closing postgres connection
	a.postgresDB.Pool.Close()

	if err != nil {
		log.Println("failed to shutdown service", err)
	}
}

func (a *App) Run() error {
	errCh := make(chan error, 1)

	// Running http server
	a.httpServer.Run(errCh)

	log.Printf("service %v started\n", serviceName)

	// Waiting signal
	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case errRun := <-errCh:
		return errRun

	case s := <-shutdownCh:
		log.Printf("received signal: %v. Running graceful shutdown...\n", s)

		a.Close()
		log.Println("graceful shutdown completed!")
	}

	return nil
}