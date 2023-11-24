package db

import (
	"context"
	"database/sql"
	"simple_bank/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func CreateRandomUser(t *testing.T) User {
	arg := CreateUserParams{
		Username:       util.RandomUsername(),
		Name1:          util.RandomUsername(),
		Name2:          sql.NullString{String: util.RandomUsername(), Valid: true},
		Lastname1:      util.RandomUsername(),
		Lastname2:      sql.NullString{String: util.RandomUsername(), Valid: true},
		Email:          util.RandomEmail(),
		HashedPassword: "",
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.Name1, user.Name1)
	require.Equal(t, arg.Name2, user.Name2)
	require.Equal(t, arg.Lastname1, user.Lastname1)
	require.Equal(t, arg.Lastname2, user.Lastname2)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.NotZero(t, user.PasswordChangedAt)
	require.NotZero(t, user.CreatedAt)
	return user
}

func TestCreateUser(t *testing.T) {
	CreateRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := CreateRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.Name1, user2.Name1)
	require.Equal(t, user1.Name2, user2.Name2)
	require.Equal(t, user1.Lastname1, user2.Lastname1)
	require.Equal(t, user1.Lastname2, user2.Lastname2)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}

func TestUpdateUser(t *testing.T) {
	user1 := CreateRandomUser(t)
	arg := UpdateUserParams{
		Username:       user1.Username,
		Name1:          util.RandomUsername(),
		Name2:          sql.NullString{String: util.RandomUsername(), Valid: true},
		Lastname1:      util.RandomUsername(),
		Lastname2:      sql.NullString{String: util.RandomUsername(), Valid: true},
		Email:          util.RandomEmail(),
		HashedPassword: "",
	}
	user2, err := testQueries.UpdateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, arg.Username, user2.Username)
	require.Equal(t, arg.Name1, user2.Name1)
	require.Equal(t, arg.Name2, user2.Name2)
	require.Equal(t, arg.Lastname1, user2.Lastname1)
	require.Equal(t, arg.Lastname2, user2.Lastname2)
	require.Equal(t, arg.Email, user2.Email)
	require.Equal(t, arg.HashedPassword, user2.HashedPassword)
	require.WithinDuration(t, arg.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}

func TestDeleleUser(t *testing.T) {
	user1 := CreateRandomUser(t)
	err := testQueries.DeleteUser(context.Background(), user1.Username)
	require.NoError(t, err)
	user2, err := testQueries.GetUser(context.Background(), user1.Username)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, user2)
}

func TestGetUsers(t *testing.T) {
	for i := 0; i < 10; i++ {
		CreateRandomUser(t)
	}

	arg := GetUsersParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.GetUsers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 5)
	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
