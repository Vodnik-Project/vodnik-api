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
	e.Use(middleware.Recover())
	api := e.Group("/api")

	api.POST("/login", server.Login)
	api.POST("/refresh_token", server.Refresh_token)

	api.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:     &auth.AccessTokenPayload{},
		SigningKey: []byte(tokenSecret),
		Skipper:    skipper,
	}))
	user := api.Group("/user")
	user.POST("", server.CreateUser)
	user.GET("/:userid", server.GetUserData)
	user.PUT("", server.UpdateUser)
	user.DELETE("", server.DeleteUser)

	project := api.Group("/project")
	project.POST("", server.CreateProject)
	project.GET("/:projectid", server.GetProjectData, server.isInProject)
	project.PUT("/:projectid", server.UpdateProject, server.isProjectOwner)
	project.DELETE("/:projectid", server.DeleteProject, server.isProjectOwner)
	project.GET("/:projectid/users", server.GetUsersInProject, server.isInProject)
	project.POST("/:projectid/user/:userid", server.AddUserToProject, server.isInProject, server.isProjectAdmin)
	project.DELETE("/:projectid/user/:userid", server.DeleteUserFromProject, server.isInProject, server.isProjectAdmin)

	task := api.Group("/project/:projectid/task")
	task.POST("", server.CreateTask, server.isInProject, server.isProjectAdmin)
	task.GET("/:taskid", server.GetTaskData, server.isInProject)
	task.GET("/byproject/:projectid", server.GetTasksByProjectID, server.isInProject)
	task.PUT("/:taskid", server.UpdateTask, server.isInProject, server.isProjectAdmin)
	task.DELETE("/:taskid", server.DeleteTask, server.isInProject, server.isProjectAdmin)
	task.GET("/:taskid/users", server.GetUsersInTask, server.isInProject)
	task.POST("/:taskid/user/:username", server.AddUserToTask, server.isInProject, server.isProjectAdmin)
	task.DELETE("/:taskid/user/:username", server.DeleteUserFromTask, server.isInProject, server.isProjectAdmin)

	server.e = e
	return server
}

func (s *Server) StartServer(addr string) error {
	return s.e.Start(addr)
}
