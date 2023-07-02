package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/ak-karimzai/bank-api/internel/util"
	"github.com/jackc/pgx/v5/pgxpool"
)

var testQueries Store
var testDb *pgxpool.Pool

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../../")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}
	testDb, err = pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to database: ", err)
	}

	testQueries = NewStore(testDb)
	os.Exit(m.Run())
}
