package api

import (
	"database/sql"
	"net/http"
	db "simple_bank/db/sqlc"
	util "simple_bank/util"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createUserRequest struct {
	Username  string `json:"username" binding:"required,alphanum,min=3,max=24"`
	Name1     string `json:"name1" binding:"required,min=2"`
	Name2     string `json:"name2"`
	Lastname1 string `json:"lastname1" binding:"required,min=2"`
	Lastname2 string `json:"lastname2"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,alphanum,min=8,max=60"`
}
type createUserResponse struct {
	Username  string  `json:"username" binding:"required,alphanum,min=3,max=24"`
	Name1     string  `json:"name1" binding:"required,min=2"`
	Name2     *string `json:"name2"`
	Lastname1 string  `json:"lastname1" binding:"required,min=2"`
	Lastname2 *string `json:"lastname2"`
	Email     string  `json:"email" binding:"required,email"`
	Password  string  `json:"password" binding:"required,alphanum,min=8,max=60"`
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
	var name2 sql.NullString
	if req.Name2 != "" {
		name2 = sql.NullString{String: req.Name2, Valid: true}

	}
	var lastname2 sql.NullString
	if req.Lastname2 != "" {
		lastname2 = sql.NullString{String: req.Lastname2, Valid: true}
	}
	arg := db.CreateUserParams{
		Username:       req.Username,
		Name1:          req.Name1,
		Name2:          name2,
		Lastname1:      req.Lastname1,
		Lastname2:      lastname2,
		Email:          req.Email,
		HashedPassword: hashedPassword,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, value := err.(*pq.Error); value {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	responseUser := createUserResponse{
		Username:  user.Username,
		Name1:     user.Name1,
		Lastname1: user.Lastname1,
		Email:     user.Email,
		Password:  user.HashedPassword,
	}
	responseUser.Name2 = func() *string {
		if user.Name2.Valid {
			return &user.Name2.String
		} else {
			return nil
		}
	}()
	responseUser.Lastname2 = func() *string {
		if user.Lastname2.Valid {
			return &user.Lastname2.String
		} else {
			return nil
		}
	}()
	ctx.JSON(http.StatusOK, responseUser)
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
