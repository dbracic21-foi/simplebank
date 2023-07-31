postgres:
	  docker run --name postgres3 -p 5433:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:15

createdb:
	docker exec -it postgres3 createdb --username=root --owner=root simple_bank

dropdb:	
	docker exec -it postgres3 dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5433/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5433/simple_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v  -cover ./...

server:
	 go run main.go

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test server	
