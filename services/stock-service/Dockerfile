# Dockerfile pour Stock-service (Go Gin + Redis)
FROM golang:1.21-alpine AS builder

LABEL maintainer="ICS GMAO Team"
LABEL service="stock-service"
LABEL technology="Go Gin"

RUN apk add --no-cache git ca-certificates tzdata
WORKDIR /app

# Copier go.mod et go.sum avant le code source
COPY go.mod go.sum ./
RUN go mod download

# Copier tout le code source
COPY . .

# Build statique pour Linux
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# --------------------
# Image finale minimale
# --------------------
FROM alpine:latest

RUN apk --no-cache add ca-certificates curl tzdata

# Création d'un utilisateur non-root
RUN adduser -D -s /bin/sh appuser

WORKDIR /app

# Copier le binaire
COPY --from=builder /app/main .

USER appuser

ENV ENVIRONMENT=production
ENV PORT=8004
ENV REDIS_URL=redis://redis-stock:6379
ENV JWT_SECRET=your-secret-key-here

EXPOSE 8004

HEALTHCHECK --interval=30s --timeout=3s --start-period=60s --retries=3 \
    CMD curl -f http://localhost:8004/health || exit 1

CMD ["./main"]
