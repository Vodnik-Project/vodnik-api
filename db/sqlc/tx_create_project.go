package sqlc

import (
	"context"
	"database/sql"
)

func (s SQLStore) CreateProjectTx(ctx context.Context, args CreateProjectParams) error {
	err := s.execTx(ctx, func(q *Queries) error {
		project, err := s.Queries.CreateProject(ctx, args)
		if err != nil {
			return err
		}
		_, err = s.Queries.AddUserToProject(ctx, AddUserToProjectParams{
			UserID:    args.OwnerID,
			ProjectID: project.ProjectID,
			Admin:     sql.NullBool{Bool: true, Valid: true},
		})
		if err != nil {
			return err
		}
		return err
	})
	return err
}
