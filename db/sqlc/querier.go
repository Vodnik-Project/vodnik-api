// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0

package sqlc

import (
	"context"

	"github.com/gofrs/uuid"
)

type Querier interface {
	AddUserToProject(ctx context.Context, arg AddUserToProjectParams) (Usersinproject, error)
	AddUserToTask(ctx context.Context, arg AddUserToTaskParams) (Usersintask, error)
	CreateProject(ctx context.Context, arg CreateProjectParams) (Project, error)
	// TODO: auto time for beggining without input doesn't work.
	CreateTask(ctx context.Context, arg CreateTaskParams) (Task, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteProject(ctx context.Context, projectID uuid.UUID) error
	DeleteSession(ctx context.Context, token string) error
	DeleteTask(ctx context.Context, taskID uuid.UUID) error
	DeleteUser(ctx context.Context, userID uuid.UUID) error
	DeleteUserFromProject(ctx context.Context, arg DeleteUserFromProjectParams) error
	DeleteUserFromTask(ctx context.Context, arg DeleteUserFromTaskParams) error
	GetDeviceSession(ctx context.Context, arg GetDeviceSessionParams) (RefreshToken, error)
	GetProjectData(ctx context.Context, projectID uuid.UUID) (Project, error)
	GetProjectsByUserID(ctx context.Context, userID uuid.UUID) ([]Usersinproject, error)
	GetSessionByToken(ctx context.Context, token string) (RefreshToken, error)
	GetTaskData(ctx context.Context, taskID uuid.UUID) (Task, error)
	GetTasksByProjectID(ctx context.Context, projectID uuid.UUID) ([]Task, error)
	GetTasksByUserID(ctx context.Context, userID uuid.UUID) ([]Usersintask, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetUserById(ctx context.Context, userID uuid.UUID) (User, error)
	GetUsersByProjectID(ctx context.Context, projectID uuid.UUID) ([]Usersinproject, error)
	GetUsersByTaskID(ctx context.Context, taskID uuid.UUID) ([]Usersintask, error)
	IsAdmin(ctx context.Context, arg IsAdminParams) (Usersinproject, error)
	IsUserInProject(ctx context.Context, arg IsUserInProjectParams) (Usersinproject, error)
	SetSession(ctx context.Context, arg SetSessionParams) error
	UpdateProject(ctx context.Context, arg UpdateProjectParams) (Project, error)
	// TODO: get tasks by filtering: title, tag, created_by, beggining, deadline
	UpdateTask(ctx context.Context, arg UpdateTaskParams) (Task, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error)
}

var _ Querier = (*Queries)(nil)
