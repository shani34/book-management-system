# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g cmd/main.go -o api/docs --parseDependency --parseInternal
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd

# Final stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/docs ./docs
COPY .env .

EXPOSE 8080

CMD ["./main"]