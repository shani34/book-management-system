# Use a multi-stage build
FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY . .

# Install swag and dependencies
RUN apk add --no-cache git
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Generate docs
RUN swag init -g cmd/main.go -o ./docs 

# Build application
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd

FROM alpine:latest
WORKDIR /app

# Copy built application and documentation
COPY --from=builder /app/main .
COPY --from=builder /app/docs ./docs
COPY .env .

# Ensure the directory exists before copying files
RUN mkdir -p /app/pkg/kafka

# Copy TLS files from local machine to container
COPY pkg/kafka/service.cert /app/pkg/kafka/service.cert
COPY pkg/kafka/service.key /app/pkg/kafka/service.key
COPY pkg/kafka/ca.pem /app/pkg/kafka/ca.pem

EXPOSE 8080
CMD ["./main"]