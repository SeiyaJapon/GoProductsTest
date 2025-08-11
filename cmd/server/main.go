package main

import (
	"context"
	"fmt"
	"github.com/mytheresa/go-hiring-challenge/app/catalog"
	"github.com/mytheresa/go-hiring-challenge/app/shared"
	"github.com/mytheresa/go-hiring-challenge/app/variants"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
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
	db, close := shared.New(
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_PORT"),
	)
	defer close()

	// Initialize repositories
	prodRepo := catalog.NewProductRepository(db)
	variantRepo := variants.NewVariantRepository(db, prodRepo)

	// Initialize use cases
	getCatalogUseCase := catalog.NewGetCatalogUseCase(prodRepo)
	getProductByIDUseCase := variants.NewGetProductByIDUseCase(variantRepo)

	// Inicialize HTTP handlers_variants
	catHandler := catalog.NewCatalogHandler(getCatalogUseCase)
	productHandler := variants.NewProductHandler(getProductByIDUseCase)

	// Set up routing

	router := mux.NewRouter()
	router.HandleFunc("/health", shared.HealthCheck).Methods("GET")
	router.HandleFunc("/catalog", catHandler.GetCatalog).Methods("GET")
	router.HandleFunc("/catalog/", productHandler.GetProductById).Methods("GET")

	// Set up the HTTP server
	srv := &http.Server{
		Addr:    fmt.Sprintf("localhost:%s", os.Getenv("HTTP_PORT")),
		Handler: router,
	}

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed to start: %s", err)
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
