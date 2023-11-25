package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	mockdb "simple_bank/db/mock"
	db "simple_bank/db/sqlc"
	"simple_bank/util"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type eqCreateUserParamsMatcher struct {
	arg      db.CreateUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, isOk := x.(db.CreateUserParams)
	if !isOk {
		return false
	}
	err := util.CheckPasswordHash(e.password, arg.HashedPassword)
	if err != nil {
		return false

	}
	e.arg.HashedPassword = arg.HashedPassword
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqCreateUserParams(arg db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg, password}
}

func TestCreateUserAPI(t *testing.T) {
	user, password := randomUser(t)

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"username":  user.Username,
				"name1":     user.Name1,
				"name2":     util.SqlNullStringToStringPtr(user.Name2),
				"lastname1": user.Lastname1,
				"lastname2": util.SqlNullStringToStringPtr(user.Lastname2),
				"email":     user.Email,
				"password":  password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Username:  user.Username,
					Name1:     user.Name1,
					Name2:     user.Name2,
					Lastname1: user.Lastname1,
					Lastname2: user.Lastname2,
					Email:     user.Email,
				}
				store.EXPECT().
					CreateUser(gomock.Any(), EqCreateUserParams(arg, password)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, user)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"username":  user.Username,
				"name1":     user.Name1,
				"name2":     util.SqlNullStringToStringPtr(user.Name2),
				"lastname1": user.Lastname1,
				"lastname2": util.SqlNullStringToStringPtr(user.Lastname2),
				"email":     user.Email,
				"password":  password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "DuplicateUsername",
			body: gin.H{
				"username":  user.Username,
				"name1":     user.Name1,
				"name2":     util.SqlNullStringToStringPtr(user.Name2),
				"lastname1": user.Lastname1,
				"lastname2": util.SqlNullStringToStringPtr(user.Lastname2),
				"email":     user.Email,
				"password":  password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, &pq.Error{Code: "23505"})
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name: "InvalidUsername",
			body: gin.H{
				"username":  "i#1",
				"name1":     user.Name1,
				"name2":     util.SqlNullStringToStringPtr(user.Name2),
				"lastname1": user.Lastname1,
				"lastname2": util.SqlNullStringToStringPtr(user.Lastname2),
				"email":     user.Email,
				"password":  password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidEmail",
			body: gin.H{
				"username":  user.Username,
				"name1":     user.Name1,
				"name2":     util.SqlNullStringToStringPtr(user.Name2),
				"lastname1": user.Lastname1,
				"lastname2": util.SqlNullStringToStringPtr(user.Lastname2),
				"password":  password,
				"email":     "invalid-email",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "TooShortPassword",
			body: gin.H{
				"password":  "123",
				"username":  user.Username,
				"name1":     user.Name1,
				"name2":     util.SqlNullStringToStringPtr(user.Name2),
				"lastname1": user.Lastname1,
				"lastname2": util.SqlNullStringToStringPtr(user.Lastname2),
				"email":     user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/users"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func randomUser(t *testing.T) (user db.User, password string) {
	password = util.RandomString(10)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)
	user = db.User{
		Username:       util.RandomUsername(),
		Name1:          util.RandomUsername(),
		Name2:          sql.NullString{String: "", Valid: false},
		Lastname1:      util.RandomUsername(),
		Lastname2:      sql.NullString{String: "", Valid: false},
		HashedPassword: hashedPassword,
		Email:          util.RandomEmail(),
	}
	n1 := util.RandomInt(0, 100)
	if n1%2 == 0 {
		user.Name2 = sql.NullString{String: util.RandomUsername(), Valid: true}
	}
	n1 = util.RandomInt(0, 99)
	if n1%2 == 0 {
		user.Lastname2 = sql.NullString{String: util.RandomUsername(), Valid: true}
	}
	return user, password
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotUser userResponse
	err = json.Unmarshal(data, &gotUser)

	require.NoError(t, err)
	require.Equal(t, user.Username, gotUser.Username)
	require.Equal(t, user.Name1, gotUser.Name1)
	require.Equal(t, util.SqlNullStringToStringPtr(user.Name2), gotUser.Name2)
	require.Equal(t, user.Lastname1, gotUser.Lastname1)
	require.Equal(t, util.SqlNullStringToStringPtr(user.Lastname2), gotUser.Lastname2)
	require.Equal(t, user.Email, gotUser.Email)
}
