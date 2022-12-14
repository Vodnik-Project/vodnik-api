package sqlc

import (
	"database/sql"
	"log"
	"testing"

	"github.com/Vodnik-Project/vodnik-api/util"
	_ "github.com/jackc/pgx/stdlib"
)

var TestDB *sql.DB
var TestQueries *Queries

func TestMain(t *testing.T) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Printf("can't load env variables: %v", err)
		return
	}
	TestDB, err = sql.Open(config.DB_DRIVER, config.DB_SOURCE_TEST)
	if err != nil {
		log.Printf("can't open database: %v", err)
		return
	}
	TestQueries = New(TestDB)
}
