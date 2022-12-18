package sqlc

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
	CreateProjectTx(ctx context.Context, args CreateProjectParams) error
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

func (s SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
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
