package sqlc

import (
	"database/sql"
	"fmt"

	"github.com/labstack/echo/v4"
)

type Store interface {
	Querier
	CreateProjectTx(c echo.Context, args CreateProjectParams) error
	CreateTaskTx(c echo.Context, args CreateTaskParams) error
}

type SQLStore struct {
	db *sql.DB
	*Queries
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

func (s SQLStore) execTx(c echo.Context, fn func(*Queries) error) error {
	ctx := c.Request().Context()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		rberr := tx.Rollback()
		if rberr != nil {
			return fmt.Errorf("rberr: %v, txerr: %v", rberr, err)
		}
		return err
	}
	return tx.Commit()
}
