# ---- Build stage ----
FROM golang:1.22-alpine AS builder

RUN apk add --no-cache ca-certificates git

WORKDIR /app

# Cache deps
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build CLI+server binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o postflow .

# ---- Runtime stage ----
FROM alpine:3.19

RUN apk add --no-cache ca-certificates && update-ca-certificates

WORKDIR /app

COPY --from=builder /app/postflow /app/postflow

# Expose server port (matches default :8080)
EXPOSE 8080

# Default PORT; override at runtime if needed
ENV PORT=:8080

# Entrypoint: run server via CLI subcommand
CMD ["/app/postflow", "serve"]
