postgres:
		docker run --name postgres12 -p 7999:5432 -e POSTGRES_PASSWORD=secret -e POSTGRES_USER=root -d postgres:12-alpine

createdb:
		docker exec -it postgres12 createdb --username=root --username=root simple_bank

dropdb:
		docker exec -it postgres12 dropdb simple_bank

migrateup:
		migrate -path db/migration -database "postgresql://root:secret@localhost:7999/simple_bank?sslmode=disable" -verbose up
migrateup1:
		migrate -path db/migration -database "postgresql://root:secret@localhost:7999/simple_bank?sslmode=disable" -verbose up 1

migratedown:
		migrate -path db/migration -database "postgresql://root:secret@localhost:7999/simple_bank?sslmode=disable" -verbose down
migratedown1:
		migrate -path db/migration -database "postgresql://root:secret@localhost:7999/simple_bank?sslmode=disable" -verbose down 1

sqlc:
		sqlc generate

test:
		go test -v -cover ./...

server: 
		go run main.go

mockdb:
		mockgen -package mockdb -destination db/mock/store.go simple_bank/db/sqlc Store

proto:
		rm -f pb/*.go
		protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
    proto/*.proto

evans:
		~/evans --host localhost --port 9090 --package pb -r repl

.PHONY: postgres createdb dropdb migrateup migrateup1 migratedown migratedown1 sqlc test server mockdb proto evans