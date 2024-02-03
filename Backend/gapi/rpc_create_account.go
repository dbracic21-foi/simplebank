package gapi

import (
	"context"

	db "github.com/dbracic21-foi/simplebank/db/sqlc"
	"github.com/dbracic21-foi/simplebank/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateAccount(ctx context.Context, req *pb.CreateAccountRequest) (*pb.CreateAccountResponse, error) {
	arg := db.CreateAccountParams{
		Owner:    req.GetOwner(),
		Balance:  req.GetBalance(),
		Currency: req.GetCurrency(),
	}

	account, err := server.store.CreateAccount(ctx, arg)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user %s", err)

	}

	rsp := &pb.CreateAccountResponse{
		Account: convertAccount(account),
	}
	return rsp, nil

}
