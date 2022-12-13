// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0

package sqlc

import (
	"context"
	"database/sql"
	"fmt"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

func Prepare(ctx context.Context, db DBTX) (*Queries, error) {
	q := Queries{db: db}
	var err error
	if q.addUserToProjectStmt, err = db.PrepareContext(ctx, addUserToProject); err != nil {
		return nil, fmt.Errorf("error preparing query AddUserToProject: %w", err)
	}
	if q.createProjectStmt, err = db.PrepareContext(ctx, createProject); err != nil {
		return nil, fmt.Errorf("error preparing query CreateProject: %w", err)
	}
	if q.createTaskStmt, err = db.PrepareContext(ctx, createTask); err != nil {
		return nil, fmt.Errorf("error preparing query CreateTask: %w", err)
	}
	if q.createUserStmt, err = db.PrepareContext(ctx, createUser); err != nil {
		return nil, fmt.Errorf("error preparing query CreateUser: %w", err)
	}
	if q.deleteProjectStmt, err = db.PrepareContext(ctx, deleteProject); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteProject: %w", err)
	}
	if q.deleteTaskStmt, err = db.PrepareContext(ctx, deleteTask); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteTask: %w", err)
	}
	if q.deleteUserStmt, err = db.PrepareContext(ctx, deleteUser); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteUser: %w", err)
	}
	if q.getProjectStmt, err = db.PrepareContext(ctx, getProject); err != nil {
		return nil, fmt.Errorf("error preparing query GetProject: %w", err)
	}
	if q.getProjectUsersStmt, err = db.PrepareContext(ctx, getProjectUsers); err != nil {
		return nil, fmt.Errorf("error preparing query GetProjectUsers: %w", err)
	}
	if q.getProjectsStmt, err = db.PrepareContext(ctx, getProjects); err != nil {
		return nil, fmt.Errorf("error preparing query GetProjects: %w", err)
	}
	if q.getTaskStmt, err = db.PrepareContext(ctx, getTask); err != nil {
		return nil, fmt.Errorf("error preparing query GetTask: %w", err)
	}
	if q.getTasksStmt, err = db.PrepareContext(ctx, getTasks); err != nil {
		return nil, fmt.Errorf("error preparing query GetTasks: %w", err)
	}
	if q.getUserStmt, err = db.PrepareContext(ctx, getUser); err != nil {
		return nil, fmt.Errorf("error preparing query GetUser: %w", err)
	}
	if q.getUserProjectsStmt, err = db.PrepareContext(ctx, getUserProjects); err != nil {
		return nil, fmt.Errorf("error preparing query GetUserProjects: %w", err)
	}
	if q.updateProjectStmt, err = db.PrepareContext(ctx, updateProject); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateProject: %w", err)
	}
	if q.updateTaskStmt, err = db.PrepareContext(ctx, updateTask); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateTask: %w", err)
	}
	if q.updateUserStmt, err = db.PrepareContext(ctx, updateUser); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateUser: %w", err)
	}
	return &q, nil
}

