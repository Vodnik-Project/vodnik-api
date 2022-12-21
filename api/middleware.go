package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/Vodnik-Project/vodnik-api/db/sqlc"
	log "github.com/Vodnik-Project/vodnik-api/logger"
	"github.com/Vodnik-Project/vodnik-api/util"
	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
)

func skipper(c echo.Context) bool {
	if c.Request().Method == "POST" {
		switch c.Path() {
		case "/user":
			return true
		case "/login":
			return true
		case "/refresh_token":
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
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"message": "an error occurred while processing your request",
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
				"message": "an error occurred while processing your request",
				"traceid": traceid,
			})
		}
		c.Set("projectUUID", projectUUID)
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
