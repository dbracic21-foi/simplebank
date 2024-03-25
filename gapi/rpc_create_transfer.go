package gapi

import (
	"context"

	db "github.com/dbracic21-foi/simplebank/db/sqlc"
	"github.com/dbracic21-foi/simplebank/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateTransfer(ctx context.Context, req *pb.CreateTransferRequest) (*pb.CreateTransferResponse, error) {
	arg := db.CreateTransfersParams{
		FromAccountID: req.GetFrom(),
		ToAccountID:   req.GetTo(),
		Amount:        req.GetAmount(),
	}

	transfer, err := server.store.CreateTransfers(ctx, arg)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create transfer %s", err)
	}

	rsp := &pb.CreateTransferResponse{
		Transfer: convertTransfers(transfer),
	}
	return rsp, nil
}
