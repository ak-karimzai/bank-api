package main

import (
	"database/sql"
	"log"

	"github.com/ak-karimzai/ak-karimzai/simpleb/api"
	"github.com/ak-karimzai/ak-karimzai/simpleb/internal/db"
	"github.com/ak-karimzai/ak-karimzai/simpleb/util"

	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	if err = conn.Ping(); err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server")
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
