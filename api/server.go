package api

import (
	"fmt"

	db "github.com/Klaygogo/simplebank/db/sqlc"
	token "github.com/Klaygogo/simplebank/token"
	"github.com/Klaygogo/simplebank/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	config     util.Config
	store      db.Store
	router     *gin.Engine
	tokenMaker token.Maker
}

func (server *Server) setupRoutes() {
	authRouter := server.router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRouter.POST("/accounts", server.createAccount)
	authRouter.GET("/accounts/:id", server.getAccount)
	authRouter.GET("/accounts", server.listAccount)

	authRouter.POST("/transfers", server.createTransfer)
	authRouter.GET("/transfers/:id", server.getTransfer)
	authRouter.GET("/transfers", server.listTransfer)

	server.router.POST("/users", server.createUser)
	server.router.POST("/users/login", server.loginUser)
	authRouter.GET("/users/:username", server.getUser)

	authRouter.POST("/tokens/renew_access", server.renewAccessToken)
}

func NewServer(store db.Store, config util.Config) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSyncKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create tokenMaker: %w", err)
	}
	r := gin.Default()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}
	server := &Server{
		config:     config,
		store:      store,
		router:     r,
		tokenMaker: tokenMaker,
	}
	server.setupRoutes()
	return server, nil
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
