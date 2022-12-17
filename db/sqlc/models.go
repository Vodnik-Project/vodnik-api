// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0

package sqlc

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Project struct {
	ProjectID uuid.UUID      `json:"project_id"`
	Title     string         `json:"title"`
	Info      sql.NullString `json:"info"`
	OwnerID   uuid.UUID      `json:"owner_id"`
	CreatedAt sql.NullTime   `json:"created_at"`
}

type RefreshToken struct {
	Token       string    `json:"token"`
	UserID      uuid.UUID `json:"user_id"`
	Fingerprint string    `json:"fingerprint"`
	Device      string    `json:"device"`
}

type Task struct {
	TaskID    uuid.UUID      `json:"task_id"`
	ProjectID uuid.UUID      `json:"project_id"`
	Title     string         `json:"title"`
	Info      sql.NullString `json:"info"`
	Tag       sql.NullString `json:"tag"`
	CreatedBy uuid.UUID      `json:"created_by"`
	CreatedAt sql.NullTime   `json:"created_at"`
	Beggining sql.NullTime   `json:"beggining"`
	Deadline  sql.NullTime   `json:"deadline"`
	Color     sql.NullString `json:"color"`
}

type User struct {
	UserID       uuid.UUID      `json:"user_id"`
	Username     string         `json:"username"`
	Email        string         `json:"email"`
	PassHash     string         `json:"pass_hash"`
	JoinedAt     time.Time      `json:"joined_at"`
	Bio          sql.NullString `json:"bio"`
	ProfilePhoto sql.NullString `json:"profile_photo"`
}

type Usersetting struct {
	UserID   uuid.UUID    `json:"user_id"`
	Darkmode sql.NullBool `json:"darkmode"`
}

type Usersinproject struct {
	ProjectID uuid.UUID    `json:"project_id"`
	UserID    uuid.UUID    `json:"user_id"`
	AddedAt   sql.NullTime `json:"added_at"`
	Admin     sql.NullBool `json:"admin"`
}

type Usersintask struct {
	TaskID  uuid.UUID    `json:"task_id"`
	UserID  uuid.UUID    `json:"user_id"`
	AddedAt sql.NullTime `json:"added_at"`
	Admin   sql.NullBool `json:"admin"`
}
