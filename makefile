postgres:
	docker run --name postgres12 -e POSTGRES_PASSWORD=secret -e POSTGRES_USER=root -p 5432:5432 --network app-network -d postgres:12-alpine;

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root simpleBank;

migrateup:
	migrate -path db/migration -database "postgres://root:secret@localhost:5432/simpleBank?sslmode=disable" -verbose up;

migratedown:
	migrate -path db/migration -database "postgres://root:secret@localhost:5432/simpleBank?sslmode=disable" -verbose down;

migratedown1:
	migrate -path db/migration -database "postgres://root:secret@localhost:5432/simpleBank?sslmode=disable" -verbose down 1;

migrateup1:
	migrate -path db/migration -database "postgres://root:secret@localhost:5432/simpleBank?sslmode=disable" -verbose up 1;
sqlc:
	sqlc generate;

dropdb:
	docker exec -it postgres12 dropdb simpleBank;
	
test:
	go test -v -cover ./...

server:
	go run main.go

mockdb:
	mockgen -package mockdb -destination ./db/mock/store.go github.com/Klaygogo/simplebank/db/sqlc Store 

.PHONY: postgres createdb migrateup migratedown sqlc dropdb test server mockdb migratedown1 migrateup1