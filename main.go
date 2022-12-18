package main

import (
	"database/sql"
	"log"

	"github.com/Vodnik-Project/vodnik-api/api"
	"github.com/Vodnik-Project/vodnik-api/auth"
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
	q := sqlc.NewStore(db)
	token := auth.NewTokenMaker(auth.Token{
		Secret:               []byte(config.JWT_SECRET_KEY),
		AccessTokenDuration:  config.ACCESS_TOKEN_DURATION,
		RefreshTokenDuration: config.REFRESH_TOKEN_DURATION,
	})
	server := api.NewServer(q, config.JWT_SECRET_KEY, token)
	err = server.StartServer(config.SERVER_PORT)
	if err != nil {
		log.Fatalf("can't start server: %v", err)
	}
}
