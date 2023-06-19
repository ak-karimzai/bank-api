run: server
	./server

restart: clean run

server:
	go build -o server

sqlc:
	sqlc generate

test: 
	go test -v -cover ./...

db_docs:
	dbdocs build doc/db.dbml

db_schema:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml

mock:
	mockgen --destination=internal/db/mock/store.go --package=mockdb github.com/ak-karimzai/ak-karimzai/simpleb/internal/db Store

proto:
	rm -f pb/*.pb.go
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
		--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
		proto/*.proto

evans:
	evans --host localhost --port 9090 -r repl

clean: 
	rm -f server

.PHONY: db_schema db_docs proto mock db_schema db_docs sqlc test clean run restart evans