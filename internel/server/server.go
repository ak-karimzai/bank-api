package server

import (
	"github.com/ak-karimzai/bank-api/internel/db"
	"github.com/ak-karimzai/bank-api/internel/handlers"
	"github.com/gin-gonic/gin"
)

type Server struct {
	AccountHandler handlers.AccountHandler
	Router         *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{
		AccountHandler: *handlers.NewAccountHandler(store),
	}
	router := gin.Default()

	router.POST("/accounts", server.AccountHandler.CreateAccount)
	router.GET("/accounts", server.AccountHandler.ListAccounts)
	router.GET("/accounts/:id", server.AccountHandler.GetAccount)

	server.Router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.Router.Run(address)
}
