package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/Vodnik-Project/vodnik-api/db/sqlc"
	"github.com/Vodnik-Project/vodnik-api/util"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/labstack/echo/v4"
)

type CreateUserReqParams struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Bio      string `json:"bio"`
}

func (s *Server) CreateUser(c echo.Context) error {
	var user CreateUserReqParams
	err := c.Bind(&user)
	if err != nil {
		msg := fmt.Sprintf("can't bind received data: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"msg": msg})
	}

	err = util.CheckEmpty(user, []string{"Username", "Email", "Password"})
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, echo.Map{"err": err.Error()})
	}
	passHash := util.PassHash(user.Password)
	createdUser, err := s.queries.CreateUser(c.Request().Context(), sqlc.CreateUserParams{
		Username: user.Username,
		Email:    user.Email,
		PassHash: passHash,
		Bio:      sql.NullString{String: user.Bio, Valid: true},
	})
	if err != nil {
		msg := fmt.Sprintf("can't craete user: %v", err)
		return c.JSON(http.StatusForbidden, echo.Map{"msg": msg})
	}
	return c.JSON(http.StatusOK, createdUser)
}

type userDataRespond struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Bio      string `json:"bio"`
	JoinedAt string `json:"joinedAt"`
}

func (s *Server) GetUserData(c echo.Context) error {
	ctx := c.Request().Context()
	username := util.GetUsername(c)
	userData, err := s.queries.GetUserByUsername(ctx, username)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "can't get data from db")
	}
	return c.JSON(http.StatusOK, userDataRespond{
		Username: userData.Username,
		Email:    userData.Email,
		Bio:      userData.Bio.String,
		JoinedAt: userData.JoinedAt.Format(time.RFC3339),
	})
}

type updateUserRequest struct {
	NewUsername string `json:"username"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Bio         string `json:"bio"`
}

func (s *Server) UpdateUser(c echo.Context) error {
	ctx := c.Request().Context()
	username := util.GetUsername(c)
	var updateData updateUserRequest
	err := c.Bind(&updateData)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, "can't bind input data")
	}
	if updateData.Password != "" {
		updateData.Password = util.PassHash(updateData.Password)
	}
	_, err = s.queries.UpdateUser(ctx, sqlc.UpdateUserParams{
		Username:    username,
		NewUsername: updateData.NewUsername,
		Email:       updateData.Email,
		PassHash:    updateData.Password,
		Bio:         updateData.Bio,
	})
	if err != nil {
		msg := fmt.Errorf("can't update data to db: %v", err)
		return c.JSON(http.StatusInternalServerError, msg.Error())
	}
	return c.JSON(http.StatusOK, "user data updated succesfully")
}

func (s *Server) DeleteUser(c echo.Context) error {
	return nil
}