func (q *Queries) Close() error {
	var err error
	if q.addUserToProjectStmt != nil {
		if cerr := q.addUserToProjectStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing addUserToProjectStmt: %w", cerr)
		}
	}
	if q.createProjectStmt != nil {
		if cerr := q.createProjectStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createProjectStmt: %w", cerr)
		}
	}
	if q.createTaskStmt != nil {
		if cerr := q.createTaskStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createTaskStmt: %w", cerr)
		}
	}
	if q.createUserStmt != nil {
		if cerr := q.createUserStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createUserStmt: %w", cerr)
		}
	}
	if q.deleteProjectStmt != nil {
		if cerr := q.deleteProjectStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteProjectStmt: %w", cerr)
		}
	}
	if q.deleteTaskStmt != nil {
		if cerr := q.deleteTaskStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteTaskStmt: %w", cerr)
		}
	}
	if q.deleteUserStmt != nil {
		if cerr := q.deleteUserStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteUserStmt: %w", cerr)
		}
	}
	if q.getProjectStmt != nil {
		if cerr := q.getProjectStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getProjectStmt: %w", cerr)
		}
	}
	if q.getProjectUsersStmt != nil {
		if cerr := q.getProjectUsersStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getProjectUsersStmt: %w", cerr)
		}
	}
	if q.getProjectsStmt != nil {
		if cerr := q.getProjectsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getProjectsStmt: %w", cerr)
		}
	}
	if q.getTaskStmt != nil {
		if cerr := q.getTaskStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getTaskStmt: %w", cerr)
		}
	}
	if q.getTasksStmt != nil {
		if cerr := q.getTasksStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getTasksStmt: %w", cerr)
		}
	}
	if q.getUserStmt != nil {
		if cerr := q.getUserStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getUserStmt: %w", cerr)
		}
	}
	if q.getUserProjectsStmt != nil {
		if cerr := q.getUserProjectsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getUserProjectsStmt: %w", cerr)
		}
	}
	if q.updateProjectStmt != nil {
		if cerr := q.updateProjectStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateProjectStmt: %w", cerr)
		}
	}
	if q.updateTaskStmt != nil {
		if cerr := q.updateTaskStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateTaskStmt: %w", cerr)
		}
	}
	if q.updateUserStmt != nil {
		if cerr := q.updateUserStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateUserStmt: %w", cerr)
		}
	}
	return err
}

func (q *Queries) exec(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (sql.Result, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).ExecContext(ctx, args...)
	case stmt != nil:
		return stmt.ExecContext(ctx, args...)
	default:
		return q.db.ExecContext(ctx, query, args...)
	}
}

func (q *Queries) query(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (*sql.Rows, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryContext(ctx, args...)
	default:
		return q.db.QueryContext(ctx, query, args...)
	}
}

func (q *Queries) queryRow(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) *sql.Row {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryRowContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryRowContext(ctx, args...)
	default:
		return q.db.QueryRowContext(ctx, query, args...)
	}
}

type Queries struct {
	db                   DBTX
	tx                   *sql.Tx
	addUserToProjectStmt *sql.Stmt
	createProjectStmt    *sql.Stmt
	createTaskStmt       *sql.Stmt
	createUserStmt       *sql.Stmt
	deleteProjectStmt    *sql.Stmt
	deleteTaskStmt       *sql.Stmt
	deleteUserStmt       *sql.Stmt
	getProjectStmt       *sql.Stmt
	getProjectUsersStmt  *sql.Stmt
	getProjectsStmt      *sql.Stmt
	getTaskStmt          *sql.Stmt
	getTasksStmt         *sql.Stmt
	getUserStmt          *sql.Stmt
	getUserProjectsStmt  *sql.Stmt
	updateProjectStmt    *sql.Stmt
	updateTaskStmt       *sql.Stmt
	updateUserStmt       *sql.Stmt
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db:                   tx,
		tx:                   tx,
		addUserToProjectStmt: q.addUserToProjectStmt,
		createProjectStmt:    q.createProjectStmt,
		createTaskStmt:       q.createTaskStmt,
		createUserStmt:       q.createUserStmt,
		deleteProjectStmt:    q.deleteProjectStmt,
		deleteTaskStmt:       q.deleteTaskStmt,
		deleteUserStmt:       q.deleteUserStmt,
		getProjectStmt:       q.getProjectStmt,
		getProjectUsersStmt:  q.getProjectUsersStmt,
		getProjectsStmt:      q.getProjectsStmt,
		getTaskStmt:          q.getTaskStmt,
		getTasksStmt:         q.getTasksStmt,
		getUserStmt:          q.getUserStmt,
		getUserProjectsStmt:  q.getUserProjectsStmt,
		updateProjectStmt:    q.updateProjectStmt,
		updateTaskStmt:       q.updateTaskStmt,
		updateUserStmt:       q.updateUserStmt,
	}
}
