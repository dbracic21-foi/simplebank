package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq" // Import the PostgreSQL driver package
	"github.com/techschool/simplebank/api"
	db "github.com/techschool/simplebank/db/sqlc"
)

const (
	dbDriver     = "postgres" // Update driver name to "postgres"
	dbSource     = "postgresql://root:secret@localhost:5433/simple_bank?sslmode=disable"
	serverAdress = "0.0.0.0:8080"
)

func main() {

	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {

		log.Fatal("Cannot connect to db", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAdress)

	if err != nil {
		log.Fatal("Cannot start server", err)
	}

}
