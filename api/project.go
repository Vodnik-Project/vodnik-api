package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/Vodnik-Project/vodnik-api/db/sqlc"
	"github.com/Vodnik-Project/vodnik-api/util"
	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
)

type CreateProjectRequest struct {
	Title string `json:"title"`
	Info  string `json:"info"`
}

func (s Server) CreateProject(c echo.Context) error {
	ctx := c.Request().Context()
	var project CreateProjectRequest
	err := c.Bind(&project)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, "can't bind data")
	}
	err = util.CheckEmpty(project, []string{"Title", "Info"})
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	userid := util.GetFieldFromPayload(c, "UserID")
	userUUID, err := uuid.FromString(userid)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "can't parse uuid")
	}
	userID, err := s.store.GetUserById(ctx, userUUID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "can't get userid from db")
	}
	err = s.store.CreateProjectTx(ctx, sqlc.CreateProjectParams{
		Title:   project.Title,
		Info:    sql.NullString{String: project.Info, Valid: true},
		OwnerID: uuid.NullUUID{UUID: userID.UserID, Valid: true},
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "can't create project")
	}
	return c.JSON(http.StatusOK, "project created successfully")
}

type GetProjectDataRespond struct {
	Title     string `json:"title"`
	Info      string `json:"info"`
	Owner     string `json:"owner"`
	CreatedAt string `json:"created_at"`
}

func (s Server) GetProjectData(c echo.Context) error {
	ctx := c.Request().Context()
	projectID := c.Param("projectid")
	userid := util.GetFieldFromPayload(c, "UserID")
	userUUID, err := uuid.FromString(userid)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "can't parse uuid")
	}
	user, err := s.store.GetUserById(ctx, userUUID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	projectUUID, err := uuid.FromString(projectID)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, "projectid is wrong")
	}
	_, err = s.store.IsUserInProject(ctx, sqlc.IsUserInProjectParams{
		UserID:    user.UserID,
		ProjectID: projectUUID,
	})
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "can't access to project data")
	}
	project, err := s.store.GetProjectData(ctx, projectUUID)
	if err != nil {
		return c.JSON(http.StatusNotFound, "project not found")
	}

	return c.JSON(http.StatusOK, GetProjectDataRespond{
		Title:     project.Title,
		Info:      project.Info.String,
		Owner:     user.Username,
		CreatedAt: project.CreatedAt.Time.Format(time.RFC3339),
	})
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
