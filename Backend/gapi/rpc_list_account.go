package gapi

import (
	"context"

	db "github.com/dbracic21-foi/simplebank/db/sqlc"
	"github.com/dbracic21-foi/simplebank/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) ListAccounts(ctx context.Context, req *pb.ListAccountsRequest) (*pb.ListAccountsResponse, error) {

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
