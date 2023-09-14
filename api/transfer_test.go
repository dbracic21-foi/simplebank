package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "github.com/dbracic21-foi/simplebank/db/mock"
	db "github.com/dbracic21-foi/simplebank/db/sqlc"
	"github.com/dbracic21-foi/simplebank/util"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestTransferAPI(t *testing.T) {
	amount := int64(10)
	account1 := randomAccount()
	account2 := randomAccount()
	account3 := randomAccount()
	account4 := randomAccount()

	account1.Currency = util.USD
	account2.Currency = util.USD
	account3.Currency = util.CAD
	account4.Currency = util.USD

	testCases := []struct {
		name           string
		body           gin.H
		buildstubs     func(store *mockdb.MockStore)
		checkResponses func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        util.USD,
			},
			buildstubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccounts(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, nil)
				store.EXPECT().GetAccounts(gomock.Any(), gomock.Eq(account2.ID)).Times(1).Return(account2, nil)

				arg := db.TransferTxParams{
					FromAccountID: account1.ID,
					ToAccountID:   account2.ID,
					Amount:        amount,
				}
				store.EXPECT().TransfersTx(gomock.Any(), gomock.Eq(arg)).Times(1)
			},
			checkResponses: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "FromAccountNotFound ",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        util.USD,
			},
			buildstubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccounts(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(db.Accounts{}, sql.ErrNoRows)
				store.EXPECT().GetAccounts(gomock.Any(), gomock.Eq(account2.ID)).Times(0)
				store.EXPECT().TransfersTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponses: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "ToAccountNotFound",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        util.USD,
			},
			buildstubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccounts(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, nil)
				store.EXPECT().GetAccounts(gomock.Any(), gomock.Eq(account2.ID)).Times(1).Return(db.Accounts{}, sql.ErrNoRows)
				store.EXPECT().TransfersTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponses: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "FromAccountCurrencyMissMatch",
			body: gin.H{
				"from_account_id": account2.ID,
				"to_account_id":   account4.ID,
				"amount":          amount,
				"currency":        util.CAD,
			},
			buildstubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccounts(gomock.Any(), gomock.Eq(account2.ID)).Times(1).Return(account2, nil)
				store.EXPECT().GetAccounts(gomock.Any(), gomock.Eq(account4.ID)).Times(0)
				store.EXPECT().TransfersTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponses: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "ToAccountCurrencyMismatch",
			body: gin.H{
				"from_account_id": account3.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        util.CAD,
			},
			buildstubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccounts(gomock.Any(), gomock.Eq(account3.ID)).Times(1).Return(account3, nil)
				store.EXPECT().GetAccounts(gomock.Any(), gomock.Eq(account2.ID)).Times(1).Return(account2, nil)
				store.EXPECT().TransfersTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponses: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidCurrency",
			body: gin.H{
				"from_account_id": account2.ID,
				"to_account_id":   account3.ID,
				"amount":          amount,
				"currency":        "XYZ",
			},
			buildstubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccounts(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().TransfersTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponses: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "NegativeAmount",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account4.ID,
				"amount":          -amount,
				"currency":        util.EUR,
			},
			buildstubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccounts(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().TransfersTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponses: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "GetAccountError",
			body: gin.H{
				"from_account_id": account2.ID,
				"to_account_id":   account4.ID,
				"amount":          amount,
				"currency":        util.EUR,
			},
			buildstubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccounts(gomock.Any(), gomock.Any()).Times(1).Return(db.Accounts{}, sql.ErrConnDone)
				store.EXPECT().TransfersTx(gomock.Any(), gomock.Any()).Times(0)

			},
			checkResponses: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "TransferTxError",
			body: gin.H{
				"from_account_id": account2.ID,
				"to_account_id":   account4.ID,
				"amount":          amount,
				"currency":        util.USD,
			},
			buildstubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccounts(gomock.Any(), gomock.Eq(account2.ID)).Times(1).Return(account2, nil)
				store.EXPECT().GetAccounts(gomock.Any(), gomock.Eq(account4.ID)).Times(1).Return(account4, nil)
				store.EXPECT().TransfersTx(gomock.Any(), gomock.Any()).Times(1).Return(db.TransferTxResult{}, sql.ErrTxDone)
			},
			checkResponses: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildstubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/transfers"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponses(recorder)
		})
	}

}
