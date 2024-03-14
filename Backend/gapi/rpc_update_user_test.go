package gapi

import (
	"context"
	"database/sql"
	"testing"
	"time"

	mockdb "github.com/dbracic21-foi/simplebank/db/mock"
	db "github.com/dbracic21-foi/simplebank/db/sqlc"
	"github.com/dbracic21-foi/simplebank/pb"
	"github.com/dbracic21-foi/simplebank/token"
	"github.com/dbracic21-foi/simplebank/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestUpdateUserAPI(t *testing.T) {
	user, _ := randomUser(t)
	newName := util.RandomOwner()
	newEmail := util.RandomEmail()
	invalidEmail := "invalidEmail"

	testcases := []struct {
		name           string
		req            *pb.UpdateUserRequest
		buildstubs     func(mockStore *mockdb.MockStore)
		buildContext   func(t *testing.T, tokenMaker token.Maker) context.Context
		checkResponses func(t *testing.T, res *pb.UpdateUserResponse, err error)
	}{	
		{
			name: "OK",
			req: &pb.UpdateUserRequest{
				Username: user.Username,
				FullName: &newName,
				Email:    &newEmail,
			},
			buildstubs: func(store *mockdb.MockStore) {
				arg := db.UpdateUserParams{
					Username: user.Username,
					FullName: sql.NullString{
						String: newName,
						Valid:  true,
					},
					Email: sql.NullString{
						String: newEmail,
						Valid:  true,
					},
				}
				updateUser := db.Users{
					Username:          user.Username,
					HashedPassword:    user.HashedPassword,
					FullName:          newName,
					Email:             newEmail,
					PasswordChangedAt: user.PasswordChangedAt,
					CreatedAt:         user.CreatedAt,
					IsEmailVerified:   user.IsEmailVerified,
				}
				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(updateUser, nil)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, user.Username, time.Minute)

			},

			checkResponses: func(t *testing.T, res *pb.UpdateUserResponse, err error) {
				require.NoError(t, err)
				require.NotEmpty(t, res)
				updateUser := res.GetUser()
				require.Equal(t, user.Username, updateUser.Username)
				require.Equal(t, newName, updateUser.FullName)
				require.Equal(t, newEmail, updateUser.Email)

			},
		},
		{
			name: "UserNotFound",
			req: &pb.UpdateUserRequest{
				Username: user.Username,
				//TODO : Password: password,
				FullName: &newName,
				Email:    &newEmail,
			},
			buildstubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Users{}, sql.ErrNoRows)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, user.Username, time.Minute)

			},

			checkResponses: func(t *testing.T, res *pb.UpdateUserResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.NotFound, st.Code())

			},
		},
		{
			name: "ExpiredToken",
			req: &pb.UpdateUserRequest{
				Username: user.Username,
				//TODO : Password: password,
				FullName: &newName,
				Email:    &newEmail,
			},
			buildstubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, user.Username, -time.Minute)

			},

			checkResponses: func(t *testing.T, res *pb.UpdateUserResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Unauthenticated, st.Code())

			},
		},
		{
			name: "InvalidEmail",
			req: &pb.UpdateUserRequest{
				Username: user.Username,
				//TODO : Password: password,
				FullName: &newName,
				Email:    &invalidEmail,
			},
			buildstubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, user.Username, time.Minute)

			},

			checkResponses: func(t *testing.T, res *pb.UpdateUserResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.InvalidArgument, st.Code())

			},
		},
		{
			name: "No authorization",
			req: &pb.UpdateUserRequest{
				Username: user.Username,
				//TODO : Password: password,
				FullName: &newName,
				Email:    &newEmail,
			},
			buildstubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return context.Background()

			},

			checkResponses: func(t *testing.T, res *pb.UpdateUserResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Unauthenticated, st.Code())

			},
		},
	}
	for i := range testcases {
		tc := testcases[i]
		t.Run(tc.name, func(t *testing.T) {
			storectrl := gomock.NewController(t)
			defer storectrl.Finish()
			store := mockdb.NewMockStore(storectrl)

			tc.buildstubs(store)
			server := newTestServer(t, store, nil)
			ctx := tc.buildContext(t, server.tokenMaker)
			res, err := server.UpdateUser(ctx, tc.req)
			tc.checkResponses(t, res, err)
		})
	}
}
