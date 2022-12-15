// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: task.sql

package sqlc

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const createTask = `-- name: CreateTask :one

INSERT INTO tasks (
    project_id, title, info, tag, created_by, beggining, deadline, color
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING task_id, project_id, title, info, tag, created_by, created_at, beggining, deadline, color
`

type CreateTaskParams struct {
	ProjectID uuid.UUID      `json:"project_id"`
	Title     string         `json:"title"`
	Info      sql.NullString `json:"info"`
	Tag       sql.NullString `json:"tag"`
	CreatedBy uuid.UUID      `json:"created_by"`
	Beggining sql.NullTime   `json:"beggining"`
	Deadline  sql.NullTime   `json:"deadline"`
	Color     sql.NullString `json:"color"`
}

// TODO: auto time for beggining without input doesn't work.
func (q *Queries) CreateTask(ctx context.Context, arg CreateTaskParams) (Task, error) {
	row := q.queryRow(ctx, q.createTaskStmt, createTask,
		arg.ProjectID,
		arg.Title,
		arg.Info,
		arg.Tag,
		arg.CreatedBy,
		arg.Beggining,
		arg.Deadline,
		arg.Color,
	)
	var i Task
	err := row.Scan(
		&i.TaskID,
		&i.ProjectID,
		&i.Title,
		&i.Info,
		&i.Tag,
		&i.CreatedBy,
		&i.CreatedAt,
		&i.Beggining,
		&i.Deadline,
		&i.Color,
	)
	return i, err
}

const deleteTask = `-- name: DeleteTask :exec
DELETE FROM tasks
WHERE task_id = $1
`

func (q *Queries) DeleteTask(ctx context.Context, taskID uuid.UUID) error {
	_, err := q.exec(ctx, q.deleteTaskStmt, deleteTask, taskID)
	return err
}

const getTaskData = `-- name: GetTaskData :one
SELECT task_id, project_id, title, info, tag, created_by, created_at, beggining, deadline, color FROM tasks
WHERE task_id = $1
`

func (q *Queries) GetTaskData(ctx context.Context, taskID uuid.UUID) (Task, error) {
	row := q.queryRow(ctx, q.getTaskDataStmt, getTaskData, taskID)
	var i Task
	err := row.Scan(
		&i.TaskID,
		&i.ProjectID,
		&i.Title,
		&i.Info,
		&i.Tag,
		&i.CreatedBy,
		&i.CreatedAt,
		&i.Beggining,
		&i.Deadline,
		&i.Color,
	)
	return i, err
}

const getTasksInProject = `-- name: GetTasksInProject :many
SELECT task_id, project_id, title, info, tag, created_by, created_at, beggining, deadline, color FROM tasks
WHERE project_id = $1
`

func (q *Queries) GetTasksInProject(ctx context.Context, projectID uuid.UUID) ([]Task, error) {
	rows, err := q.query(ctx, q.getTasksInProjectStmt, getTasksInProject, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Task
	for rows.Next() {
		var i Task
		if err := rows.Scan(
			&i.TaskID,
			&i.ProjectID,
			&i.Title,
			&i.Info,
			&i.Tag,
			&i.CreatedBy,
			&i.CreatedAt,
			&i.Beggining,
			&i.Deadline,
			&i.Color,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateTask = `-- name: UpdateTask :one

UPDATE tasks SET
  title =     COALESCE(NULLIF($1, 'NULL'), title),
  info =      COALESCE(NULLIF($2, 'NULL'), info),
  tag =       COALESCE(NULLIF($3, 'NULL'), tag),
  beggining = COALESCE($4, beggining),
  deadline =  COALESCE($5, deadline),
  color =     COALESCE(NULLIF($6, 'NULL'), color)
WHERE task_id = $7
RETURNING task_id, project_id, title, info, tag, created_by, created_at, beggining, deadline, color
`

type UpdateTaskParams struct {
	Title     interface{}  `json:"title"`
	Info      interface{}  `json:"info"`
	Tag       interface{}  `json:"tag"`
	Beggining sql.NullTime `json:"beggining"`
	Deadline  sql.NullTime `json:"deadline"`
	Color     interface{}  `json:"color"`
	ID        uuid.UUID    `json:"id"`
}

// TODO: get tasks by filtering: title, tag, created_by, beggining, deadline
func (q *Queries) UpdateTask(ctx context.Context, arg UpdateTaskParams) (Task, error) {
	row := q.queryRow(ctx, q.updateTaskStmt, updateTask,
		arg.Title,
		arg.Info,
		arg.Tag,
		arg.Beggining,
		arg.Deadline,
		arg.Color,
		arg.ID,
	)
	var i Task
	err := row.Scan(
		&i.TaskID,
		&i.ProjectID,
		&i.Title,
		&i.Info,
		&i.Tag,
		&i.CreatedBy,
		&i.CreatedAt,
		&i.Beggining,
		&i.Deadline,
		&i.Color,
	)
	return i, err
}
