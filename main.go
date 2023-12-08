package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"

	"github.com/dbracic21-foi/simplebank/api"
	db "github.com/dbracic21-foi/simplebank/db/sqlc"
	"github.com/dbracic21-foi/simplebank/gapi"
	"github.com/dbracic21-foi/simplebank/pb"
	"github.com/dbracic21-foi/simplebank/util"
	_ "github.com/golang/mock/mockgen/model"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq" // Import the PostgreSQL driver package
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load config:", err)

	}
	conn, err := sql.Open(config.DBDRIVER, config.DBSOURCE)
	if err != nil {

		log.Fatal("Cannot connect to db", err)
	}

	store := db.NewStore(conn)
	go rungGatewayServer(config, store)
	rungRPCServer(config, store)

}

func rungRPCServer(config util.Config, store db.Store) {

	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("Cannot start server", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAdress)
	if err != nil {
		log.Fatal("Cannot start server", err)
	}
	log.Printf("start grpc server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start grpc server", err)
	}

}

func rungGatewayServer(config util.Config, store db.Store) {

	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("Cannot start server", err)
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
		log.Fatal("cannot register handler server : ", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpMux)

	listener, err := net.Listen("tcp", config.HTTPServerAdress)
	if err != nil {
		log.Fatal("Cannot create a listener", err)
	}
	log.Printf("start HTTP gateway server at %s", listener.Addr().String())
	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatal("cannot start HTTP gateway server : ", err)
	}

}
func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("Cannot start server", err)
	}
	err = server.Start(config.HTTPServerAdress)
	if err != nil {
		log.Fatal("Cannot start server", err)
	}

}
