package api

import (
	"database/sql"
	"fmt"
	"net/http"
	db "simple_bank/db/sqlc"
	util "simple_bank/util"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type createUserRequest struct {
	Username  string  `json:"username" binding:"required,alphanum,min=3,max=24"`
	Name1     string  `json:"name1" binding:"required,min=2"`
	Name2     *string `json:"name2"`
	Lastname1 string  `json:"lastname1" binding:"required,min=2"`
	Lastname2 *string `json:"lastname2"`
	Email     string  `json:"email" binding:"required,email"`
	Password  string  `json:"password" binding:"required,alphanum,min=8,max=60"`
}
type userResponse struct {
	Username          string    `json:"username" binding:"required,alphanum,min=3,max=24"`
	Name1             string    `json:"name1" binding:"required,min=2"`
	Name2             *string   `json:"name2"`
	Lastname1         string    `json:"lastname1" binding:"required,min=2"`
	Lastname2         *string   `json:"lastname2"`
	Email             string    `json:"email" binding:"required,email"`
	PasswordChangedAt time.Time `json:"passwordChangedAt"`
	CreatedAt         time.Time `json:"createdAt"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	arg := db.CreateUserParams{
		Username:       req.Username,
		Name1:          req.Name1,
		Name2:          util.StringPtrToSqlNullString(req.Name2),
		Lastname1:      req.Lastname1,
		Lastname2:      util.StringPtrToSqlNullString(req.Lastname2),
		Email:          req.Email,
		HashedPassword: hashedPassword,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, value := err.(*pq.Error); value {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	responseUser := newUserResponse(user)
	ctx.JSON(http.StatusOK, responseUser)
}

type loginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum,min=3,max=24"`
	Password string `json:"password" binding:"required,alphanum,min=8,max=60"`
}

type loginUserResponse struct {
	SessionId             uuid.UUID    `json:"sessionId"`
	AccessToken           string       `json:"accessToken"`
	AccessTokenExpiresAt  time.Time    `json:"accessTokenExpiresAt"`
	RefreshToken          string       `json:"refreshToken"`
	RefreshTokenExpiresAt time.Time    `json:"refreshTokenExpiresAt"`
	User                  userResponse `json:"user"`
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	err = util.CheckPasswordHash(req.Password, user.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(user.Username, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(user.Username, server.config.RefreshTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:               refreshPayload.Id,
		Username:         user.Username,
		AccessToken:      accessToken,
		AccessExpiresAt:  accessPayload.ExpiredAt,
		RefreshToken:     refreshToken,
		RefreshExpiresAt: refreshPayload.ExpiredAt,
		UserAgent:        util.StringToSqlNullString(ctx.Request.UserAgent()),
		ClientIp:         util.StringToSqlNullString(ctx.ClientIP()),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	response := loginUserResponse{
		SessionId:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
		User:                  newUserResponse(user),
	}
	ctx.JSON(http.StatusOK, response)
}

type renewAccessTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type renewAccessTokenResponse struct {
	AccessToken          string    `json:"accessToken"`
	AccessTokenExpiresAt time.Time `json:"accessTokenExpiresAt"`
}

func (server *Server) renewAccessToken(ctx *gin.Context) {
	var req renewAccessTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	refreshPayload, err := server.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	session, err := server.store.GetSession(ctx, refreshPayload.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if session.IsBlocked {
		err := fmt.Errorf("session is blocked")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	if session.Username != refreshPayload.Username {
		err := fmt.Errorf("session user mismatch")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	if session.RefreshToken != req.RefreshToken {
		err := fmt.Errorf("session token mismatch")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	accessToken, accessPayload, err := server.tokenMaker.CreateToken(session.Username, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	session, err = server.store.UpdateSessionAccess(ctx, db.UpdateSessionAccessParams{
		ID:              refreshPayload.Id,
		AccessToken:     accessToken,
		AccessExpiresAt: accessPayload.ExpiredAt,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	response := renewAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiredAt,
	}
	ctx.JSON(http.StatusOK, response)
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		Username:          user.Username,
		Name1:             user.Name1,
		Name2:             util.SqlNullStringToStringPtr(user.Name2),
		Lastname1:         user.Lastname1,
		Lastname2:         util.SqlNullStringToStringPtr(user.Lastname2),
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
}

// type getUserRequest struct {
// 	ID int64 `uri:"id" binding:"required,min=1"`
// }

// func (server *Server) getUser(ctx *gin.Context) {
// 	var req getUserRequest
// 	if err := ctx.ShouldBindUri(&req); err != nil {
// 		ctx.JSON(http.StatusBadRequest, errorResponse(err))
// 		return
// 	}
// 	user, err := server.store.GetUser(ctx, req.ID)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			ctx.JSON(http.StatusNotFound, errorResponse(err))
// 		}
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return
// 	}
// 	ctx.JSON(http.StatusOK, user)
// }

// type getUsersRequest struct {
// 	Offset *int32 `form:"offset" binding:"required,min=0"`
// 	Limit  int32  `form:"limit" binding:"required,min=5,max=100"`
// }

// func (server *Server) getUsers(ctx *gin.Context) {
// 	var req getUsersRequest
// 	if err := ctx.ShouldBindQuery(&req); err != nil {
// 		ctx.JSON(http.StatusBadRequest, errorResponse(err))
// 		return
// 	}
// 	arg := db.GetUsersParams{
// 		Limit:  req.Limit,
// 		Offset: *req.Offset * req.Limit,
// 	}
// 	users, err := server.store.GetUsers(ctx, arg)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			ctx.JSON(http.StatusNotFound, errorResponse(err))
// 		}
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return
// 	}
// 	ctx.JSON(http.StatusOK, users)
// }

// type updateUserRequest struct {
// 	Amount int64 `json:"amount" binding:"required,gte=0"`
// }

// func (server *Server) updateUser(ctx *gin.Context) {
// 	var req1 getUserRequest
// 	var req2 updateUserRequest
// 	if err := ctx.ShouldBindUri(&req1); err != nil {
// 		ctx.JSON(http.StatusBadRequest, errorResponse(err))
// 		return
// 	}
// 	if err := ctx.ShouldBindJSON(&req2); err != nil {
// 		ctx.JSON(http.StatusBadRequest, errorResponse(err))
// 		return
// 	}
// 	arg := db.AddAmountUserParams{
// 		ID:     req1.ID,
// 		Amount: req2.Amount,
// 	}
// 	user, err := server.store.AddAmountUser(ctx, arg)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return
// 	}
// 	ctx.JSON(http.StatusOK, user)

// }

// func (server *Server) deleteUser(ctx *gin.Context) {
// 	var req getUserRequest
// 	if err := ctx.ShouldBindUri(&req); err != nil {
// 		ctx.JSON(http.StatusBadRequest, errorResponse(err))
// 		return
// 	}
// 	err := server.store.DeleteUser(ctx, req.ID)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			ctx.JSON(http.StatusNotFound, errorResponse(err))
// 		}
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return
// 	}
// 	ctx.JSON(http.StatusOK, "User deleted successfully")
// }
