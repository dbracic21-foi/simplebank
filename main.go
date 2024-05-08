package main

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"

	worker "github.com/dbracic21-foi/simplebank/Worker"
	"github.com/dbracic21-foi/simplebank/api"
	db "github.com/dbracic21-foi/simplebank/db/sqlc"
	_ "github.com/dbracic21-foi/simplebank/doc/statik"
	"github.com/dbracic21-foi/simplebank/gapi"
	"github.com/dbracic21-foi/simplebank/mail"
	"github.com/dbracic21-foi/simplebank/pb"
	"github.com/dbracic21-foi/simplebank/util"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang/mock/mockgen/model"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rakyll/statik/fs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

var interruptSignal = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot load config:")

	}
	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	ctx, stop := signal.NotifyContext(context.Background(), interruptSignal...)
	defer stop()

	connpool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {

		log.Fatal().Msg("Cannot connect to db")
	}

	runDBMigration(config.MigrationURL, config.DBSource)

	store := db.NewStore(connpool)

	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}
	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)

	waitGroup, ctx := errgroup.WithContext(ctx)
	go runTaskProcessor(ctx, waitGroup, config, redisOpt, store)
	go rungGatewayServer(ctx, waitGroup, config, store, taskDistributor)
	rungRPCServer(ctx, waitGroup, config, store, taskDistributor)
	err = waitGroup.Wait()
	if err != nil {
		log.Fatal().Err(err).Msg("error during wait group")

	}

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

func runTaskProcessor(
	ctx context.Context,
	waitGroup *errgroup.Group,
	config util.Config,
	redisOpt asynq.RedisClientOpt,
	store db.Store,

) {
	mailer := mail.NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store, mailer)
	log.Info().Msg("start task processor")
	err := taskProcessor.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to  start task processor")
	}
	waitGroup.Go(func() error {
		<-ctx.Done()
		log.Info().Msg("graceful shutdwon  task processor")
		taskProcessor.Shutdown()
		log.Info().Msg("task processor stoped")
		return nil
	})

}

func rungRPCServer(
	ctx context.Context,
	waitGroup *errgroup.Group,
	config util.Config,
	store db.Store,
	taskDistributor worker.TaskDistributor,
) {

	server, err := gapi.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot start server")
	}
	grpcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)
	grpcServer := grpc.NewServer(grpcLogger)
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listenerGRPC, err := net.Listen("tcp", config.GRPCServerAdress)
	if err != nil {
		log.Fatal().Msg("Cannot start gRPC server")
	}
	waitGroup.Go(func() error {
		log.Info().Msgf("start gRPC server at %s", listenerGRPC.Addr().String())

		err = grpcServer.Serve(listenerGRPC)
		if err != nil {
			if errors.Is(err, grpc.ErrServerStopped) {
				return nil
			}
			log.Error().Err(err).Msg("grpc server failed to serve")
			return err
		}
		return nil

	})
	waitGroup.Go(func() error {
		<-ctx.Done()
		log.Info().Msg("graceful shutdown gRPC server")
		grpcServer.GracefulStop()
		log.Info().Msg("gRPC server stoped")
		return nil
	})

}

func rungGatewayServer(
	ctx context.Context,
	waitGroup *errgroup.Group,
	config util.Config,
	store db.Store,
	taskDistributor worker.TaskDistributor,
) {

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
	//ctx, cancel := context.WithCancel(context.Background())
	//defer cancel()

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
	httpServer := &http.Server{
		Handler: gapi.HttpLogger(mux),
		Addr:    config.HTTPServerAdress,
	}

	//	listenerHTTP, err := net.Listen("tcp", config.HTTPServerAdress)
	//	if err != nil {
	//		log.Fatal().Msg("Cannot create a listener for HTTP server")

	waitGroup.Go(func() error {
		log.Info().Msgf("start HTTP gateway server at %s", httpServer.Addr)
		err = httpServer.ListenAndServe()
		if err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return nil
			}
			log.Error().Err(err).Msg("http server failed to serve")
			return err
		}
		return nil
	})
	waitGroup.Go(func() error {
		<-ctx.Done()
		log.Info().Msg("graceful shutdown HTTP gateway server")

		err := httpServer.Shutdown(context.Background())
		if err != nil {
			log.Error().Err(err).Msg("http server failed to shutdown")
			return err
		}
		log.Info().Msg("HTTP gateway server stoped")
		return nil

	})

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
