# Instalation
api uses postgresql as database. consider installing it.

then set environment variables. Ex:
```
DB_DRIVER=pgx
DB_SOURCE=postgres://USER@localhost:5432/task-mng?sslmode=disable
DB_SOURCE_TEST=postgres://USER@localhost:5432/task-mng-test?sslmode=disable
DB_USER=USER
DB_NAME=task-mng
DB_NAME_TEST=task-mng-test
SERVER_PORT=:8080
JWT_SECRET_KEY=b912f68841bbfcb2ace33db11880e343b62c85fdbcff
REFRESH_TOKEN_DURATION=2160h
ACCESS_TOKEN_DURATION=5m
```
Then create database:
```
make createdb dbschema
```

Ready to go:
```
go run main.go
```
