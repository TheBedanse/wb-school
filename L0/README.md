# How to start

## Create .env
```
DB_PASSWORD=password
HTTP_PORT=8080
KAFKA_BROKERS=kafka:9092
POSTGRES_HOST=postgres
```

## And type terminal

```
docker-compose up --build
```

## API

 /

## Test
go test ./internal/models
go test ./internal/cache
go test ./internal/service
go test ./internal/handler