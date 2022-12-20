package sqlc

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
	CreateProjectTx(ctx context.Context, args CreateProjectParams) (interface{}, error)
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

func (s SQLStore) execTx(ctx context.Context, fn func(*Queries) (interface{}, error)) (interface{}, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	q := New(tx)
	n, err := fn(q)
	if err != nil {
		rberr := tx.Rollback()
		if rberr != nil {
			return nil, fmt.Errorf("rberr: %v, txerr: %v", rberr, err)
		}
		return nil, err
	}
	return n, tx.Commit()
}
