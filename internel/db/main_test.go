package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	dbSource = `host=localhost port=5432 user=postgres password=postgres dbname=bank-api sslmode=disable`
)

var testQueries *Store
var testDb *pgxpool.Pool

func TestMain(m *testing.M) {
	var err error
	testDb, err = pgxpool.New(context.Background(), dbSource)
	if err != nil {
		log.Fatal("cannot connect to database: ", err)
	}

	testQueries = NewStore(testDb)
	os.Exit(m.Run())
}
