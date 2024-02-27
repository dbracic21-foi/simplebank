package gapi

import (
	"testing"
	"time"

	worker "github.com/dbracic21-foi/simplebank/Worker"
	db "github.com/dbracic21-foi/simplebank/db/sqlc"
	"github.com/dbracic21-foi/simplebank/util"
	"github.com/stretchr/testify/require"
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
