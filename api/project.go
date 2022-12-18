package api

import (
	"database/sql"
	"net/http"

	"github.com/Vodnik-Project/vodnik-api/db/sqlc"
	"github.com/Vodnik-Project/vodnik-api/util"
	"github.com/labstack/echo/v4"
)

type CreateProjectRequest struct {
	Title string `json:"title"`
	Info  string `json:"info"`
}

func (s Server) CreateProject(c echo.Context) error {
	ctx := c.Request().Context()
	username := util.GetUsername(c)
	var project CreateProjectRequest
	err := c.Bind(&project)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, "can't bind data")
	}
	err = util.CheckEmpty(project, []string{"Title", "Info"})
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	userID, err := s.store.GetUserByUsername(ctx, username)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "can't get userid from db")
	}
	err = s.store.CreateProjectTx(ctx, sqlc.CreateProjectParams{
		Title:   project.Title,
		Info:    sql.NullString{String: project.Info, Valid: true},
		OwnerID: userID.UserID,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "can't create project")
	}
	return c.JSON(http.StatusOK, "project created successfully")
}

func (s Server) GetProjectData(c echo.Context) error {
	return nil
}

func (s Server) UpdateProject(c echo.Context) error {
	return nil
}

func (s Server) DeleteProject(c echo.Context) error {
	return nil
}

func (s Server) GetUsersInProject(c echo.Context) error {
	return nil
}

func (s Server) AddUserToProject(c echo.Context) error {
	return nil
}

func (s Server) DeleteUserFromProject(c echo.Context) error {
	return nil
}
