package gapi

import (
	"context"

	db "github.com/dbracic21-foi/simplebank/db/sqlc"
	"github.com/dbracic21-foi/simplebank/pb"
	"github.com/dbracic21-foi/simplebank/util"
	"github.com/dbracic21-foi/simplebank/val"
	"github.com/lib/pq"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreatUserRequest) (*pb.CreatUserResponse, error) {
	violations := validateCreateUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}
	hashedPassword, err := util.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to hash password : %s", err)

	}

	arg := db.CreateUserParams{
		Username:       req.GetUsername(),
		HashedPassword: hashedPassword,
		FullName:       req.GetFullName(),
		Email:          req.GetEmail(),
	}

	user, err := server.store.CreateUser(ctx, arg)

	if err != nil {
		if pqerr, ok := err.(*pq.Error); ok {
			switch pqerr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "username already exists %s", err)

			}
		}
		return nil, status.Errorf(codes.Internal, "failed to create user %s", err)

	}
	rsp := &pb.CreatUserResponse{
		User: convertUser(user),
	}
	return rsp, nil

}
func validateCreateUserRequest(req *pb.CreatUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))

	}

	if err := val.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))

	}
	if err := val.ValidateFullName(req.GetFullName()); err != nil {
		violations = append(violations, fieldViolation("full_name", err))

	}
	if err := val.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, fieldViolation("email", err))

	}
	return violations
}
