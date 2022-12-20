package sqlc

import (
	"context"
	"database/sql"
)

func (s SQLStore) CreateProjectTx(ctx context.Context, args CreateProjectParams) (interface{}, error) {
	project, err := s.execTx(ctx, func(q *Queries) (interface{}, error) {
		project, err := s.Queries.CreateProject(ctx, args)
		if err != nil {
			return nil, err
		}
		_, err = s.Queries.AddUserToProject(ctx, AddUserToProjectParams{
			UserID:    args.OwnerID.UUID,
			ProjectID: project.ProjectID,
			Admin:     sql.NullBool{Bool: true, Valid: true},
		})
		if err != nil {
			return nil, err
		}
		return project, err
	})
	return project, err
}
