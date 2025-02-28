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
COPY --from=builder /app/main .
COPY --from=builder /app/docs ./docs
COPY .env .

EXPOSE 8080
CMD ["./main"]