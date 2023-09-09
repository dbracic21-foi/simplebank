package main

import (
	"database/sql"
	"log"

	"github.com/dbracic21-foi/simplebank/api"
	db "github.com/dbracic21-foi/simplebank/db/sqlc"
	"github.com/dbracic21-foi/simplebank/util"
	_ "github.com/golang/mock/mockgen/model"
	_ "github.com/lib/pq" // Import the PostgreSQL driver package
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
	server := api.NewServer(t,store)

	err = server.Start(config.ServerAdress)

	if err != nil {
		log.Fatal("Cannot start server", err)
	}

}
