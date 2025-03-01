# Book Management API

A REST API for managing books with Redis caching and Kafka event streaming, built with Go, Gin, PostgreSQL, and Swagger documentation.

## Features

- CRUD operations for books
- Redis caching for GET endpoints
- Kafka event streaming for write operations
- Swagger API documentation
- PostgreSQL database
- Pagination support
- Proper error handling

## Installation

### Prerequisites
- Go 1.23
- PostgreSQL
- Redis
- Kafka
- Swag CLI (`go install github.com/swaggo/swag/cmd/swag@latest`)
- Docker (optional)

### Local Setup

1. **Clone the repository**
   ```bash
   git clone https://github.com/shani34/book-management-system.git
   cd book-management-system
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Generate Swagger documentation**
   ```bash
   swag init -g cmd/main.go -o docs
   ```

4. **Start services**
   ```bash
   docker-compose up --build  # Starts PostgreSQL, Redis, Kafka
   ```
5. **Add .env file from docker compose file**
    ```bash
   DB_HOST=localhost  # modify according to your local set up
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=passowrd
   DB_NAME=railway
   DB_SSL_MODE=disable

   REDIS_HOST=localhost
   REDIS_PORT=6379
   REDIS_PASSWORD=password
   REDIS_DB=0

   KAFKA_BROKERS=localhost:9092
   SERVER_PORT=8080
   SERVER_READ_TIMEOUT=10s
   SERVER_WRITE_TIMEOUT=10s
    ```

## Production Deployment

### Deploy on Railway

Live API URL:
[https://book-management-system-production-7d0e.up.railway.app](https://book-management-system-production-7d0e.up.railway.app)


## API Documentation

### Base URLs:
- **Local:** [http://localhost:8080/api/v1](http://localhost:8080/api/v1)
- **Production:** [https://book-management-system-production-7d0e.up.railway.app/api/v1](https://book-management-system-production-7d0e.up.railway.app/api/v1)

### Books Endpoints

#### Create Book
```http
POST https://book-management-system-production-7d0e.up.railway.app/api/v1/books
Content-Type: application/json

{
  "title": "The Go Programming Language",
  "author": "Alan Donovan",
  "year": 2015
}
```

#### Get All Books
```http
GET https://book-management-system-production-7d0e.up.railway.app/api/v1/books?limit=10&offset=0
```

#### Get Single Book
```http
GET https://book-management-system-production-7d0e.up.railway.app/api/v1/books/{id}
```

#### Update Book
```http
PUT https://book-management-system-production-7d0e.up.railway.app/api/v1/books/{id}
Content-Type: application/json

{
  "title": "Updated Title",
  "author": "Updated Author",
  "year": 2023
}
```

#### Delete Book
```http
DELETE https://book-management-system-production-7d0e.up.railway.app/api/v1/books/{id}
```

## Swagger UI
- **Local:** [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)
- **Production:** [https://book-management-system-production-7d0e.up.railway.app/swagger/index.html](https://book-management-system-production-7d0e.up.railway.app/swagger/index.html)

## Postman Collection
To run the API in Postman, import the collection file and configure your environment variables 

## Technologies Used
- **Framework:** Gin
- **Database:** PostgreSQL
- **Caching:** Redis
- **Event Streaming:** Kafka
- **Documentation:** Swagger
- **ORM:** GORM
- **Containerization:** Docker


## Contact
- **Shani Kumar** - [shani.mnnit18@gmail.com](mailto:shani.mnnit18@gmail.com)


