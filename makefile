username := postgres
password := postgres
host := localhost
port := 5432
database := bank-api
db_string := "host=${host} \
							port=${port} \
							user=${username} \
							password=${password} \
							dbname=${database} \
							sslmode=disable"
db_container_name := bank-api-db
db_doc_filename := diagram.dbml
FileName ?= none

start_container:
	sudo docker start ${db_container_name}

stop_container:
	sudo docker start ${db_container_name}

postgres:
	sudo docker run --name ${db_container_name} --network=bank-network -p 5432:5432 -e POSTGRES_USER=${username} -e POSTGRES_PASSWORD=${password} -e POSTGRES_DB=${database} -d postgres:14-alpine

create_db:
	sudo docker exec -it ${db_container_name} \
		createdb --username=${username} --owner=${username} ${database}

drop_db:
	sudo docker exec -it ${db_container_name} \
		dropdb --username=${username} ${database}

migrate_up:
	cd db/migration && goose postgres ${db_string} up

migrate_down:
	cd db/migration && goose postgres ${db_string} down

create_migrate_file:
	cd db/migration && goose postgres ./ create ${FileName} sql

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run cmd/web/main.go

db_docs:
	dbdocs build docs/${db_doc_filename}

db_schema:
	dbml2sql --postgres -o docs/schema.sql docs/${db_doc_filename}

mock:
	mockgen -package mockdb -destination=internel/db/mock/store.go github.com/ak-karimzai/bank-api/internel/db Store

.PHONY: postgres createdb dropdb migrate_up migrate_down create_migrate_file sqlc test server