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

type projectDataRequest struct {
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
	var project projectDataRequest
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
	err = s.store.CreateProjectTx(c, sqlc.CreateProjectParams{
		Title:   project.Title,
		Info:    sql.NullString{String: project.Info, Valid: true},
		OwnerID: uuid.NullUUID{UUID: userUUID, Valid: true},
	})
	if err != nil {
		traceid := util.RandomString(8)
		log.Logger.Err(err).Str("traceid", traceid).Msg("")
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "an error occurred while processing your request",
			"traceid": traceid,
		})
	}
	createdProjectData := c.Get("project").(sqlc.Project)
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
	projectUUID := c.Get("projectUUID").(uuid.UUID)
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
		log.Logger.Err(err).Str("traceid", traceid).Msg("")
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "an error occurred while processing your request",
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

type updateProjectRequest struct {
	Title   string `json:"title"`
	Info    string `json:"info"`
	OwnerID string `json:"owner_id"`
}

func (s Server) UpdateProject(c echo.Context) error {
	ctx := c.Request().Context()
	var updateProjectData updateProjectRequest
	err := c.Bind(&updateProjectData)
	if err != nil {
		traceid := util.RandomString(8)
		log.Logger.Err(err).Str("traceid", traceid).Msg("can't parse input data")
		return c.JSON(http.StatusUnprocessableEntity, echo.Map{
			"message": "invalid input",
			"traceid": traceid,
		})
	}
	if updateProjectData == (updateProjectRequest{}) {
		traceid := util.RandomString(8)
		log.Logger.Err(errors.New("input data is empty")).Str("traceid", traceid).Msg("")
		return c.JSON(http.StatusUnprocessableEntity, echo.Map{
			"message": "input data is empty",
			"traceid": traceid,
		})
	}
	var ownerUUID uuid.UUID
	if updateProjectData.OwnerID != "" {
		ownerUUID, err = uuid.FromString(updateProjectData.OwnerID)
		if err != nil {
			traceid := util.RandomString(8)
			log.Logger.Err(err).Str("traceid", traceid).Msg("invalid ownerID")
			return c.JSON(http.StatusUnprocessableEntity, echo.Map{
				"message": "invalid ownerID",
				"traceid": traceid,
			})
		}
		_, err = s.store.GetUserById(ctx, ownerUUID)
		if err != nil {
			traceid := util.RandomString(8)
			log.Logger.Err(err).Str("traceid", traceid).Msg("")
			return c.JSON(http.StatusNotFound, echo.Map{
				"message": "no user found to change ownership",
				"traceid": traceid,
			})
		}
	}
	updatedProject, err := s.store.UpdateProject(ctx, sqlc.UpdateProjectParams{
		Title:     updateProjectData.Title,
		Info:      updateProjectData.Info,
		OwnerID:   ownerUUID,
		ProjectID: c.Get("projectUUID").(uuid.UUID),
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
	ctx := c.Request().Context()
	err := s.store.DeleteProject(ctx, c.Get("projectUUID").(uuid.UUID))
	if err != nil {
		traceid := util.RandomString(8)
		log.Logger.Err(err).Str("traceid", traceid).Msg("")
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "an error occurred while processing your request",
			"traceid": traceid,
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "project deleted successfully",
	})
}

type usersInProjectResponse struct {
	UserID         string `json:"userid"`
	Username       string `json:"username"`
	Bio            string `json:"bio"`
	AddedToProject string `json:"added_to_project"`
	Admin          bool   `json:"admin"`
}

func (s Server) GetUsersInProject(c echo.Context) error {
	ctx := c.Request().Context()
	projectUUID := c.Get("projectUUID").(uuid.UUID)
	usersInProject, err := s.store.GetUsersByProjectID(ctx, projectUUID)
	if err != nil {
		traceid := util.RandomString(8)
		log.Logger.Err(err).Str("traceid", traceid).Msg("")
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "an error occurred while processing your request",
			"traceid": traceid,
		})
	}
	var responseData []usersInProjectResponse
	for _, i := range usersInProject {
		user, err := s.store.GetUserById(ctx, i.UserID)
		if err != nil {
			traceid := util.RandomString(8)
			log.Logger.Err(err).Str("traceid", traceid).Msg("")
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"message": "an error occurred while processing your request",
				"traceid": traceid,
			})
		}
		responseData = append(responseData, usersInProjectResponse{
			UserID:         user.UserID.String(),
			Username:       user.Username,
			Bio:            user.Bio.String,
			AddedToProject: i.AddedAt.Time.Format(time.RFC3339),
			Admin:          i.Admin.Bool,
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "users found",
		"count":   len(usersInProject),
		"users":   responseData,
	})
}

func (s Server) AddUserToProject(c echo.Context) error {
	ctx := c.Request().Context()
	userToAdd := c.Param("userid")
	userToAddUUID, err := uuid.FromString(userToAdd)
	if err != nil {
		traceid := util.RandomString(8)
		log.Logger.Err(err).Str("traceid", traceid).Msg("")
		return c.JSON(http.StatusUnprocessableEntity, echo.Map{
			"message": "invalid userid",
			"traceid": traceid,
		})
	}
	_, err = s.store.AddUserToProject(ctx, sqlc.AddUserToProjectParams{
		UserID:    userToAddUUID,
		ProjectID: c.Get("projectUUID").(uuid.UUID),
	})
	if err != nil {
		traceid := util.RandomString(8)
		log.Logger.Err(err).Str("traceid", traceid).Msg("")
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "an error occurred while processing your request",
			"traceid": traceid,
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "user added to project successfully",
	})
}

func (s Server) DeleteUserFromProject(c echo.Context) error {
	ctx := c.Request().Context()
	userToDelete := c.Param("userid")
	userToDeleteUUID, err := uuid.FromString(userToDelete)
	if err != nil {
		traceid := util.RandomString(8)
		log.Logger.Err(err).Str("traceid", traceid).Msg("")
		return c.JSON(http.StatusUnprocessableEntity, echo.Map{
			"message": "invalid userid",
			"traceid": traceid,
		})
	}
	err = s.store.DeleteUserFromProject(ctx, sqlc.DeleteUserFromProjectParams{
		UserID:    userToDeleteUUID,
		ProjectID: c.Get("projectUUID").(uuid.UUID),
	})
	if err != nil {
		traceid := util.RandomString(8)
		log.Logger.Err(err).Str("traceid", traceid).Msg("")
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "an error occurred while processing your request",
			"traceid": traceid,
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "user deleted from project successfully",
	})
}
