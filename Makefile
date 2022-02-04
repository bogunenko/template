createdb:
	docker exec -i mysql mysql -u user -ppass db < ./db/schema.sql

dropdb:
	docker exec -i mysql mysql -u user -ppass db < ./db/drop.sql

sqlc:
	sqlc generate -f db/sqlc.yaml

test:
	go test -v -cover ./...

go:
	go run main.go

build:
	docker-compose build

up:
	docker-compose up -d

down:
	docker-compose down

.PHONY: createdb dropdb sqlc test go build up down
