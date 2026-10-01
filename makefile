# Variáveis do Projeto
APP_NAME=payflow-engine
MAIN_FILE=cmd/api/main.go
DOCKER_COMPOSE_FILE=docker-compose.yml

.PHONY: help build run test test-coverage lint clean docker-up docker-down docker-logs migrate-up migrate-down

help: ## Exibe este menu de ajuda
	@echo "Comandos disponíveis no PayFlow Engine:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-18s\033[0m %s\n", $$1, $$2}'

build: ## Compila o binário da aplicação
	@echo "--> Compilando o binário..."
	@go build -o bin/$(APP_NAME) $(MAIN_FILE)

run: ## Executa a aplicação localmente via Go
	@echo "--> Executando $(APP_NAME)..."
	@go run $(MAIN_FILE)

test: ## Executa os testes unitários
	@echo "--> Rodando testes unitários..."
	@go test -v -race ./...

test-coverage: ## Executa os testes e gera o relatório de cobertura
	@echo "--> Gerando relatório de cobertura de testes..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "--> Relatório gerado em coverage.html"

lint: ## Executa o linter Go (requer golangci-lint instalado)
	@echo "--> Rodando linter..."
	@golangci-lint run

clean: ## Remove arquivos compilados e temporários
	@echo "--> Limpando arquivos temporários..."
	@rm -rf bin/ coverage.out coverage.html

docker-up: ## Suba todos os containers do Docker Compose em segundo plano
	@echo "--> Subindo ambiente com Docker Compose..."
	@docker compose -f $(DOCKER_COMPOSE_FILE) up -d --build

docker-down: ## Para e remove todos os containers e volumes do Docker
	@echo "--> Parando ambiente Docker..."
	@docker compose -f $(DOCKER_COMPOSE_FILE) down -v

docker-logs: ## Exibe os logs dos containers em tempo real
	@docker compose -f $(DOCKER_COMPOSE_FILE) logs -f

deps: ## Baixa e organiza as dependências do módulo Go
	@echo "--> Baixando dependências..."
	@go mod tidy
	@go mod download
