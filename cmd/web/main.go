package main

import (
	"context"
	"net"
	"net/http"
	"os"

	"github.com/ak-karimzai/bank-api/internel/db"
	grpcserver "github.com/ak-karimzai/bank-api/internel/grpc_server"
	"github.com/ak-karimzai/bank-api/internel/pb"
	"github.com/ak-karimzai/bank-api/internel/server"
	"github.com/ak-karimzai/bank-api/internel/util"
	"github.com/ak-karimzai/bank-api/internel/worker"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hibiken/asynq"
	"github.com/pressly/goose/v3"
	"github.com/rakyll/statik/fs"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	_ "github.com/ak-karimzai/bank-api/internel/docs/statik"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Error().Msgf("cannot load the configurations: %v", err)
	}

	if config.Environment == util.DevelopmentEnvironment {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	runDbMigrations(config)

	dbConn, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Error().Msgf("cannot connect to db: ", err)
	}

	err = dbConn.Ping(context.Background())
	if err != nil {
		log.Error().Msgf("cannot connect to db: ", err)
	}

	store := db.NewStore(dbConn)

	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisServerAddress,
	}

	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)

	go runTaskProcessor(redisOpt, store)
	go runGatewayServer(config, store, taskDistributor)
	runRpcServer(config, store, taskDistributor)
}

func runDbMigrations(config util.Config) {
	db, err := goose.OpenDBWithDriver(config.DBDriver, config.DBSource)
	if err != nil {
		log.Error().Msgf("goose: failed to open DB: %v\n", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Error().Msgf("goose: failed to close DB: %v\n", err)
		}
	}()

	if err := goose.Run("up", db, config.MigrationDir); err != nil {
		log.Error().Msgf("goose %v: %v", "up", err)
	}
}

func runTaskProcessor(
	redisOpt asynq.RedisClientOpt, store db.Store) {
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store)
	log.Info().Msg("start task processor")
	err := taskProcessor.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start task processor")
	}
}

func runHttpServer(
	config util.Config, store db.Store) {
	httpServer, err := server.NewServer(config, store)
	if err != nil {
		log.Error().Msgf("error while creating server: %v", err)
	}

	log.Info().Msgf("HTTP server is litening in port: %v", config.GRPCServerAddress)
	err = httpServer.Start(config.HTTPServerAddress)
	if err != nil {
		log.Error().Msgf("cannot start server: %v", err)
	}
}

func runRpcServer(
	config util.Config, store db.Store, taskDistributor worker.TaskDistributor) {
	server, err := grpcserver.NewGRPCServer(config, store, taskDistributor)
	if err != nil {
		log.Error().Msgf("error while creating server: %v", err)
	}

	grpcLogger := grpc.UnaryInterceptor(grpcserver.GRPCLogger)
	grpcServer := grpc.NewServer(grpcLogger)
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Error().Msgf("cannot create listener")
	}

	log.Info().Msgf("gRPC server is litening in port: %v", config.GRPCServerAddress)
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Error().Msgf("cannot start gRPC server")
	}
}

func runGatewayServer(
	config util.Config, store db.Store, taskDistributor worker.TaskDistributor) {
	server, err := grpcserver.NewGRPCServer(config, store, taskDistributor)
	if err != nil {
		log.Error().Msgf("error while creating server: %v", err)
	}

	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOption)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Error().Msgf("cannot register handler server: %v", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	statikFS, err := fs.New()
	if err != nil {
		log.Error().Msgf("cannot create statik fs: %v", err)
	}

	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))
	mux.Handle("/swagger/", swaggerHandler)

	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Error().Msgf("cannot create listener")
	}

	log.Info().Msgf("HTTP gateway server is litening in port: %v", config.HTTPServerAddress)
	handler := grpcserver.HttpLogger(mux)
	err = http.Serve(listener, handler)
	if err != nil {
		log.Error().Msgf("cannot start HTTP gateway server")
	}
}
