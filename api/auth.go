package api

import (
	"fmt"
	"net/http"

	"github.com/Vodnik-Project/vodnik-api/util"
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
		return c.JSON(http.StatusUnprocessableEntity, "can't bind data")
	}
	err = util.CheckEmpty(reqData, []string{"Email", "Password"})
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	user, err := s.queries.GetUserByEmail(c.Request().Context(), reqData.Email)
	if err != nil {
		return c.JSON(http.StatusNotFound, "user not found")
	}
	if reqData.Password != user.PassHash {
		return c.JSON(http.StatusUnauthorized, "password is wrong")
	}
	accessToken, err := s.tokenMaker.CreateAccessToken(user.UserID)
	if err != nil {
		msg := fmt.Errorf("can't create access token: %v", err)
		return c.JSON(http.StatusInternalServerError, msg.Error())
	}

	refreshToken, err := s.tokenMaker.CreateRefreshToken()
	if err != nil {
		msg := fmt.Errorf("can't create refresh token: %v", err)
		return c.JSON(http.StatusInternalServerError, msg.Error())
	}
	return c.JSON(http.StatusOK, loginRespond{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

func (s Server) Refresh_token(c echo.Context) error {
	return nil
}
