package sqlc

import (
	"database/sql"

	"github.com/labstack/echo/v4"
)

func (s SQLStore) CreateProjectTx(c echo.Context, args CreateProjectParams) error {
	err := s.execTx(c, func(q *Queries) error {
		ctx := c.Request().Context()
		project, err := s.Queries.CreateProject(ctx, args)
		if err != nil {
			return err
		}
		_, err = s.Queries.AddUserToProject(ctx, AddUserToProjectParams{
			UserID:    args.OwnerID.UUID,
			ProjectID: project.ProjectID,
			Admin:     sql.NullBool{Bool: true, Valid: true},
		})
		c.Set("project", project)
		return err
	})
	return err
}
