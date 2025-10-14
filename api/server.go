package api

import (
	db "github.com/Klaygogo/simplebank/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func (server *Server) setupRoutes() {
	server.router.POST("/accounts", server.createAccount)
	server.router.GET("/accounts/:id", server.getAccount)
	server.router.GET("/accounts", server.listAccount)

	server.router.POST("/transfers", server.createTransfer)
	server.router.GET("/transfers/:id", server.getTransfer)
	server.router.GET("/transfers", server.listTransfer)
}

func NewServer(store db.Store) *Server {
	r := gin.Default()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}
	server := &Server{
		store:  store,
		router: r,
	}
	server.setupRoutes()
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
