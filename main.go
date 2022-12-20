package main

import (
	"database/sql"

	"github.com/Vodnik-Project/vodnik-api/api"
	"github.com/Vodnik-Project/vodnik-api/auth"
	"github.com/Vodnik-Project/vodnik-api/db/sqlc"
	log "github.com/Vodnik-Project/vodnik-api/logger"
	"github.com/Vodnik-Project/vodnik-api/util"
)

func main() {
	log.LoggerInit("./log.log")
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Logger.Fatal().
			Err(err).
			Msg("can't load env varialbes")
	}
	db, err := sql.Open(config.DB_DRIVER, config.DB_SOURCE)
	if err != nil {
		log.Logger.Fatal().
			Err(err).
			Msg("can't connect to db")
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
		log.Logger.Fatal().
			Err(err).
			Msg("can't start server")
	}
}
