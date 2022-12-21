package api

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/Vodnik-Project/vodnik-api/db/sqlc"
	log "github.com/Vodnik-Project/vodnik-api/logger"
	"github.com/Vodnik-Project/vodnik-api/util"
	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
)

type CreateProjectRequest struct {
	Title string `json:"title"`
	Info  string `json:"info"`
}

type projectDataResponse struct {
	ProjectID string `json:"project_id"`
	Title     string `json:"title"`
	Info      string `json:"info"`
	OwnerID   string `json:"owner_id"`
	CreatedAt string `json:"created_at"`
}

func (s Server) CreateProject(c echo.Context) error {
	ctx := c.Request().Context()
	var project CreateProjectRequest
	err := c.Bind(&project)
	if err != nil {
		traceid := util.RandomString(8)
		log.Logger.Err(err).Str("traceid", traceid).Msg("can't parse input data")
		return c.JSON(http.StatusUnprocessableEntity, echo.Map{
			"message": "invalid input",
			"traceid": traceid,
		})
	}
	err = util.CheckEmpty(project, []string{"Title"})
	if err != nil {
		traceid := util.RandomString(8)
		log.Logger.Err(err).Str("traceid", traceid).Msg(err.Error())
		return c.JSON(http.StatusUnprocessableEntity, echo.Map{
			"message": err.Error(),
			"traceid": traceid,
		})
	}
	userid := util.GetFieldFromPayload(c, "UserID")
	userUUID, err := uuid.FromString(userid)
	if err != nil {
		traceid := util.RandomString(8)
		log.Logger.Err(err).Str("traceid", traceid).Msg("")
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "an error occurred while processing your request",
			"traceid": traceid,
		})
	}
	userID, err := s.store.GetUserById(ctx, userUUID)
	if err != nil {
		traceid := util.RandomString(8)
		log.Logger.Err(err).Str("traceid", traceid).Msg("")
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "an error occurred while processing your request",
			"traceid": traceid,
		})
	}
	createdProject, err := s.store.CreateProjectTx(ctx, sqlc.CreateProjectParams{
		Title:   project.Title,
		Info:    sql.NullString{String: project.Info, Valid: true},
		OwnerID: uuid.NullUUID{UUID: userID.UserID, Valid: true},
	})
	if err != nil {
		traceid := util.RandomString(8)
		log.Logger.Err(err).Str("traceid", traceid).Msg("")
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "an error occurred while processing your request",
			"traceid": traceid,
		})
	}
	createdProjectData := createdProject.(sqlc.Project)
	responseData := projectDataResponse{
		ProjectID: createdProjectData.ProjectID.String(),
		Title:     createdProjectData.Title,
		Info:      createdProjectData.Info.String,
		OwnerID:   createdProjectData.OwnerID.UUID.String(),
		CreatedAt: createdProjectData.CreatedAt.Time.Format(time.RFC3339),
	}
	log.Logger.Info().Msgf("project created: %+v", responseData)
	return c.JSON(http.StatusOK, echo.Map{
		"message": "project created successfully",
		"project": responseData,
	})
}

func (s Server) GetProjectData(c echo.Context) error {
	ctx := c.Request().Context()
	projectID := c.Param("projectid")
	projectUUID, err := uuid.FromString(projectID)
	if err != nil {
		traceid := util.RandomString(8)
		log.Logger.Err(err).Str("traceid", traceid).Msg("")
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "an error occurred while processing your request",
			"traceid": traceid,
		})
	}
	project, err := s.store.GetProjectData(ctx, projectUUID)
	if err != nil {
		traceid := util.RandomString(8)
		if err == sql.ErrNoRows {
			log.Logger.Err(err).Str("traceid", traceid).Msg("")
			return c.JSON(http.StatusNotFound, echo.Map{
				"message": "no project found",
				"traceid": traceid,
			})
		}
	}
	userid := util.GetFieldFromPayload(c, "UserID")
	userUUID, err := uuid.FromString(userid)
	if err != nil {
		traceid := util.RandomString(8)
		log.Logger.Err(err).Str("traceid", traceid).Msg("")
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "an error occurred while processing your request",
			"traceid": traceid,
		})
	}
	_, err = s.store.IsUserInProject(ctx, sqlc.IsUserInProjectParams{
		UserID:    userUUID,
		ProjectID: projectUUID,
	})
	if err != nil {
		traceid := util.RandomString(8)
		log.Logger.Err(err).Str("traceid", traceid).Msg("no access to project")
		return c.JSON(http.StatusForbidden, echo.Map{
			"message": "no access to project",
			"traceid": traceid,
		})
	}
	responseData := projectDataResponse{
		ProjectID: project.ProjectID.String(),
		Title:     project.Title,
		Info:      project.Info.String,
		OwnerID:   project.OwnerID.UUID.String(),
		CreatedAt: project.CreatedAt.Time.Format(time.RFC3339),
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "project found",
		"project": responseData,
	})
}

