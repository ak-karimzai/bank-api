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

clean: 
	rm -f server

.PHONY: db_schema db_docs 