# --- Сборка ---
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Кэшируем зависимости отдельно от исходников.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Статически слинкованный бинарник.
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/crm ./cmd/api

# --- Финальный образ ---
FROM alpine:3.20

RUN adduser -D -g '' appuser
WORKDIR /app

COPY --from=builder /app/crm /app/crm

USER appuser
EXPOSE 8080

ENTRYPOINT ["/app/crm"]
