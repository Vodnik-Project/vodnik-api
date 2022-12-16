package api

import (
	"database/sql"
	"fmt"
	"net/http"

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

	createdUser, err := s.queries.CreateUser(c.Request().Context(), sqlc.CreateUserParams{
		Username: user.Username,
		Email:    user.Email,
		PassHash: user.Password,
		Bio:      sql.NullString{String: user.Bio, Valid: true},
	})
	if err != nil {
		msg := fmt.Sprintf("can't craete user: %v", err)
		return c.JSON(http.StatusForbidden, echo.Map{"msg": msg})
	}
	return c.JSON(http.StatusOK, createdUser)
}

func (s *Server) GetUserData(c echo.Context) error {
	return nil
}

func (s *Server) UpdateUser(c echo.Context) error {
	return nil
}

func (s *Server) DeleteUser(c echo.Context) error {
	return nil
}
