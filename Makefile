build:
	@go build -o bin/gender_reveal_service ./app.go

test:
	@go test -v ./...

run: build
	@./bin/gender_reveal_service

migration:
	@migrate create -ext sql -dir cmd/model/migrate/migrations -seq $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
	@go run cmd/model/migrate/main.go up

migrate-down:
	@go run cmd/model/migrate/main.go down

migrate-force:
	@go run cmd/model/migrate/main.go force