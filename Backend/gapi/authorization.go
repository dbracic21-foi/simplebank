package gapi

import (
	"context"
	"fmt"
	"strings"

	"github.com/dbracic21-foi/simplebank/token"
	"google.golang.org/grpc/metadata"
)

const (
	authoriaztionHeader = "authorization"
	authorizationBearer = "bearer"
)

func (server *Server) authorizationUser(ctx context.Context, accesibleRoles []string) (*token.Payload, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("missing metadata")
	}
	values := md.Get(authoriaztionHeader)
	if len(values) == 0 {
		return nil, fmt.Errorf("missing authorization header")
	}
	authHeader := values[0]
	fields := strings.Fields(authHeader)
	if len(fields) < 2 {
		return nil, fmt.Errorf("invalid authorization header format")
	}
	authType := strings.ToLower(fields[0])
	if authType != authorizationBearer {
		return nil, fmt.Errorf("unsupported authorization type")
	}
	accessToken := fields[1]
	payload, err := server.tokenMaker.VerifyToken(accessToken)
	if err != nil {
		return nil, fmt.Errorf("invalid access token : %s", err)
	}
	if !hasPermission(payload.Role, accesibleRoles) {
		return nil, fmt.Errorf("unauthorized")
	}

	return payload, nil
}

func hasPermission(userRole string, accesibleRoles []string) bool {
	for _, role := range accesibleRoles {
		if userRole == role {
			return true
		}

	}
	return false
}
