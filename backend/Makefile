DATABASE_URL=mysql://kokomed:secret@tcp(172.17.0.2:3306)/kokomed?multiStatements=true&parseTime=true
MIGRATIONS_PATH=file:///home/emilio-cliff/kokomed-fin/backend/internal/mysql/migrations

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

migrateUp:
	migrate -source ${MIGRATIONS_PATH} -database "${DATABASE_URL}" -verbose up

migrateDown:
	migrate -source ${MIGRATIONS_PATH} -database "${DATABASE_URL}" -verbose down ${VERSION}

.PHONY: test race-test sqlc run coverage createMigrate mock migrateUp migrateDown
	