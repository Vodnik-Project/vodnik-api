package main

import (
	"database/sql"
	"log"

	"github.com/Vodnik-Project/vodnik-api/api"
	"github.com/Vodnik-Project/vodnik-api/db/sqlc"
	"github.com/Vodnik-Project/vodnik-api/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatalf("can't load env variables: %v", err)
	}
	db, err := sql.Open(config.DB_DRIVER, config.DB_SOURCE)
	if err != nil {
		log.Fatalf("can't connect to db: %v", err)
	}
	q := sqlc.New(db)
	server := api.NewServer(q)
	err = server.StartServer(config.SERVER_PORT)
	if err != nil {
		log.Fatalf("can't start server: %v", err)
	}
}
