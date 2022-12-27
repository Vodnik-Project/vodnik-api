// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: users_in_project.sql

package sqlc

import (
	"context"
	"database/sql"

	"github.com/gofrs/uuid"
)

const addUserToProject = `-- name: AddUserToProject :one
INSERT INTO usersinproject (
    user_id, project_id, admin
) VALUES (
    $1, $2, $3
) RETURNING project_id, user_id, added_at, admin
`

type AddUserToProjectParams struct {
	UserID    uuid.UUID    `json:"user_id"`
	ProjectID uuid.UUID    `json:"project_id"`
	Admin     sql.NullBool `json:"admin"`
}

func (q *Queries) AddUserToProject(ctx context.Context, arg AddUserToProjectParams) (Usersinproject, error) {
	row := q.queryRow(ctx, q.addUserToProjectStmt, addUserToProject, arg.UserID, arg.ProjectID, arg.Admin)
	var i Usersinproject
	err := row.Scan(
		&i.ProjectID,
		&i.UserID,
		&i.AddedAt,
		&i.Admin,
	)
	return i, err
}

const deleteUserFromProject = `-- name: DeleteUserFromProject :exec
DELETE FROM usersinproject
WHERE user_id = $1 AND project_id = $2
`

type DeleteUserFromProjectParams struct {
	UserID    uuid.UUID `json:"user_id"`
	ProjectID uuid.UUID `json:"project_id"`
}

func (q *Queries) DeleteUserFromProject(ctx context.Context, arg DeleteUserFromProjectParams) error {
	_, err := q.exec(ctx, q.deleteUserFromProjectStmt, deleteUserFromProject, arg.UserID, arg.ProjectID)
	return err
}

const getProjectsByUserID = `-- name: GetProjectsByUserID :many
SELECT projects.project_id, projects.title, projects.info, projects.owner_id, projects.created_at, 
       usersinproject.admin 
FROM projects
INNER JOIN usersinproject 
ON projects.project_id=usersinproject.project_id 
WHERE usersinproject.user_id=$1
`

type GetProjectsByUserIDRow struct {
	ProjectID uuid.UUID      `json:"project_id"`
	Title     string         `json:"title"`
	Info      sql.NullString `json:"info"`
	OwnerID   uuid.NullUUID  `json:"owner_id"`
	CreatedAt sql.NullTime   `json:"created_at"`
	Admin     sql.NullBool   `json:"admin"`
}

func (q *Queries) GetProjectsByUserID(ctx context.Context, userID uuid.UUID) ([]GetProjectsByUserIDRow, error) {
	rows, err := q.query(ctx, q.getProjectsByUserIDStmt, getProjectsByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetProjectsByUserIDRow
	for rows.Next() {
		var i GetProjectsByUserIDRow
		if err := rows.Scan(
			&i.ProjectID,
			&i.Title,
			&i.Info,
			&i.OwnerID,
			&i.CreatedAt,
			&i.Admin,
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

const getUsersByProjectID = `-- name: GetUsersByProjectID :many
SELECT users.user_id, users.username, users.bio,
       usersinproject.added_at, usersinproject.admin 
FROM users
INNER JOIN usersinproject
ON users.user_id=usersinproject.user_id
WHERE usersinproject.project_id = $1
`

type GetUsersByProjectIDRow struct {
	UserID   uuid.UUID      `json:"user_id"`
	Username string         `json:"username"`
	Bio      sql.NullString `json:"bio"`
	AddedAt  sql.NullTime   `json:"added_at"`
	Admin    sql.NullBool   `json:"admin"`
}

func (q *Queries) GetUsersByProjectID(ctx context.Context, projectID uuid.UUID) ([]GetUsersByProjectIDRow, error) {
	rows, err := q.query(ctx, q.getUsersByProjectIDStmt, getUsersByProjectID, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetUsersByProjectIDRow
	for rows.Next() {
		var i GetUsersByProjectIDRow
		if err := rows.Scan(
			&i.UserID,
			&i.Username,
			&i.Bio,
			&i.AddedAt,
			&i.Admin,
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

const isUserInProject = `-- name: IsUserInProject :one
SELECT project_id, user_id, added_at, admin FROM usersinproject
WHERE user_id = $1 AND project_id = $2
`

type IsUserInProjectParams struct {
	UserID    uuid.UUID `json:"user_id"`
	ProjectID uuid.UUID `json:"project_id"`
}

func (q *Queries) IsUserInProject(ctx context.Context, arg IsUserInProjectParams) (Usersinproject, error) {
	row := q.queryRow(ctx, q.isUserInProjectStmt, isUserInProject, arg.UserID, arg.ProjectID)
	var i Usersinproject
	err := row.Scan(
		&i.ProjectID,
		&i.UserID,
		&i.AddedAt,
		&i.Admin,
	)
	return i, err
}
