package grpcserver

import (
	"github.com/ak-karimzai/bank-api/internel/db"
	"github.com/ak-karimzai/bank-api/internel/pb"
	"github.com/ak-karimzai/bank-api/internel/repository"
	"github.com/ak-karimzai/bank-api/internel/token"
	"github.com/ak-karimzai/bank-api/internel/util"
	"github.com/ak-karimzai/bank-api/internel/worker"
)

type GRPCServer struct {
	pb.UnimplementedSimpleBankServer
	AccountRepo     repository.AccountRepository
	TransferRepo    repository.TransferRepository
	UserRepo        repository.UserRepository
	SessionRepo     repository.SessionRepository
	Config          util.Config
	TokenMaker      token.Maker
	TaskDistributor worker.TaskDistributor
}

func NewGRPCServer(
	config util.Config, store db.Store, taskDistributor worker.TaskDistributor) (*GRPCServer, error) {
	tokenMaker, err := token.NewJWTMaker(config.TokenSymmtricKey)
	if err != nil {
		return nil, err
	}

	server := &GRPCServer{
		AccountRepo:     store,
		TransferRepo:    store,
		UserRepo:        store,
		SessionRepo:     store,
		Config:          config,
		TokenMaker:      tokenMaker,
		TaskDistributor: taskDistributor,
	}

	return server, nil
}
