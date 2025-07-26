package main

import (
	"context"
	"fmt"
	"github.com/mytheresa/go-hiring-challenge/app/catalog/application/usecases"
	"github.com/mytheresa/go-hiring-challenge/app/catalog/infrastructure/http/handlers"
	"github.com/mytheresa/go-hiring-challenge/app/catalog/infrastructure/persistence/repository"
	"github.com/mytheresa/go-hiring-challenge/app/shared/infrastructure/http/handlers_shared"
	"github.com/mytheresa/go-hiring-challenge/app/shared/infrastructure/persistence/database"
	"github.com/mytheresa/go-hiring-challenge/app/variants/application/usecases_variants"
	"github.com/mytheresa/go-hiring-challenge/app/variants/infrastructure/http/handlers_variants"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	// signal handling for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Initialize database connection
	db, close := database.New(
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_PORT"),
	)
	defer close()

	// Initialize repositories
	prodRepo := repository.NewProductRepository(db)

	// Initialize use cases
	getCatalogUseCase := usecases.NewGetCatalogUseCase(prodRepo)
	getProductByIDUseCase := usecases_variants.NewGetProductByIDUseCase(prodRepo)

	// Inicialize HTTP handlers_variants
	catHandler := handlers.NewCatalogHandler(getCatalogUseCase)
	productHandler := handlers_variants.NewProductHandler(getProductByIDUseCase)

	// Set up routing
	mux := http.NewServeMux()
	mux.HandleFunc("/health", handlers_shared.HealthCheck)
	mux.HandleFunc("/catalog", catHandler.GetCatalog)
	mux.HandleFunc("/catalog/", productHandler.GetProductById)

	// Set up the HTTP server
	srv := &http.Server{
		Addr:    fmt.Sprintf("localhost:%s", os.Getenv("HTTP_PORT")),
		Handler: mux,
	}

	// Start the server
	go func() {
		log.Printf("Starting server on http://%s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %s", err)
		}

		log.Println("Server stopped gracefully")
	}()

	<-ctx.Done()
	log.Println("Shutting down server...")
	srv.Shutdown(ctx)
	stop()
}
