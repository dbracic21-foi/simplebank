package gapi

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"testing"

	worker "github.com/dbracic21-foi/simplebank/Worker"
	mockwk "github.com/dbracic21-foi/simplebank/Worker/mock"
	mockdb "github.com/dbracic21-foi/simplebank/db/mock"
	db "github.com/dbracic21-foi/simplebank/db/sqlc"
	"github.com/dbracic21-foi/simplebank/pb"
	"github.com/dbracic21-foi/simplebank/util"
	"github.com/golang/mock/gomock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type eqCreateUserTxParamsMatcher struct {
	arg      db.CreateUserTxParams
	password string
	user     db.Users
}

func (expected eqCreateUserTxParamsMatcher) Matches(x interface{}) bool {
	actualArg, ok := x.(db.CreateUserTxParams)

	if !ok {
		return false
	}
	err := util.CheckPassword(expected.password, actualArg.HashedPassword)
	if err != nil {
		return false
	}
	expected.arg.HashedPassword = actualArg.HashedPassword
	if !reflect.DeepEqual(expected.arg.CreateUserParams, actualArg.CreateUserParams) {
		return false

	}

	err = actualArg.AfterCreate(expected.user)
	return err == nil

}

func (e eqCreateUserTxParamsMatcher) String() string {
	return fmt.Sprintf("Matches arg  %v and password  %v", e.arg, e.password)
}
func EqCreateUserTxParams(arg db.CreateUserTxParams, password string, user db.Users) gomock.Matcher {
	return eqCreateUserTxParamsMatcher{arg, password, user}
}
func randomUser(t *testing.T) (user db.Users, password string) {
	password = util.RadnomString(6)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	user = db.Users{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		Email:          util.RandomEmail(),
		FullName:       util.RandomOwner(),
	}
	return
}

func TestCreateUserAPI(t *testing.T) {
	user, password := randomUser(t)

	testcases := []struct {
		name           string
		req            *pb.CreatUserRequest
		buildstubs     func(mockStore *mockdb.MockStore, taskDistributor *mockwk.MockTaskDistributor)
		checkResponses func(t *testing.T, res *pb.CreatUserResponse, err error)
	}{
		{
			name: "OK",
			req: &pb.CreatUserRequest{

				Username: user.Username,
				Password: password,
				FullName: user.FullName,
				Email:    user.Email,
			},
			buildstubs: func(store *mockdb.MockStore, taskDistributor *mockwk.MockTaskDistributor) {
				arg := db.CreateUserTxParams{
					CreateUserParams: db.CreateUserParams{
						Username: user.Username,
						FullName: user.FullName,
						Email:    user.Email,
					},
				}
				store.EXPECT().
					CreateUserTx(gomock.Any(), EqCreateUserTxParams(arg, password, user)).
					Times(1).
					Return(db.CreateUserTxResult{Users: user}, nil)

				taskPayload := &worker.PayloadSendVerifyEmail{
					Username: user.Username,
				}

				taskDistributor.EXPECT().
					DistributeTaskSendVerifyEmail(gomock.Any(), taskPayload, gomock.Any()).
					Times(1).
					Return(nil)
			},
			checkResponses: func(t *testing.T, res *pb.CreatUserResponse, err error) {
				require.NoError(t, err)
				require.NotEmpty(t, res)
				createdUserResponse := res.GetUser()
				require.Equal(t, user.Username, createdUserResponse.Username)
				require.Equal(t, user.FullName, createdUserResponse.FullName)
				require.Equal(t, user.Email, createdUserResponse.Email)

			},
		},
		{
			name: "InternalError",
			req: &pb.CreatUserRequest{
				Username: user.Username,
				Password: password,
				FullName: user.FullName,
				Email:    user.Email,
			},
			buildstubs: func(store *mockdb.MockStore, taskDistributor *mockwk.MockTaskDistributor) {
				store.EXPECT().
					CreateUserTx(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.CreateUserTxResult{}, sql.ErrConnDone)

				taskDistributor.EXPECT().
					DistributeTaskSendVerifyEmail(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
			},

			checkResponses: func(t *testing.T, res *pb.CreatUserResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Internal, st.Code())

			},
		},
		{
			name: "DuplicateUserName",
			req: &pb.CreatUserRequest{
				Username: user.Username,
				Password: password,
				FullName: user.FullName,
				Email:    user.Email,
			},
			buildstubs: func(store *mockdb.MockStore, taskDistributor *mockwk.MockTaskDistributor) {
				store.EXPECT().
					CreateUserTx(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.CreateUserTxResult{}, &pq.Error{Code: "23505"})
				taskDistributor.EXPECT().
					DistributeTaskSendVerifyEmail(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponses: func(t *testing.T, res *pb.CreatUserResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.AlreadyExists, st.Code())
			},
		},

		{
			name: "InvalidEmail",
			req: &pb.CreatUserRequest{
				Username: user.Username,
				Password: password,
				FullName: user.FullName,
				Email:    "invalid-email",
			},
			buildstubs: func(store *mockdb.MockStore, taskDistributor *mockwk.MockTaskDistributor) {
				store.EXPECT().
					CreateUserTx(gomock.Any(), gomock.Any()).
					Times(0)
				taskDistributor.EXPECT().
					DistributeTaskSendVerifyEmail(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponses: func(t *testing.T, res *pb.CreatUserResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.InvalidArgument, st.Code())

			},
		},
	}
	for i := range testcases {
		tc := testcases[i]
		t.Run(tc.name, func(t *testing.T) {
			storectrl := gomock.NewController(t)
			defer storectrl.Finish()
			store := mockdb.NewMockStore(storectrl)

			taskctrl := gomock.NewController(t)
			defer taskctrl.Finish()

			taskDistributor := mockwk.NewMockTaskDistributor(taskctrl)
			tc.buildstubs(store, taskDistributor)

			server := newTestServer(t, store, taskDistributor)
			res, err := server.CreateUser(context.Background(), tc.req)
			tc.checkResponses(t, res, err)
		})
	}
}