type UpdateProjectRequest struct {
	Title string `json:"title"`
	Info  string `json:"info"`
}

func (s Server) UpdateProject(c echo.Context) error {
	ctx := c.Request().Context()
	projectID := c.Param("projectid")
	projectUUID, err := uuid.FromString(projectID)
	if err != nil {
		traceid := util.RandomString(8)
		log.Logger.Err(err).Str("traceid", traceid).Msg("")
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "an error occurred while processing your request",
			"traceid": traceid,
		})
	}
	project, err := s.store.GetProjectData(ctx, projectUUID)
	if err != nil {
		traceid := util.RandomString(8)
		if err == sql.ErrNoRows {
			log.Logger.Err(err).Str("traceid", traceid).Msg("")
			return c.JSON(http.StatusNotFound, echo.Map{
				"message": "no project found",
				"traceid": traceid,
			})
		}
	}
	userid := util.GetFieldFromPayload(c, "UserID")
	userUUID, err := uuid.FromString(userid)
	if err != nil {
		traceid := util.RandomString(8)
		log.Logger.Err(err).Str("traceid", traceid).Msg("")
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "an error occurred while processing your request",
			"traceid": traceid,
		})
	}
	if project.OwnerID.UUID != userUUID {
		traceid := util.RandomString(8)
		log.Logger.Err(errors.New("only owner of project can edit the project")).Str("traceid", traceid).Msg("")
		return c.JSON(http.StatusForbidden, echo.Map{
			"message": "only owner of project can edit the project",
			"traceid": traceid,
		})
	}
	var updateProjectData UpdateProjectRequest
	err = c.Bind(&updateProjectData)
	if err != nil {
		traceid := util.RandomString(8)
		log.Logger.Err(err).Str("traceid", traceid).Msg("can't parse input data")
		return c.JSON(http.StatusUnprocessableEntity, echo.Map{
			"message": "invalid input",
			"traceid": traceid,
		})
	}
	if updateProjectData == (UpdateProjectRequest{}) {
		traceid := util.RandomString(8)
		log.Logger.Err(errors.New("input data is empty")).Str("traceid", traceid).Msg("")
		return c.JSON(http.StatusUnprocessableEntity, echo.Map{
			"message": "input data is empty",
			"traceid": traceid,
		})
	}
	updatedProject, err := s.store.UpdateProject(ctx, sqlc.UpdateProjectParams{
		Title:     updateProjectData.Title,
		Info:      updateProjectData.Info,
		ProjectID: projectUUID,
	})
	if err != nil {
		traceid := util.RandomString(8)
		log.Logger.Err(err).Str("traceid", traceid).Msg("")
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "an error occurred while processing your request",
			"traceid": traceid,
		})
	}
	responseData := projectDataResponse{
		ProjectID: updatedProject.ProjectID.String(),
		Title:     updatedProject.Title,
		Info:      updatedProject.Info.String,
		OwnerID:   updatedProject.OwnerID.UUID.String(),
		CreatedAt: updatedProject.CreatedAt.Time.Format(time.RFC3339),
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "project updated successfully",
		"project": responseData,
	})
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
