package main

import (
	"context"
	"database/sql"
	"log"
	"moneyget/internal/domain"
	"moneyget/internal/domain/service"
	"moneyget/internal/infrastructure/sqlite"
	"moneyget/internal/interface/handler"
	"moneyget/internal/interface/router"
	"moneyget/internal/usecase"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := initDB()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database: %v\n", err)
		}
	}()

	// Domain Services
	eventDispatcher := service.NewEventDispatcher()
	eventStore := service.NewEventStore(sqlite.NewEventStoreDB(db))
	strategyService := service.NewInvestmentStrategyService()
	passwordService, jwtService := initServices()

	// Event Handlers
	setupEventHandlers(eventDispatcher, eventStore)

	// Infrastructure Layer
	txManager := sqlite.NewTransactionManager(db)
	userRepo := sqlite.NewUserRepository(db)
	investmentRepo := sqlite.NewInvestmentRepository(db)
	portfolioRepo := sqlite.NewPortfolioRepository(db)

	// Application Layer (Use Cases)
	userUsecase := usecase.NewUserUsecase(userRepo, passwordService)
	investmentUsecase := usecase.NewInvestmentUseCase(
		investmentRepo,
		portfolioRepo,
		txManager,
		eventDispatcher,
		strategyService,
	)
	portfolioUsecase := usecase.NewPortfolioUseCase(
		portfolioRepo,
		txManager,
		eventDispatcher,
		strategyService,
	)

	// Interface Layer (Handlers)
	userHandler := handler.NewUserHandler(userUsecase, jwtService)
	investmentHandler := handler.NewInvestmentHandler(investmentUsecase)
	portfolioHandler := handler.NewPortfolioHandler(portfolioUsecase)

	// Setup and start server
	srv := setupServer(userHandler, investmentHandler, portfolioHandler, jwtService)

	// Start the server
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v\n", err)
		}
	}()

	log.Println("Server started on :8080")
	gracefulShutdown(srv)
}

func setupEventHandlers(dispatcher *service.EventDispatcher, store *service.EventStore) {
	// Investment Created Event Handler
	dispatcher.Subscribe(func(event domain.DomainEvent) {
		if err := store.SaveEvent(event); err != nil {
			log.Printf("Failed to store InvestmentCreated event: %v\n", err)
		}
	})

	// Portfolio Updated Event Handler
	dispatcher.Subscribe(func(event domain.DomainEvent) {
		if err := store.SaveEvent(event); err != nil {
			log.Printf("Failed to store PortfolioUpdated event: %v\n", err)
		}
	})
}

// 既存の関数は維持
func initDB() (*sql.DB, error) {
	dbPath := filepath.Join(".", "moneyget.db")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// コネクションプールの設定
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	if err := sqlite.RunMigrations(db); err != nil {
		return nil, err
	}

	return db, nil
}

func initServices() (service.PasswordService, service.JWTService) {
	passwordService := service.NewPasswordService()
	jwtService := service.NewJWTService("your-secret-key-here")
	return passwordService, jwtService
}

func setupServer(
	userHandler *handler.UserHandler,
	investmentHandler *handler.InvestmentHandler,
	portfolioHandler *handler.PortfolioHandler,
	jwtService service.JWTService,
) *http.Server {
	return &http.Server{
		Addr: ":8080",
		Handler: router.NewRouter(
			userHandler,
			investmentHandler,
			portfolioHandler,
			jwtService,
		),
	}
}

func gracefulShutdown(srv *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v\n", err)
	}

	log.Println("Server gracefully stopped")
}
