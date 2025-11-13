# Start service

### Create .env
```
DB_PASSWORD=password
HTTP_PORT=8080
KAFKA_BROKERS=kafka:9092
POSTGRES_HOST=postgres
```

### And type terminal

```
docker-compose up --build
```
__OR__
```
make docker-up
```
### API Endpoints

```GET /``` - List orders
```GET /order/{order_uid}``` - Details order
```GET /api/order/{order_uid}``` - Details order in JSON

### Test
```
go test ./internal/models
go test ./internal/cache
go test ./internal/service
go test ./internal/handler
```
__OR__
```
make test
make test-verbose
make test-coverage
```

### Migrate
```
make migrate-up
make migrate-down
make migrate-status
```

### Local dev
```
make run-app
```
Local development does not use Kafka

### Project structure
```
L0
├───cmd
│   ├───app
│   ├───generator
│   └───migrate
├───html
├───internal
│   ├───cache
│   ├───config
│   ├───database
│   ├───handler
│   ├───interfaces
│   ├───kafka
│   ├───mocks
│   ├───models
│   └───service
└───schema
```

### Build
```
make build
make build-app
make build-generator
make build-migrate
```

### Docker
```
make docker-build
make docker-up
make docker-down
make docker-restart
make docker-logs
```

### Generator mocks
```
make generate-mocks
```