package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/raman-vhd/task-management-api/auth"
	"github.com/raman-vhd/task-management-api/db/sqlc"
	log "github.com/raman-vhd/task-management-api/logger"
	"github.com/raman-vhd/task-management-api/util"
)

func skipper(c echo.Context) bool {
	if c.Request().Method == "POST" {
		switch c.Path() {
		case "/api/user":
			return true
		case "/api/login":
			return true
		case "/api/refresh_token":
			return true
		}
	}
	return false
}

func (s Server) isProjectOwner(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		projectID := c.Param("projectid")
		projectUUID, err := uuid.FromString(projectID)
		if err != nil {
			traceid := util.RandomString(8)
			log.Logger.Err(err).Str("traceid", traceid).Msg("")
			return c.JSON(http.StatusUnprocessableEntity, echo.Map{
				"message": "invalid projectID",
				"traceid": traceid,
			})
		}
		c.Set("projectUUID", projectUUID)
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
		userid := c.Get("user").(*jwt.Token).Claims.(*auth.AccessTokenPayload).UserID
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
			log.Logger.Err(errors.New("only owner of project can modify the project")).Str("traceid", traceid).Msg("")
			return c.JSON(http.StatusForbidden, echo.Map{
				"message": "only owner of project can modify the project",
				"traceid": traceid,
			})
		}
		return next(c)
	}
}

func (s Server) isInProject(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		projectID := c.Param("projectid")
		projectUUID, err := uuid.FromString(projectID)
		if err != nil {
			traceid := util.RandomString(8)
			log.Logger.Err(err).Str("traceid", traceid).Msg("")
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"message": "invalid projectID",
				"traceid": traceid,
			})
		}
		c.Set("projectUUID", projectUUID)
		userid := c.Get("user").(*jwt.Token).Claims.(*auth.AccessTokenPayload).UserID
		userUUID, err := uuid.FromString(userid)
		if err != nil {
			traceid := util.RandomString(8)
			log.Logger.Err(err).Str("traceid", traceid).Msg("")
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"message": "an error occurred while processing your request",
				"traceid": traceid,
			})
		}
		c.Set("userUUID", userUUID)
		userInProject, err := s.store.IsUserInProject(ctx, sqlc.IsUserInProjectParams{
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
		c.Set("admin", userInProject.Admin.Bool)
		return next(c)
	}
}

func (s Server) isProjectAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if !c.Get("admin").(bool) {
			traceid := util.RandomString(8)
			log.Logger.Err(errors.New("not an admin")).Str("traceid", traceid).Msg("")
			return c.JSON(http.StatusForbidden, echo.Map{
				"message": "only admins can modify users in project",
				"traceid": traceid,
			})
		}
		return next(c)
	}
}

func (s Server) getTaskID(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		taskID := c.Param("taskid")
		taskUUID, err := uuid.FromString(taskID)
		if err != nil {
			traceid := util.RandomString(8)
			log.Logger.Err(err).Str("traceid", traceid).Msg("invalid taskid")
			return c.JSON(http.StatusUnprocessableEntity, echo.Map{
				"message": "invalid taskID",
				"traceid": traceid,
			})
		}
		c.Set("taskUUID", taskUUID)
		return next(c)
	}
}
