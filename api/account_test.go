package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mockdb "github.com/dbracic21-foi/simplebank/db/mock"
	db "github.com/dbracic21-foi/simplebank/db/sqlc"
	"github.com/dbracic21-foi/simplebank/token"
	"github.com/dbracic21-foi/simplebank/util"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetAccountAPI(t *testing.T) {
	user, _ := randomUser(t)
	account := randomAccount(user.Username)
	testcases := []struct {
		name           string
		accountID      int64
		setupAuth      func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildstubs     func(mockStore *mockdb.MockStore)
		checkResponses func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildstubs: func(mockStore *mockdb.MockStore) {
				mockStore.EXPECT().
					GetAccounts(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			checkResponses: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyAccount(t, recorder.Body, account)
			},
		},
		{
			name:      "UnathorizedUser",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "unathorized", time.Minute)
			},
			buildstubs: func(mockStore *mockdb.MockStore) {
				mockStore.EXPECT().
					GetAccounts(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			checkResponses: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:      "NoAuthorization",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildstubs: func(mockStore *mockdb.MockStore) {
				mockStore.EXPECT().
					GetAccounts(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponses: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},

		{
			name:      "Not Found",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildstubs: func(govno *mockdb.MockStore) {
				govno.EXPECT().
					GetAccounts(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Accounts{}, sql.ErrNoRows)
			},
			checkResponses: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "Internal server error",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildstubs: func(mockStore *mockdb.MockStore) {
				mockStore.EXPECT().
					GetAccounts(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Accounts{}, sql.ErrConnDone)
			},
			checkResponses: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "Bad request",
			accountID: 0,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildstubs: func(mockStore *mockdb.MockStore) {
				mockStore.EXPECT().
					GetAccounts(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponses: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}
	for i := range testcases {
		tc := testcases[i]

		t.Run(tc.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)

			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			//BUILD STUBS
			tc.buildstubs(store)
			//start test server

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts/%d", tc.accountID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)
			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponses(t, recorder)

		})

	}
}

// TODO: CREATE AND DELETE ACC DO ON THE SAME WAY
func TestCreateAccountAPI(t *testing.T) {
	user, _ := randomUser(t)
	account := randomAccount(user.Username)

	testCases := []struct {
		name           string
		body           gin.H
		buildstubs     func(mockStore *mockdb.MockStore)
		checkResponses func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"Owner":    account.Owner,
				"Currency": account.Currency,
			},

			buildstubs: func(mockStore *mockdb.MockStore) {
				arg := db.CreateAccountParams{
					Owner:    account.Owner,
					Currency: account.Currency,
					Balance:  0,
				}
				mockStore.EXPECT().
					CreateAccount(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(account, nil)

			},
			checkResponses: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"owner":    account.Owner,
				"currency": account.Currency,
			},

			buildstubs: func(mockStore *mockdb.MockStore) {

				mockStore.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Accounts{}, sql.ErrConnDone)
			},
			checkResponses: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidCurrency",
			body: gin.H{
				"owner":    account.Owner,
				"currency": "invalid currency",
			},

			buildstubs: func(mockStore *mockdb.MockStore) {

				mockStore.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponses: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)

			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			//BUILD STUBS
			tc.buildstubs(store)
			//start test server

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)
			url := "/accounts"

			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponses(t, recorder)
		})
	}
}

func randomAccount(owner string) db.Accounts {
	return db.Accounts{
		ID:       util.RandomInt(1, 1000),
		Owner:    owner,
		Currency: util.RandomCurrency(),
		Balance:  util.RadnomMoney(),
	}
}

func requireBodyAccount(t *testing.T, body *bytes.Buffer, account db.Accounts) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotAccount db.Accounts
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)
	require.Equal(t, account, gotAccount)

}
func requireBodyAccounts(t *testing.T, body *bytes.Buffer, accounts []db.Accounts) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotAccounts db.Accounts
	err = json.Unmarshal(data, &gotAccounts)
	require.NoError(t, err)
	require.Equal(t, accounts, gotAccounts)
}
