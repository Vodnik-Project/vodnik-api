package sqlc

import (
	"context"
	"database/sql"
	"testing"

	"github.com/Vodnik-Project/vodnik-api/util"
	"github.com/stretchr/testify/require"
)

func createrandomuser(t *testing.T) User {
	randomuser := CreateUserParams{
		Username: util.RandomString(5),
		Email:    util.RandomEmail(),
		PassHash: util.RandomString(10),
		Bio:      sql.NullString{String: util.RandomString(30), Valid: true},
	}

	user, err := testQueries.CreateUser(context.Background(), randomuser)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, randomuser.Username, user.Username)
	require.Equal(t, randomuser.Email, user.Email)
	require.Equal(t, randomuser.PassHash, user.PassHash)
	require.Equal(t, randomuser.Bio, user.Bio)

	return user
}
func TestCreateUser(t *testing.T) {
	createrandomuser(t)
}

func TestGetUserById(t *testing.T) {
	randomuser := createrandomuser(t)

	user, err := testQueries.GetUserById(context.Background(), randomuser.UserID)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, randomuser.Username, user.Username)
	require.Equal(t, randomuser.Email, user.Email)
	require.Equal(t, randomuser.PassHash, user.PassHash)
	require.Equal(t, randomuser.Bio, user.Bio)
}

func TestGetUserByEmail(t *testing.T) {
	randomuser := createrandomuser(t)

	user, err := testQueries.GetUserByEmail(context.Background(), randomuser.Email)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, randomuser.Username, user.Username)
	require.Equal(t, randomuser.Email, user.Email)
	require.Equal(t, randomuser.PassHash, user.PassHash)
	require.Equal(t, randomuser.Bio, user.Bio)
}

func TestGetUserByUsername(t *testing.T) {
	randomuser := createrandomuser(t)

	user, err := testQueries.GetUserByUsername(context.Background(), randomuser.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, randomuser.Username, user.Username)
	require.Equal(t, randomuser.Email, user.Email)
	require.Equal(t, randomuser.PassHash, user.PassHash)
	require.Equal(t, randomuser.Bio, user.Bio)
}

func TestUpdateUser(t *testing.T) {
	randomuser := createrandomuser(t)
	newUsername := util.RandomString(5)
	user, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: newUsername,
		ID:       randomuser.UserID,
	})
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, newUsername, user.Username)
	require.Equal(t, randomuser.Email, user.Email)
	require.Equal(t, randomuser.PassHash, user.PassHash)
	require.Equal(t, randomuser.Bio, user.Bio)
	require.Equal(t, randomuser.UserID, user.UserID)
	require.Equal(t, randomuser.JoinDate, user.JoinDate)
}

func TestDeleteUser(t *testing.T) {
	randomuser := createrandomuser(t)

	err := testQueries.DeleteUser(context.Background(), randomuser.UserID)
	require.NoError(t, err)

	user, err := testQueries.GetUserById(context.Background(), randomuser.UserID)
	require.Error(t, err)
	require.EqualError(t, sql.ErrNoRows, err.Error())
	require.Empty(t, user)
}
