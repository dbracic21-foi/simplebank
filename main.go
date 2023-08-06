package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq" // Import the PostgreSQL driver package
	"github.com/techschool/simplebank/api"
	db "github.com/techschool/simplebank/db/sqlc"
	"github.com/techschool/simplebank/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load config:", err)

	}
	conn, err := sql.Open(config.DBDRIVER, config.DBSOURCE)
	if err != nil {

		log.Fatal("Cannot connect to db", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAdress)

	if err != nil {
		log.Fatal("Cannot start server", err)
	}

}
