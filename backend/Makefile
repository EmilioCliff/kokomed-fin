test:
	go test -v ./...

race-test:
	go test -v -race ./...

coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out 
	go tool cover -html=coverage.out -o coverage.html

sqlc:
	cd ./.envs/configs/ && sqlc generate

run:
	cd cmd/server && go run main.go

build:
	cd cmd/server && go build -o main .

mock:
	mockgen -package mockdb -destination ./internal/mysql/mockdb/mock.go github.com/EmilioCliff/kokomed-fin/backend/internal/mysql/generated Querier

createMigrate:
	migrate create -ext sql -dir ./internal/mysql/migrations/ -seq "${SCHEMA}"

migrateup:
	migrate -path ./internal/mysql/migrations -database "mysql://kokomed:secret@tcp(localhost:3306)/kokomed-fin" -verbose up

migratedown:
	migrate -path ./internal/mysql/migrations -database "mysql://kokomed:secret@tcp(localhost:3306)/kokomed-fin" -verbose down

cert:
	cd ./internal/cert && ./gen.sh

.PHONY: test race-test sqlc run coverage createMigrate mock cert
	