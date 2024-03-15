package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq" // Import the PostgreSQL driver package
)

const (
	dbDriver = "postgres" // Update driver name to "postgres"
	dbSource = "postgresql://root:secret@localhost:5433/simple_bank?sslmode=disable"
)

var testStore Store

func TestMain(m *testing.M) {
	var err error
	conpool, err := pgxpool.New(context.Background(), dbSource)
	if err != nil {

		log.Fatal("Cannot connect to db", err)
	}

	testStore = NewStore(conpool)

	os.Exit(m.Run())
}
