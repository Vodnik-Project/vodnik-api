# Task Management API
User can create projects and define tasks in every project.

# Installation
Api uses postgresql as database. consider installing it.

Set environment variables. Ex:
```
DB_DRIVER=pgx
DB_SOURCE=postgres://USER@localhost:5432/task-mng?sslmode=disable
DB_USER=USER
DB_NAME=task-mng
SERVER_PORT=:8080
JWT_SECRET_KEY=b912f68841bbfcb2ace33db11880e343b62c85fdbcff
REFRESH_TOKEN_DURATION=2160h
ACCESS_TOKEN_DURATION=5m
```
Create database:
```
make createdb dbschema
```
Download required libraries:
```
go mod download
```
Ready to go:
```
go run main.go
```

# Technologies & Tools
- PostgreSQL
- [Echo](https://github.com/labstack/echo)
- [Zerolog](https://github.com/rs/zerolog)
- JWT
- [Sqlc](https://github.com/kyleconroy/sqlc)

# To-Do
- [ ] Dockerizing
- [ ] API Tests
- [ ] Notification System
