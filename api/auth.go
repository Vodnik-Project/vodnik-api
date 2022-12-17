package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/Vodnik-Project/vodnik-api/auth"
	"github.com/Vodnik-Project/vodnik-api/db/sqlc"
	"github.com/Vodnik-Project/vodnik-api/util"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRespond struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (s Server) Login(c echo.Context) error {
	var reqData loginRequest
	err := c.Bind(&reqData)
	if err != nil {
		msg := fmt.Errorf("can't bind data: %v", err)
		return c.JSON(http.StatusUnprocessableEntity, msg.Error())
	}
	err = util.CheckEmpty(reqData, []string{"Email", "Password"})
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	user, err := s.queries.GetUserByEmail(c.Request().Context(), reqData.Email)
	if err != nil {
		return c.JSON(http.StatusNotFound, "user not found")
	}
	if err = util.CheckPassword(reqData.Password, user.PassHash); err != nil {
		return c.JSON(http.StatusUnauthorized, err.Error())
	}
	accessToken, err := s.tokenMaker.CreateAccessToken(user.Username)
	if err != nil {
		msg := fmt.Errorf("can't create access token: %v", err)
		return c.JSON(http.StatusInternalServerError, msg.Error())
	}
	sessionID := util.GetSessionID(c.Request().UserAgent(), c.Request().Header.Get("Accept-Language"))
	refreshToken, err := s.tokenMaker.CreateRefreshToken()
	if err != nil {
		msg := fmt.Errorf("can't create refresh token: %v", err)
		return c.JSON(http.StatusInternalServerError, msg.Error())
	}
	oldSession, err := s.queries.GetDeviceSession(c.Request().Context(), sqlc.GetDeviceSessionParams{
		Username:    user.Username,
		Fingerprint: sessionID,
	})
	if err != sql.ErrNoRows {
		s.queries.DeleteSession(c.Request().Context(), oldSession.Token)
	}
	err = s.queries.SetSession(c.Request().Context(), sqlc.SetSessionParams{
		Token:       refreshToken,
		Username:    user.Username,
		Fingerprint: sessionID,
		Device:      c.Request().UserAgent(),
	})
	if err != nil {
		msg := fmt.Errorf("can't save refresh token to db: %v", err)
		return c.JSON(http.StatusUnauthorized, msg.Error())
	}
	return c.JSON(http.StatusOK, loginRespond{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

type refreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
	Username     string `'json:"username"`
}

func (s Server) Refresh_token(c echo.Context) error {
	var refreshToken refreshTokenRequest
	err := c.Bind(&refreshToken)
	if err != nil {
		msg := fmt.Errorf("can't bind data: %v", err)
		return c.JSON(http.StatusUnprocessableEntity, msg.Error())
	}
	var payload auth.RefreshTokenPayload
	_, err = jwt.ParseWithClaims(refreshToken.RefreshToken, &payload, func(t *jwt.Token) (interface{}, error) {
		return []byte(s.tokenSecret), nil
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	session, err := s.queries.GetSessionByToken(c.Request().Context(), refreshToken.RefreshToken)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "Not a valid refresh token")
	}
	if refreshToken.Username != session.Username {
		return c.JSON(http.StatusUnauthorized, "Username is wrong")
	}
	sessionID := util.GetSessionID(c.Request().UserAgent(), c.Request().Header.Get("Accept-Language"))

	if sessionID != session.Fingerprint {
		return c.JSON(http.StatusUnauthorized, "Not valid user")
	}

	accessToken, err := s.tokenMaker.CreateAccessToken(refreshToken.Username)
	if err != nil {
		msg := fmt.Errorf("can't create access token: %v", err)
		return c.JSON(http.StatusInternalServerError, msg.Error())
	}
	return c.JSON(http.StatusOK, echo.Map{"access_token": accessToken})
}
