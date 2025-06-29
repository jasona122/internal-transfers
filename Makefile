include .env
export

.PHONY: build test tidy vet run migrate

build:
	go build -o bin/internal-transfers ./internal/main.go

test:
	go test -v -cover ./...

tidy:
	go mod tidy

fmt:
	go fmt ./...

vet:
	go vet ./...

lint:
	golangci-lint run ./...

run:
	go build -o bin/internal-transfers ./internal/main.go
	./bin/internal-transfers

migrate:
	docker-compose run --rm \
	  -e PGPASSWORD=${DB_PASSWORD} \
	  postgres \
	  bash -c "for f in /migrations/*.up.sql; do echo Applying \$${f}; psql -h ${DB_HOST} -U ${DB_USER} -d ${DB_NAME} -f \$${f}; done"
