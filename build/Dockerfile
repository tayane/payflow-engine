# ---------------------------------------------------
# Estágio 1: Compilação do código em Go
# ---------------------------------------------------
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Instala certificados de CA caso a aplicação faça chamadas HTTPS externas
RUN apk add --no-cache ca-certificates

# Copia os arquivos de gerenciamento de dependências
COPY go.mod go.sum ./
RUN go mod download

# Copia todo o código-fonte do projeto
COPY . .

# Compila o binário estático sem dependências CGO
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o payflow-engine ./cmd/api

# ---------------------------------------------------
# Estágio 2: Imagem final minimalista
# ---------------------------------------------------
FROM alpine:latest

WORKDIR /app

# Copia os certificados e o binário gerado no estágio anterior
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/payflow-engine .

EXPOSE 8080

CMD ["./payflow-engine"]