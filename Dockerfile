# ---- Build stage ----
FROM golang:1.22-alpine AS builder

RUN apk add --no-cache ca-certificates git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# ADDED: -ldflags="-s -w" to shrink the binary size since SQL is now inside
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o postflow .

# ---- Runtime stage ----
FROM alpine:3.19

RUN apk add --no-cache ca-certificates && update-ca-certificates

# SECURITY: Best practice - don't run as root
RUN adduser -D postuser
USER postuser

WORKDIR /app

COPY --from=builder /app/postflow /app/postflow

EXPOSE 8080
ENV PORT=:8080

# The binary now carries the migrations, so this just works!
CMD ["/app/postflow", "serve"]