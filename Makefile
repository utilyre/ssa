include .env
export

run:
	@gow -s -e go,mod,html,css,js run main.go

lint:
	@golangci-lint run ./...

test:
	@go test -v ./...

db:
	@docker start ssa_db || \
		docker run -d \
		-p ${DB_PORT}:5432 \
		-v ssa:/var/lib/postgresql/data \
		-e POSTGRES_USER=${DB_USER} \
		-e POSTGRES_PASSWORD=${DB_PASS} \
		--name ssa_db \
		postgres:15.3-alpine3.18

migrate:
	@goose -dir migrations postgres "postgresql://${DB_USER}:${DB_PASS}@${DB_HOST}:${DB_PORT}" up
