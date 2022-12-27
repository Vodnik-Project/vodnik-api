// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: users_in_task.sql

package sqlc

import (
	"context"
	"database/sql"

	"github.com/gofrs/uuid"
)

const addUserToTask = `-- name: AddUserToTask :one
INSERT INTO usersintask (
    user_id, task_id, admin
) VALUES (
    $1, $2, $3
) RETURNING task_id, user_id, added_at, admin
`

type AddUserToTaskParams struct {
	UserID uuid.UUID    `json:"user_id"`
	TaskID uuid.UUID    `json:"task_id"`
	Admin  sql.NullBool `json:"admin"`
}

func (q *Queries) AddUserToTask(ctx context.Context, arg AddUserToTaskParams) (Usersintask, error) {
	row := q.queryRow(ctx, q.addUserToTaskStmt, addUserToTask, arg.UserID, arg.TaskID, arg.Admin)
	var i Usersintask
	err := row.Scan(
		&i.TaskID,
		&i.UserID,
		&i.AddedAt,
		&i.Admin,
	)
	return i, err
}

const deleteUserFromTask = `-- name: DeleteUserFromTask :exec
DELETE FROM usersintask
WHERE user_id = $1 AND task_id = $2
`

type DeleteUserFromTaskParams struct {
	UserID uuid.UUID `json:"user_id"`
	TaskID uuid.UUID `json:"task_id"`
}

func (q *Queries) DeleteUserFromTask(ctx context.Context, arg DeleteUserFromTaskParams) error {
	_, err := q.exec(ctx, q.deleteUserFromTaskStmt, deleteUserFromTask, arg.UserID, arg.TaskID)
	return err
}

const getTasksByUserID = `-- name: GetTasksByUserID :many
SELECT tasks.task_id, tasks.title, tasks.info, tasks.tag, 
        tasks.created_by, tasks.created_at, tasks.beggining, 
        tasks.deadline, tasks.color
FROM tasks
INNER JOIN usersintask 
ON tasks.task_id=usersintask.task_id 
WHERE usersintask.user_id=$1 AND tasks.project_id=$2
`

type GetTasksByUserIDParams struct {
	UserID    uuid.UUID `json:"user_id"`
	ProjectID uuid.UUID `json:"project_id"`
}

type GetTasksByUserIDRow struct {
	TaskID    uuid.UUID      `json:"task_id"`
	Title     string         `json:"title"`
	Info      sql.NullString `json:"info"`
	Tag       sql.NullString `json:"tag"`
	CreatedBy uuid.UUID      `json:"created_by"`
	CreatedAt sql.NullTime   `json:"created_at"`
	Beggining sql.NullTime   `json:"beggining"`
	Deadline  sql.NullTime   `json:"deadline"`
	Color     sql.NullString `json:"color"`
}

func (q *Queries) GetTasksByUserID(ctx context.Context, arg GetTasksByUserIDParams) ([]GetTasksByUserIDRow, error) {
	rows, err := q.query(ctx, q.getTasksByUserIDStmt, getTasksByUserID, arg.UserID, arg.ProjectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetTasksByUserIDRow
	for rows.Next() {
		var i GetTasksByUserIDRow
		if err := rows.Scan(
			&i.TaskID,
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

const getUsersByTaskID = `-- name: GetUsersByTaskID :many
SELECT users.user_id, users.username, users.bio,
       usersintask.added_at
FROM users
INNER JOIN usersintask
ON users.user_id=usersintask.user_id
WHERE usersintask.task_id = $1
`

type GetUsersByTaskIDRow struct {
	UserID   uuid.UUID      `json:"user_id"`
	Username string         `json:"username"`
	Bio      sql.NullString `json:"bio"`
	AddedAt  sql.NullTime   `json:"added_at"`
}

func (q *Queries) GetUsersByTaskID(ctx context.Context, taskID uuid.UUID) ([]GetUsersByTaskIDRow, error) {
	rows, err := q.query(ctx, q.getUsersByTaskIDStmt, getUsersByTaskID, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetUsersByTaskIDRow
	for rows.Next() {
		var i GetUsersByTaskIDRow
		if err := rows.Scan(
			&i.UserID,
			&i.Username,
			&i.Bio,
			&i.AddedAt,
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
