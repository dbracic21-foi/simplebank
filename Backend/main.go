package main

import (
	"context"
	"database/sql"
	"net"
	"net/http"
	"os"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	worker "github.com/dbracic21-foi/simplebank/Worker"
	"github.com/dbracic21-foi/simplebank/api"
	db "github.com/dbracic21-foi/simplebank/db/sqlc"
	_ "github.com/dbracic21-foi/simplebank/doc/statik"
	"github.com/dbracic21-foi/simplebank/gapi"
	"github.com/dbracic21-foi/simplebank/pb"
	"github.com/dbracic21-foi/simplebank/util"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang/mock/mockgen/model"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq" // Import the PostgreSQL driver package
	"github.com/rakyll/statik/fs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Msg("Cannot load config:")

	}
	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {

		log.Fatal().Msg("Cannot connect to db")
	}

	runDBMigration(config.MigrationURL, config.DBSource)

	store := db.NewStore(conn)

	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}
	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)
	go runTaskProcessor(redisOpt, store)
	go rungGatewayServer(config, store, taskDistributor)
	rungRPCServer(config, store, taskDistributor)

}
func runDBMigration(migrationURL string, dbSource string) {

	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal().Msg("cannot create new migrate instance")
	}
	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().Msg("failed to run migrate up:")
	}
	log.Info().Msg("db migrated succesfuly")

}

func runTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store) {

	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store)
	log.Info().Msg("start task processor")
	err := taskProcessor.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start task processor")
	}

}

func rungRPCServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) {

	server, err := gapi.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal().Msg("Cannot start server")
	}
	grpcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)
	grpcServer := grpc.NewServer(grpcLogger)
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listenerGRPC, err := net.Listen("tcp", config.GRPCServerAdress)
	if err != nil {
		log.Fatal().Msg("Cannot start gRPC server")
	}
	log.Info().Msgf("start gRPC server at %s", listenerGRPC.Addr().String())
	err = grpcServer.Serve(listenerGRPC)
	if err != nil {
		log.Fatal().Msg("cannot start gRPC server")
	}

}

func rungGatewayServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) {

	server, err := gapi.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal().Msg("Cannot start server")
	}
	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})
	grpMux := runtime.NewServeMux(jsonOption)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pb.RegisterSimpleBankHandlerServer(ctx, grpMux, server)
	if err != nil {
		log.Fatal().Msg("cannot register handler server : ")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpMux)

	statikFS, err := fs.New()
	if err != nil {
		log.Fatal().Msg("cannot create statik fs : ")
	}
	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))

	mux.Handle("/swagger/", swaggerHandler)

	listenerHTTP, err := net.Listen("tcp", config.HTTPServerAdress)
	if err != nil {
		log.Fatal().Msg("Cannot create a listener for HTTP server")
	}
	log.Info().Msgf("start HTTP gateway server at %s", listenerHTTP.Addr().String())
	handler := gapi.HttpLogger(mux)
	err = http.Serve(listenerHTTP, handler)
	if err != nil {
		log.Fatal().Msg("cannot start HTTP gateway server")
	}

}
func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal().Msg("Cannot start server")
	}
	err = server.Start(config.HTTPServerAdress)
	if err != nil {
		log.Fatal().Msg("Cannot start server")
	}

}
