package gapi

import (
	"context"
	"fmt"
	"testing"
	"time"

	worker "github.com/dbracic21-foi/simplebank/Worker"
	db "github.com/dbracic21-foi/simplebank/db/sqlc"
	"github.com/dbracic21-foi/simplebank/token"
	"github.com/dbracic21-foi/simplebank/util"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
)

func newTestServer(t *testing.T, store db.Store, TaskDistributor worker.TaskDistributor) *Server {
	config := util.Config{
		TokenSymetricKey:    util.RadnomString(32),
		AccessTokenDuration: time.Minute,
	}
	server, err := NewServer(config, store, TaskDistributor)
	require.NoError(t, err)

	return server

}
func newContextWithBearerToken(t *testing.T, tokenMaker token.Maker, username string, role string, duration time.Duration) context.Context {
	ctx := context.Background()
	accessToken, _, err := tokenMaker.CreateToken(username, role, duration)
	require.NoError(t, err)
	bearerToken := fmt.Sprintf("%s %s ", authorizationBearer, accessToken)
	md := metadata.MD{
		authoriaztionHeader: []string{
			bearerToken,
		},
	}
	return metadata.NewIncomingContext(ctx, md)
}
