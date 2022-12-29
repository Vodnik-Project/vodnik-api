package api

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/Vodnik-Project/vodnik-api/auth"
	"github.com/Vodnik-Project/vodnik-api/db/sqlc"
	log "github.com/Vodnik-Project/vodnik-api/logger"
	"github.com/Vodnik-Project/vodnik-api/types"
	"github.com/Vodnik-Project/vodnik-api/util"
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt"
	"github.com/jackc/pgx"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/validator.v2"
)

func (s *Server) CreateUser(c echo.Context) error {
	ctx := c.Request().Context()
	var user types.CreateUserParams
	err := c.Bind(&user)
	if err != nil {
		traceid := util.RandomString(8)
		log.Logger.Err(err).Str("traceid", traceid).Msg("can't parse input data")
		return c.JSON(http.StatusUnprocessableEntity, echo.Map{
			"message": "invalid input",
			"traceid": traceid,
		})
	}
	if err = validator.Validate(user); err != nil {
		traceid := util.RandomString(8)
		log.Logger.Err(err).Str("traceid", traceid).Msg("invalid input data")
		return c.JSON(http.StatusUnprocessableEntity, echo.Map{
			"message": "invalid input data",
			"error":   err.Error(),
			"traceid": traceid,
		})
	}
	passHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		traceid := util.RandomString(8)
		log.Logger.Err(err).Str("traceid", traceid).Msg("can't generate password hash")
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "an error occurred while processing your request",
			"traceid": traceid,
		})
	}
	createdUser, err := s.store.CreateUser(ctx, sqlc.CreateUserParams{
		Username: user.Username,
		Email:    user.Email,
		PassHash: string(passHash),
		Bio:      sql.NullString{String: user.Bio, Valid: true},
	})
	if err != nil {
		pgerr := err.(pgx.PgError)
		traceid := util.RandomString(8)
		if pgerr.Code == "23505" {
			switch pgerr.ConstraintName {
			case "users_email_key":
				log.Logger.Err(err).Str("traceid", traceid).Msg("")
				return c.JSON(http.StatusConflict, echo.Map{
					"message": "user with same email already exist",
					"traceid": traceid,
				})
			case "users_username_key":
				log.Logger.Err(err).Str("traceid", traceid).Msg("")
				return c.JSON(http.StatusConflict, echo.Map{
					"message": "user with same username already exist",
					"traceid": traceid,
				})
			}
		}
		log.Logger.Err(err).Str("traceid", traceid).Msg("")
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "an error occurred while processing your request",
			"traceid": traceid,
		})
	}
	responseData := types.UserData{
		UserID:   createdUser.UserID.String(),
		Username: createdUser.Username,
		Email:    createdUser.Email,
		Bio:      createdUser.Bio.String,
		JoinedAt: createdUser.JoinedAt.Format(time.RFC3339),
	}
	log.Logger.Info().Msgf("user created: %+v", responseData)
	return c.JSON(http.StatusOK, echo.Map{
		"message": "user created successfully",
		"user":    responseData,
	})
}

func (s *Server) GetUserData(c echo.Context) error {
	ctx := c.Request().Context()
	userid := c.Param("userid")
	userUUID, err := uuid.FromString(userid)
	if err != nil {
		traceid := util.RandomString(8)
		log.Logger.Err(err).Str("traceid", traceid).Msg("")
		return c.JSON(http.StatusUnprocessableEntity, echo.Map{
			"message": "invalid userid",
			"traceid": traceid,
		})
	}
	userData, err := s.store.GetUserById(ctx, userUUID)
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
	responseData := types.UserData{
		UserID:   userData.UserID.String(),
		Username: userData.Username,
		Email:    userData.Email,
		Bio:      userData.Bio.String,
		JoinedAt: userData.JoinedAt.Format(time.RFC3339),
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "user found",
		"user":    responseData,
	})
}

func (s *Server) UpdateUser(c echo.Context) error {
	ctx := c.Request().Context()
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
	var updateData types.UpdateUserParams
	err = c.Bind(&updateData)
	if err != nil {
		traceid := util.RandomString(8)
		log.Logger.Err(err).Str("traceid", traceid).Msg("")
		return c.JSON(http.StatusUnprocessableEntity, echo.Map{
			"message": "invalid input data",
			"traceid": traceid,
		})
	}
	if updateData == (types.UpdateUserParams{}) {
		traceid := util.RandomString(8)
		log.Logger.Err(errors.New("input data is empty")).Str("traceid", traceid).Msg("")
		return c.JSON(http.StatusUnprocessableEntity, echo.Map{
			"message": "input data is empty",
			"traceid": traceid,
		})
	}
	passHash := ""
	if updateData.Password != "" {
		passHashByte, err := bcrypt.GenerateFromPassword([]byte(updateData.Password), bcrypt.DefaultCost)
		if err != nil {
			traceid := util.RandomString(8)
			log.Logger.Err(err).Str("traceid", traceid).Msg("can't generate password hash")
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"message": "an error occurred while processing your request",
				"traceid": traceid,
			})
		}
		passHash = string(passHashByte)
	}
	updatedUser, err := s.store.UpdateUser(ctx, sqlc.UpdateUserParams{
		Username: updateData.Username,
		Email:    updateData.Email,
		PassHash: passHash,
		Bio:      updateData.Bio,
		UserID:   userUUID,
	})
	if err != nil {
		pgerr := err.(pgx.PgError)
		traceid := util.RandomString(8)
		if pgerr.Code == "23505" {
			switch pgerr.ConstraintName {
			case "users_email_key":
				log.Logger.Err(err).Str("traceid", traceid).Msg("")
				return c.JSON(http.StatusConflict, echo.Map{
					"message": "user with same email already exist",
					"traceid": traceid,
				})
			case "users_username_key":
				log.Logger.Err(err).Str("traceid", traceid).Msg("")
				return c.JSON(http.StatusConflict, echo.Map{
					"message": "user with same username already exist",
					"traceid": traceid,
				})
			}
		}
		log.Logger.Err(err).Str("traceid", traceid).Msg("")
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "an error occurred while processing your request",
			"traceid": traceid,
		})
	}
	responseData := types.UserData{
		UserID:   updatedUser.UserID.String(),
		Username: updatedUser.Username,
		Email:    updatedUser.Email,
		Bio:      updatedUser.Bio.String,
		JoinedAt: updatedUser.JoinedAt.Format(time.RFC3339),
	}
	log.Logger.Info().Msgf("user updated successfully: %+v", responseData)
	return c.JSON(http.StatusOK, echo.Map{
		"message": "user updated successfully",
		"user":    responseData,
	})
}

func (s *Server) DeleteUser(c echo.Context) error {
	ctx := c.Request().Context()
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
	err = s.store.DeleteUser(ctx, userUUID)
	if err != nil {
		traceid := util.RandomString(8)
		log.Logger.Err(err).Str("traceid", traceid).Msg("")
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "an error occurred while processing your request",
			"traceid": traceid,
		})
	}
	log.Logger.Info().Msg("user deleted successfully")
	return c.JSON(http.StatusOK, echo.Map{
		"message": "user deleted successfully",
	})
}
