package api

import (
	"github.com/Vodnik-Project/vodnik-api/db/sqlc"
	"github.com/labstack/echo/v4"
)

type Server struct {
	queries sqlc.Queries
	e       *echo.Echo
}

func NewServer(queries *sqlc.Queries) *Server {
	e := echo.New()
	server := &Server{
		queries: *queries,
	}

	// e.POST("/login", auth.Login)
	// e.POST("/refresh_token", auth.Refresh_token)

	e.POST("/user", server.CreateUser)
	e.GET("/user", server.GetUserData)
	e.PUT("/user/:userid", server.UpdateUser)
	e.DELETE("/user/:userid", server.DeleteUser)

	// e.POST("/task", server.CreateTask)
	// e.GET("/task/:taskid", server.GetTaskData)
	// e.GET("/tasks/:projectid", server.GetTasksInProject)
	// e.PUT("/task/:taskid", server.UpdateTask)
	// e.DELETE("/task/:taskid", server.DeleteTask)

	// e.POST("/project", server.CreateProject)
	// e.GET("/project/:projectid", server.GetProjectData)
	// e.GET("/projects/:ownerid", server.GetProjectsByUserId)
	// e.PUT("/project/:projectid", server.UpdateProject)
	// e.DELETE("/project/:projectid", server.DeleteProject)

	// e.GET("/project/:projectid/users", server.GetUsersInProject)
	// e.POST("project/:projectid/:userid", server.AddUserToProject)
	// e.DELETE("project/:projectid/:userid", server.DeleteUserFromProject)

	// e.GET("/task/:taskid/users", server.GetUsersInTask)
	// e.POST("task/:taskid/:userid", server.AddUserToTask)
	// e.DELETE("task/:taskid/:userid", server.DeleteUserFromTask)

	server.e = e
	return server
}

func (s *Server) StartServer(addr string) error {
	return s.e.Start(addr)
}
