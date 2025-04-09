package service

import (
	"errors"
	"fmt"
	"inventory-service/config"
	"inventory-service/internal/adapter/http/service/handlers"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

const serverIPAddress = "127.0.0.1:%d" // Changed to 0.0.0.0 for external access

type API struct {
	server *gin.Engine
	cfg    config.HTTPServer
	addr   string

	inventoryHandler *handlers.Inventory
}

func New(cfg config.Server, inventoryUseCase InventoryUsecase) *API {
	// Setting the Gin mode
	gin.SetMode(cfg.HTTPServer.Mode)
	// Creating a new Gin Engine
	server := gin.New()

	// Applying middleware
	server.Use(gin.Logger())
	server.Use(gin.Recovery())

	// Binding inventory
	inventoryHandler := handlers.NewInventory(inventoryUseCase)

	api := &API{
		server:           server,
		cfg:              cfg.HTTPServer,
		addr:             fmt.Sprintf(serverIPAddress, cfg.HTTPServer.Port),
		inventoryHandler: inventoryHandler,
	}

	api.setupRoutes()

	return api

}
func (a *API) setupRoutes() {
	a.server.GET("/healthcheck", a.HealthCheck)

	products := a.server.Group("/products")
	{
		products.POST("/", a.inventoryHandler.Create)
		products.GET("/", a.inventoryHandler.GetList)
		products.GET("/:id", a.inventoryHandler.GetByID)
		products.PATCH("/:id", a.inventoryHandler.Update)
		products.DELETE("/:id", a.inventoryHandler.Delete)
	}
}

func (a *API) Stop() error {
	return nil
}

func (a *API) Run(errCh chan<- error) {
	go func() {
		log.Printf("HTTP server starting on: %v", a.addr)

		// No need to reinitialize `a.server` here. Just run it directly.
		if err := a.server.Run(a.addr); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- fmt.Errorf("failed to start HTTP server: %w", err)
			return
		}
	}()
}

func (a *API) HealthCheck(c *gin.Context) {
	c.JSON(200, map[string]any{
		"status": "available",
		"system_info": map[string]string{
			"address": a.addr,
			"mode":    a.cfg.Mode,
		},
	})
}
