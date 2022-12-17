// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: user.sql

package sqlc

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (
    username, email, pass_hash, bio
) VALUES (
  $1, $2, $3, $4
)
RETURNING user_id, username, email, pass_hash, joined_at, bio, profile_photo
`

type CreateUserParams struct {
	Username string         `json:"username"`
	Email    string         `json:"email"`
	PassHash string         `json:"pass_hash"`
	Bio      sql.NullString `json:"bio"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.queryRow(ctx, q.createUserStmt, createUser,
		arg.Username,
		arg.Email,
		arg.PassHash,
		arg.Bio,
	)
	var i User
	err := row.Scan(
		&i.UserID,
		&i.Username,
		&i.Email,
		&i.PassHash,
		&i.JoinedAt,
		&i.Bio,
		&i.ProfilePhoto,
	)
	return i, err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users
WHERE user_id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	_, err := q.exec(ctx, q.deleteUserStmt, deleteUser, userID)
	return err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT user_id, username, email, pass_hash, joined_at, bio, profile_photo FROM users
WHERE email = $1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.queryRow(ctx, q.getUserByEmailStmt, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.UserID,
		&i.Username,
		&i.Email,
		&i.PassHash,
		&i.JoinedAt,
		&i.Bio,
		&i.ProfilePhoto,
	)
	return i, err
}

const getUserById = `-- name: GetUserById :one
SELECT user_id, username, email, pass_hash, joined_at, bio, profile_photo FROM users
WHERE user_id = $1
`

func (q *Queries) GetUserById(ctx context.Context, userID uuid.UUID) (User, error) {
	row := q.queryRow(ctx, q.getUserByIdStmt, getUserById, userID)
	var i User
	err := row.Scan(
		&i.UserID,
		&i.Username,
		&i.Email,
		&i.PassHash,
		&i.JoinedAt,
		&i.Bio,
		&i.ProfilePhoto,
	)
	return i, err
}

const getUserByUsername = `-- name: GetUserByUsername :one
SELECT user_id, username, email, pass_hash, joined_at, bio, profile_photo FROM users
WHERE username = $1
`

func (q *Queries) GetUserByUsername(ctx context.Context, username string) (User, error) {
	row := q.queryRow(ctx, q.getUserByUsernameStmt, getUserByUsername, username)
	var i User
	err := row.Scan(
		&i.UserID,
		&i.Username,
		&i.Email,
		&i.PassHash,
		&i.JoinedAt,
		&i.Bio,
		&i.ProfilePhoto,
	)
	return i, err
}

const updateUser = `-- name: UpdateUser :one
UPDATE users SET
  username = COALESCE(NULLIF($1, ''), username),
  email = COALESCE(NULLIF($2, ''), email),
  pass_hash = COALESCE(NULLIF($3, ''), pass_hash),
  bio = COALESCE(NULLIF($4, ''), bio)
WHERE username = $5
RETURNING user_id, username, email, pass_hash, joined_at, bio, profile_photo
`

type UpdateUserParams struct {
	NewUsername interface{} `json:"new_username"`
	Email       interface{} `json:"email"`
	PassHash    interface{} `json:"pass_hash"`
	Bio         interface{} `json:"bio"`
	Username    string      `json:"username"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	row := q.queryRow(ctx, q.updateUserStmt, updateUser,
		arg.NewUsername,
		arg.Email,
		arg.PassHash,
		arg.Bio,
		arg.Username,
	)
	var i User
	err := row.Scan(
		&i.UserID,
		&i.Username,
		&i.Email,
		&i.PassHash,
		&i.JoinedAt,
		&i.Bio,
		&i.ProfilePhoto,
	)
	return i, err
}
