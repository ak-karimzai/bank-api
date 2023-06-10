username := postgres
password := postgres
host := localhost
port := 5432
dbname := test
docker-container := psql-dev

createdb:
	sudo docker exec -it ${docker-container} createdb --username=${username}  --owner=${username} ${dbname}

dropdb:
	sudo docker exec -it ${docker-container} dropdb --username=${username} --if-exists ${dbname}

migrateup:
	goose postgres  postgres://${username}:${password}@${host}:${port}/${dbname} up

migratedown:
	goose postgres postgres://${username}:${password}@${host}:${port}/${dbname} down

.PHONY: createdb 