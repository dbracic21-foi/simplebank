package gapi

import (
	"context"

	db "github.com/dbracic21-foi/simplebank/db/sqlc"
	"github.com/dbracic21-foi/simplebank/pb"
	"github.com/dbracic21-foi/simplebank/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateAccount(ctx context.Context, req *pb.CreateAccountRequest) (*pb.CreateAccountResponse, error) {
	authPayload, err := server.authorizationUser(ctx, []string{util.BankRole, util.DepositorRole})
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	if authPayload.Username != req.GetOwner() {
		return nil, status.Errorf(codes.PermissionDenied, "cannot create account for other user")

	}

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
