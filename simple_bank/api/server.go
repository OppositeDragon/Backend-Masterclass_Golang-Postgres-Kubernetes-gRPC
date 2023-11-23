package api

import (
	db "simple_bank/db/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()
	// custom validation
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.getAccounts)
	router.PATCH("/accounts/:id", server.updateAccount)
	router.DELETE("/accounts/:id", server.deleteAccount)

	router.POST("/transfers", server.createTransfer)
	// router.POST("/transfers", server.createTransfer)

	server.router = router
	return server

}
func (server *Server) Start(address string) error {
	return server.router.Run(address)

}
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
