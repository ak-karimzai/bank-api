package main

import (
	"context"
	"log"

	"github.com/ak-karimzai/bank-api/internel/db"
	"github.com/ak-karimzai/bank-api/internel/server"
	"github.com/ak-karimzai/bank-api/internel/util"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load the configurations: ", err)
	}
	dbConn, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	err = dbConn.Ping(context.Background())
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	store := db.NewStore(dbConn)
	httpServer, err := server.NewServer(config, store)
	if err != nil {
		log.Fatal("error while creating server: ", err)
	}

	err = httpServer.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
