package app

import (
	"context"
	"fmt"
	"inventory-service/config"
	httpservice "inventory-service/internal/adapter/http/service"
	postgresrepo "inventory-service/internal/adapter/postgres"
	"inventory-service/internal/usecase"
	"inventory-service/pkg/postgres"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const serviceName = "Inventory"

type Application struct {
	httpServer *httpservice.API
	postgresDB *postgres.PostgreDB
	// grpcServer *grpc.Server // Example
}

func New(ctx context.Context, config *config.Config) (*Application, error) {
	log.Printf("starting %v service\n", serviceName)
	log.Println("connecting to postgres")

	postgresDB, err := postgres.New(ctx, config.Postgres)
	if err != nil {
		return nil, fmt.Errorf("mongo: %w", err)
	}
	log.Println("connection established")

	inventoryRepo := postgresrepo.NewInventoryRepository(postgresDB.Pool)

	inventoryUseCase := usecase.NewInventory(inventoryRepo)
	httpServer := httpservice.New(config.Server, inventoryUseCase)

	app := &Application{
		httpServer: httpServer,
		postgresDB: postgresDB,
	}

	return app, nil
}

func (a *Application) Close() {
	// Closing http server
	err := a.httpServer.Stop()

	// Closing postgres connection
	a.postgresDB.Pool.Close()

	if err != nil {
		log.Println("failed to shutdown service", err)
	}
}

func (app *Application) Run() error {
	errCh := make(chan error, 1)

	// Running http server
	app.httpServer.Run(errCh)

	log.Printf("service %v started\n", serviceName)

	// Waiting signal
	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case errRun := <-errCh:
		return errRun

	case s := <-shutdownCh:
		log.Printf("received signal: %v. Running graceful shutdown...\n", s)

		app.Close()
		log.Println("graceful shutdown completed!")
	}

	return nil
}
