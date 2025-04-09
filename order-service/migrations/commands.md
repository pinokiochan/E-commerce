# Go Migrate Cheat Sheet

## Source
https://github.com/golang-migrate/migrate/

## Установка
```sh
# Установить go-migrate
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

## Создание новой миграции
```sh
migrate create -seq -ext=.sql -dir=./migrations name
```

### In this command:
- The -seq flag indicates that we want to use sequential numbering like 0001, 0002, ...
  for the migration files (instead of a Unix timestamp, which is the default).
- The -ext flag indicates that we want to give the migration files the extension .sql.
- The -dir flag indicates that we want to store the migration files in the ./migrations
directory (which will be created automatically if it doesn’t already exist).


## Применение миграций
```sh
migrate -path=./migrations -database "postgres://user:password@localhost:5432/dbname?sslmode=disable" up
```

## Откат последней миграции
```sh
migrate -path=./migrations -database "postgres://user:password@localhost:5432/dbname?sslmode=disable" down 1
```

## Откат всех миграций
```sh
migrate -path=./migrations -database "postgres://user:password@localhost:5432/dbname?sslmode=disable" down
```

## Просмотр статуса миграций
```sh
migrate -path=./migrations -database "postgres://user:password@localhost:5432/dbname?sslmode=disable" version
```

## Пропуск миграций без их выполнения
```sh
migrate -path=./migrations -database "postgres://user:password@localhost:5432/dbname?sslmode=disable" force <версия>
```

## Применение определенного количества миграций
```sh
migrate -path=./migrations -database "postgres://user:password@localhost:5432/dbname?sslmode=disable" up 2
```

psql --host=localhost --dbname=orderservice --username=temut