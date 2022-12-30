package sqlc

import (
	"database/sql"
	"log"
	"testing"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/raman-vhd/task-management-api/util"
)

var testDB *sql.DB
var testQueries *Queries

func TestMain(t *testing.T) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Printf("can't load env variables: %v", err)
		return
	}
	testDB, err = sql.Open(config.DB_DRIVER, config.DB_SOURCE_TEST)
	if err != nil {
		log.Printf("can't open database: %v", err)
		return
	}
	testQueries = New(testDB)
}
