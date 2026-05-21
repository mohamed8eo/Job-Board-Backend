db_url := env("DB_URL", "postgres://postgres:password@localhost:5434/jobboarddb?sslmode=disable")

[group('dev')]
dev:
    air

[group('prod')]
start:
    ./jobBoard

[group('prod')]
deploy:
    go build -o jobBoard
    just start

[group('db')]
up:
    goose -dir internal/db/migrations postgres "{{ db_url }}" up

[group('db')]
graphql:
    gqlgen generate

[group('db')]
down:
    goose -dir internal/db/migrations postgres "{{ db_url }}" down

[group('db')]
status:
    goose -dir internal/db/migrations postgres "{{ db_url }}" status

[group('sqlc')]
sqlc:
    sqlc generate
