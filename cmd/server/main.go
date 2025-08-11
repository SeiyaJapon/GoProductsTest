package main

import (
	"context"
	"fmt"
	"github.com/mytheresa/go-hiring-challenge/app/catalog/application"
	httpcatalog "github.com/mytheresa/go-hiring-challenge/app/catalog/infrastructure/http"
	"github.com/mytheresa/go-hiring-challenge/app/catalog/infrastructure/persistence/mysql"
	httpshared "github.com/mytheresa/go-hiring-challenge/app/shared/infrastructure/http"
	"github.com/mytheresa/go-hiring-challenge/app/shared/infrastructure/persistence"
	application2 "github.com/mytheresa/go-hiring-challenge/app/variants/application"
	http2 "github.com/mytheresa/go-hiring-challenge/app/variants/infrastructure/http"
	mysql2 "github.com/mytheresa/go-hiring-challenge/app/variants/infrastructure/persistence/mysql"
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
	db, close := persistence.NewDBConnection(
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_PORT"),
	)
	defer close()

	// Initialize repositories
	prodRepo := mysql.NewProductRepository(db)
	variantRepo := mysql2.NewVariantRepository(db, prodRepo)

	// Initialize use cases
	getCatalogUseCase := application.NewGetCatalogUseCase(prodRepo)
	getProductByIDUseCase := application2.NewGetProductByIDUseCase(variantRepo)

	// Inicialize HTTP handlers_variants
	catHandler := httpcatalog.NewCatalogHandler(getCatalogUseCase)
	productHandler := http2.NewProductHandler(getProductByIDUseCase)

	// Set up routing

	router := mux.NewRouter()
	router.HandleFunc("/health", httpshared.HealthCheck).Methods("GET")
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
