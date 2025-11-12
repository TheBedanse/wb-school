# Database Migrations Tool

__Usage:__
```
  go run cmd/migrate/main.go <command>
  ```

__Commands:__
|Столбец 1|Столбец 2|
|:-:|:-:|
  up      | Apply all pending migrations (create tables)
  down    | Rollback last migration  
  reset   | Rollback all migrations
  status  | Show migration status
  version | Show current version

__Examples:__
```
  go run cmd/migrate/main.go up
  go run cmd/migrate/main.go down
  go run cmd/migrate/main.go status
  ```

## Install goose
https://pressly.github.io/goose/installation

# .env
```
DB_PASSWORD=password
HTTP_PORT=8080
KAFKA_BROKERS=kafka:9092
POSTGRES_HOST=postgres
```