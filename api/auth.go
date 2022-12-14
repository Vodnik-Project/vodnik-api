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
	"github.com/raman-vhd/task-management-api/types"
	"github.com/raman-vhd/task-management-api/util"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/validator.v2"
)

func (s Server) Login(c echo.Context) error {
	ctx := c.Request().Context()
	var reqData types.LoginParams
	err := c.Bind(&reqData)
	if err != nil {
		traceid := util.RandomString(8)
		log.Logger.Err(err).Str("traceid", traceid).Msg("can't parse input data")
		return c.JSON(http.StatusUnprocessableEntity, echo.Map{
			"message": "invalid input",
			"traceid": traceid,
		})
	}
	if err = validator.Validate(reqData); err != nil {
		traceid := util.RandomString(8)
		log.Logger.Err(err).Str("traceid", traceid).Msg("invalid input data")
		return c.JSON(http.StatusUnprocessableEntity, echo.Map{
			"message": "invalid input data",
			"error":   err.Error(),
			"traceid": traceid,
		})
	}
	user, err := s.store.GetUserByEmail(ctx, reqData.Email)
	if err != nil {
		traceid := util.RandomString(8)
		if err == sql.ErrNoRows {
			log.Logger.Err(err).Str("traceid", traceid).Msg("")
			return c.JSON(http.StatusNotFound, echo.Map{
				"message": "no user found",
				"traceid": traceid,
			})
		}
		log.Logger.Err(err).Str("traceid", traceid).Msg("")
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "an error occurred while processing your request",
			"traceid": traceid,
		})
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.PassHash), []byte(reqData.Password)); err != nil {
		traceid := util.RandomString(8)
		log.Logger.Err(err).Str("traceid", traceid).Msg("")
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"message": "wrong password",
			"traceid": traceid,
		})
	}
	accessToken, err := s.tokenMaker.CreateAccessToken(user.UserID.String(), user.Username)
	if err != nil {
		traceid := util.RandomString(8)
		log.Logger.Err(err).Str("traceid", traceid).Msg("")
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "an error occurred while processing your request",
			"traceid": traceid,
		})
	}
	sessionID := util.GetSessionID(c.Request().UserAgent(), c.Request().Header.Get("Accept-Language"))
	refreshToken, err := s.tokenMaker.CreateRefreshToken()
	if err != nil {
		traceid := util.RandomString(8)
		log.Logger.Err(err).Str("traceid", traceid).Msg("")
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "an error occurred while processing your request",
			"traceid": traceid,
		})
	}
	oldSession, err := s.store.GetDeviceSession(ctx, sqlc.GetDeviceSessionParams{
		UserID:      user.UserID,
		Fingerprint: sessionID,
	})

	switch err {
	case sql.ErrNoRows:
		break
	case nil:
		err = s.store.DeleteSession(ctx, oldSession.Token)
		if err != nil {
			traceid := util.RandomString(8)
			log.Logger.Err(err).Str("traceid", traceid).Msg("")
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"message": "an error occurred while processing your request",
				"traceid": traceid,
			})
		}
	default:
		traceid := util.RandomString(8)
		log.Logger.Err(err).Str("traceid", traceid).Msg("")
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "an error occurred while processing your request",
			"traceid": traceid,
		})
	}
	err = s.store.SetSession(ctx, sqlc.SetSessionParams{
		Token:       refreshToken,
		UserID:      user.UserID,
		Fingerprint: sessionID,
		Device:      c.Request().UserAgent(),
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
		"access_token":  accessToken,
		"refresh_Token": refreshToken,
	})
}

func (s Server) Refresh_token(c echo.Context) error {
	ctx := c.Request().Context()
	var refreshToken types.RefreshTokenParams
	err := c.Bind(&refreshToken)
	if err != nil {
		traceid := util.RandomString(8)
		log.Logger.Err(err).Str("traceid", traceid).Msg("can't parse input data")
		return c.JSON(http.StatusUnprocessableEntity, echo.Map{
			"message": "invalid input",
			"traceid": traceid,
		})
	}
	var payload auth.RefreshTokenPayload
	_, err = jwt.ParseWithClaims(refreshToken.RefreshToken, &payload, func(t *jwt.Token) (interface{}, error) {
		return []byte(s.tokenSecret), nil
	})
	if err != nil {
		traceid := util.RandomString(8)
		log.Logger.Err(err).Str("traceid", traceid).Msg("invalid refresh token")
		return c.JSON(http.StatusUnprocessableEntity, echo.Map{
			"message": "invalid refresh token",
			"traceid": traceid,
		})
	}
	session, err := s.store.GetSessionByToken(ctx, refreshToken.RefreshToken)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			traceid := util.RandomString(8)
			log.Logger.Err(err).Str("traceid", traceid).Msg("")
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"message": "invalid refresh token",
				"traceid": traceid,
			})
		default:
			traceid := util.RandomString(8)
			log.Logger.Err(err).Str("traceid", traceid).Msg("")
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"message": "an error occurred while processing your request",
				"traceid": traceid,
			})
		}
	}
	userUUID, err := uuid.FromString(refreshToken.UserID)
	if err != nil {
		traceid := util.RandomString(8)
		log.Logger.Err(err).Str("traceid", traceid).Msg("can't parse userid from refresh token")
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "an error occurred while processing your request",
			"traceid": traceid,
		})
	}
	if session.UserID != userUUID {
		traceid := util.RandomString(8)
		log.Logger.Err(errors.New("session userid doesn't match")).Str("traceid", traceid).Msg("")
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"message": "wrong userid",
			"traceid": traceid,
		})
	}
	sessionID := util.GetSessionID(c.Request().UserAgent(), c.Request().Header.Get("Accept-Language"))
	if sessionID != session.Fingerprint {
		traceid := util.RandomString(8)
		log.Logger.Err(errors.New("unauthorized client")).Str("traceid", traceid).Msg("")
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"message": "unauthorized client",
			"traceid": traceid,
		})
	}
	user, err := s.store.GetUserById(ctx, userUUID)
	if err != nil {
		traceid := util.RandomString(8)
		log.Logger.Err(err).Str("traceid", traceid).Msg("")
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "an error occurred while processing your request",
			"traceid": traceid,
		})
	}
	accessToken, err := s.tokenMaker.CreateAccessToken(userUUID.String(), user.Username)
	if err != nil {
		traceid := util.RandomString(8)
		log.Logger.Err(err).Str("traceid", traceid).Msg("")
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "an error occurred while processing your request",
			"traceid": traceid,
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"access_token": accessToken,
	})
}
