package api

import (
	"errors"
	"fmt"
	"net/http"
	"simple_bank/token"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	ErrNoAuthorizationHeader     = errors.New("no authorization header provided")
	ErrInvalidAuthrizationHeader = errors.New("no authorization header provided")
	ErrUnsupportedAuthorization  = fmt.Errorf("unsupported authorization type")
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authPayload"
)

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			err := ErrNoAuthorizationHeader
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := ErrInvalidAuthrizationHeader
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := ErrUnsupportedAuthorization
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
		}
		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}
