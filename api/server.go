package api

import (
	"github.com/Vodnik-Project/vodnik-api/auth"
	"github.com/Vodnik-Project/vodnik-api/db/sqlc"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	queries     sqlc.Queries
	tokenMaker  auth.TokenMaker
	tokenSecret string
	e           *echo.Echo
}

func NewServer(queries *sqlc.Queries, tokenSecret string, tokenMaker auth.TokenMaker) *Server {
	e := echo.New()
	server := &Server{
		queries:     *queries,
		tokenSecret: tokenSecret,
		tokenMaker:  tokenMaker,
	}

	e.POST("/login", server.Login)
	e.POST("/refresh_token", server.Refresh_token)

	e.POST("/user", server.CreateUser)
	e.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:     &auth.AccessTokenPayload{},
		SigningKey: []byte(tokenSecret),
		Skipper:    skipper,
	}))
	user := e.Group("/user")
	// user.Use(server.userAuthorization)
	user.GET("", server.GetUserData)
	user.PUT("", server.UpdateUser)
	user.DELETE("", server.DeleteUser)

	// task := e.Group("/task")
	// task.POST("/", server.CreateTask)
	// task.GET("/:taskid", server.GetTaskData)
	// task.GET("/byproject/:projectid", server.GetTasksByProjectID)
	// task.PUT("/:taskid", server.UpdateTask)
	// task.DELETE("/:taskid", server.DeleteTask)
	// task.GET("/user/:taskid/users", server.GetUsersInTask)
	// task.POST("/user/:taskid/:username", server.AddUserToTask)
	// task.DELETE("/user/:taskid/:username", server.DeleteUserFromTask)

	// project := e.Group("/project")
	// project.POST("/", server.CreateProject)
	// project.GET("/:projectid", server.GetProjectData)
	// project.PUT("/:projectid", server.UpdateProject)
	// project.DELETE("/:projectid", server.DeleteProject)
	// project.GET("/user/:projectid/users", server.GetUsersInProject)
	// project.POST("/user/:projectid/:username", server.AddUserToProject)
	// project.DELETE("/user/:projectid/:username", server.DeleteUserFromProject)

	server.e = e
	return server
}

func (s *Server) StartServer(addr string) error {
	return s.e.Start(addr)
}
