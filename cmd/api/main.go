package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq" // Driver nulo do PostgreSQL para registrar no database/sql

	"github.com/tayane/payflow-engine/internal/handler"
	"github.com/tayane/payflow-engine/internal/repository"
	"github.com/tayane/payflow-engine/internal/usecase"
)

func main() {
	// Logger estruturado nativo do Go
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// 1. Configurações vindas de Variáveis de Ambiente (com fallback)
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "payflow")
	dbPass := getEnv("DB_PASSWORD", "payflow_secret")
	dbName := getEnv("DB_NAME", "payflow_db")
	serverPort := getEnv("PORT", "8080")

	// 2. Conexão com o PostgreSQL
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPass, dbName)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logger.Error("erro ao abrir conexao com postgres", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer db.Close()

	// Configuração do pool de conexões (Boas práticas de produção)
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		logger.Warn("banco de dados ainda nao acessivel, continuando inicialização...", slog.String("error", err.Error()))
	} else {
		logger.Info("conexao com postgres estabelecida com sucesso")
	}

	// 3. Injeção de Dependências (Clean Architecture)
	postgresRepo := repository.NewPostgresRepository(db)
	txUseCase := usecase.NewTransactionUseCase(postgresRepo)

	// 4. Handlers
	healthHandler := handler.NewHealthHandler()
	txHandler := handler.NewTransactionHandler(txUseCase)

	// 5. Roteamento HTTP (Nativo Go 1.22+)
	mux := http.NewServeMux()

	// Health Checks para Kubernetes/Docker
	mux.HandleFunc("GET /healthz", healthHandler.Healthz)
	mux.HandleFunc("GET /readyz", healthHandler.Readyz)

	// Endpoints da API REST
	mux.HandleFunc("POST /v1/transactions", txHandler.Create)

	// 6. Servidor HTTP
	server := &http.Server{
		Addr:         ":" + serverPort,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	// 7. Graceful Shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		logger.Info("servidor payflow-engine iniciado", slog.String("port", serverPort))
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("erro fatal no servidor http", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	<-stop
	logger.Info("desligando servidor graciosamente...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("erro durante o shutdown", slog.String("error", err.Error()))
	} else {
		logger.Info("servidor finalizado com sucesso")
	}
}

func getEnv(key, defaultValue string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultValue
}
