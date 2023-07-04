package server

import (
	"testing"
	"time"

	"github.com/ak-karimzai/bank-api/internel/db"
	"github.com/ak-karimzai/bank-api/internel/handlers"
	"github.com/ak-karimzai/bank-api/internel/middlewares"
	"github.com/ak-karimzai/bank-api/internel/token"
	"github.com/ak-karimzai/bank-api/internel/util"
	"github.com/ak-karimzai/bank-api/internel/validators"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

type Server struct {
	AccountHandler  *handlers.AccountHandler
	TransferHandler *handlers.TransferHandler
	UserHandler     *handlers.UserHandler
	Router          *gin.Engine
	Config          util.Config
	TokenMaker      token.Maker
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(config.TokenSymmtricKey)
	if err != nil {
		return nil, err
	}

	server := &Server{
		AccountHandler:  handlers.NewAccountHandler(store),
		TransferHandler: handlers.NewTransferHandler(store),
		UserHandler:     handlers.NewUserHandler(store, tokenMaker, &config),
		Config:          config,
		TokenMaker:      tokenMaker,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validators.ValidCurrency)
	}

	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/users", server.UserHandler.CreateUser)
	router.POST("/users/login", server.UserHandler.LoginUser)

	authRouts := router.Group("/").Use(
		middlewares.AuthMiddleware(server.TokenMaker))

	authRouts.POST("/accounts", server.AccountHandler.CreateAccount)
	authRouts.GET("/accounts", server.AccountHandler.ListAccounts)
	authRouts.GET("/accounts/:id", server.AccountHandler.GetAccount)
	authRouts.POST("/transfers", server.TransferHandler.CreateTransfer)

	server.Router = router
}

func (server *Server) Start(address string) error {
	return server.Router.Run(address)
}

func NewTestServer(t *testing.T, store db.Store) *Server {
	config := util.Config{
		TokenSymmtricKey:    util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	gin.SetMode(gin.TestMode)
	server, err := NewServer(config, store)
	require.NoError(t, err)
	return server
}
