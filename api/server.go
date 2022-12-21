package api

import (
	"github.com/Vodnik-Project/vodnik-api/auth"
	"github.com/Vodnik-Project/vodnik-api/db/sqlc"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	store       sqlc.Store
	tokenMaker  auth.TokenMaker
	tokenSecret string
	e           *echo.Echo
}

func NewServer(store sqlc.Store, tokenSecret string, tokenMaker auth.TokenMaker) *Server {
	e := echo.New()
	server := &Server{
		store:       store,
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
	user.GET("/:userid", server.GetUserData)
	user.PUT("", server.UpdateUser)
	user.DELETE("", server.DeleteUser)

	project := e.Group("/project")
	project.POST("", server.CreateProject)
	project.GET("/:projectid", server.GetProjectData, server.isInProject)
	project.PUT("/:projectid", server.UpdateProject, server.isProjectOwner)
	project.DELETE("/:projectid", server.DeleteProject, server.isProjectOwner)
	project.GET("/:projectid/users", server.GetUsersInProject, server.isInProject)
	project.POST("/:projectid/user/:userid", server.AddUserToProject, server.isInProject)
	project.DELETE("/:projectid/user/:userid", server.DeleteUserFromProject, server.isInProject)

	// task := e.Group("/task")
	// task.POST("", server.CreateTask)
	// task.GET("/:taskid", server.GetTaskData)
	// task.GET("/byproject/:projectid", server.GetTasksByProjectID)
	// task.PUT("/:taskid", server.UpdateTask)
	// task.DELETE("/:taskid", server.DeleteTask)
	// task.GET("/:taskid/users", server.GetUsersInTask)
	// task.POST("/:taskid/user/:username", server.AddUserToTask)
	// task.DELETE("/:taskid/user/:username", server.DeleteUserFromTask)

	server.e = e
	return server
}

func (s *Server) StartServer(addr string) error {
	return s.e.Start(addr)
}
