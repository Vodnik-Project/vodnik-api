package sqlc

import (
	"database/sql"

	"github.com/labstack/echo/v4"
)

func (s SQLStore) CreateTaskTx(c echo.Context, args CreateTaskParams) error {
	err := s.execTx(c, func(q *Queries) error {
		ctx := c.Request().Context()
		task, err := s.Queries.CreateTask(ctx, args)
		if err != nil {
			return err
		}
		_, err = s.Queries.AddUserToTask(ctx, AddUserToTaskParams{
			UserID: args.CreatedBy,
			TaskID: task.TaskID,
			Admin:  sql.NullBool{Bool: true, Valid: true},
		})
		c.Set("task", task)
		return err
	})
	return err
}
