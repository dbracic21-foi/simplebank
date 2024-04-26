package gapi

import (
	"context"

	db "github.com/dbracic21-foi/simplebank/db/sqlc"
	"github.com/dbracic21-foi/simplebank/pb"
	"github.com/dbracic21-foi/simplebank/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) ListAccounts(ctx context.Context, req *pb.ListAccountsRequest) (*pb.ListAccountsResponse, error) {

	authPayload, err := server.authorizationUser(ctx, []string{util.BankRole, util.DepositorRole})
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	if authPayload.Username != req.GetOwner() {
		return nil, status.Errorf(codes.PermissionDenied, "cannot list other user's accounts")

	}

	arg := db.ListAccountsParams{
		Owner:  req.GetOwner(),
		Limit:  req.GetPage(),
		Offset: req.GetPage() * 1,
	}

	listAccounts, err := server.store.ListAccounts(ctx, arg)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user %s", err)

	}
	rsp := &pb.ListAccountsResponse{
		Accounts: convertListAccounts(listAccounts)}
	return rsp, nil

}
