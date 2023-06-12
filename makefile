run: server
	./server

restart: clean run

server:
	go build -o server

sqlc:
	sqlc generate

test: 
	go test -v -cover ./...

mock:
	mockgen --destination=internal/db/mock/store.go --package=mockdb github.com/ak-karimzai/ak-karimzai/simpleb/internal/db Store

clean: 
	rm -f server