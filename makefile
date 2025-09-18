postgres:
	docker run --name postgres12 -e POSTGRES_PASSWORD=secret -e POSTGRES_USER=root -p 5432:5432 -d postgres:12-alpine;

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root simpleBank;

migrateup:
	migrate -path db/migration -database "postgres://root:secret@localhost:5432/simpleBank?sslmode=disable" -verbose up;

migratedown:
	migrate -path db/migration -database "postgres://root:secret@localhost:5432/simpleBank?sslmode=disable" -verbose down;

sqlc:
	sqlc generate;

dropdb:
	docker exec -it postgres12 dropdb simpleBank;
	
test:
	go test -v -cover ./...
