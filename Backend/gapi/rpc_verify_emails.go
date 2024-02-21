package gapi

import (
	"context"

	db "github.com/dbracic21-foi/simplebank/db/sqlc"
	"github.com/dbracic21-foi/simplebank/pb"
	"github.com/dbracic21-foi/simplebank/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) VerifyEmails(ctx context.Context, req *pb.VerifyEmailsRequest) (*pb.VerifyEmailsResponse, error) {
	violations := validateVerifyEmailsRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}
	txResult, err := server.store.VerifyEmailsTx(ctx, db.VerifyEmailsTxParams{
		EmailId:    req.GetEmailId(),
		SecretCode: req.GetSecretCode(),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal error:  %s", err)
	}

	rsp := &pb.VerifyEmailsResponse{
		IsVerifired: txResult.Users.IsEmailVerified,
	}
	return rsp, nil

}
func validateVerifyEmailsRequest(req *pb.VerifyEmailsRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateEmailId(req.GetEmailId()); err != nil {
		violations = append(violations, fieldViolation("email_id", err))

	}

	if err := val.ValidateSecretCode(req.GetSecretCode()); err != nil {
		violations = append(violations, fieldViolation("password", err))

	}

	return violations
}
