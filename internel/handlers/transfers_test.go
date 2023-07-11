package handlers_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ak-karimzai/bank-api/internel/db"
	mockdb "github.com/ak-karimzai/bank-api/internel/db/mock"
	"github.com/ak-karimzai/bank-api/internel/middlewares"
	"github.com/ak-karimzai/bank-api/internel/server"
	"github.com/ak-karimzai/bank-api/internel/token"
	"github.com/ak-karimzai/bank-api/internel/util"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestCreateTransfer(t *testing.T) {
	amount := int64(10)

	user1, _ := randomUser(t)
	user2, _ := randomUser(t)
	user3, _ := randomUser(t)

	account1 := randomAccount(user1.Username)
	account2 := randomAccount(user2.Username)
	account3 := randomAccount(user3.Username)

	account1.Currency = util.USD
	account2.Currency = util.USD
	account3.Currency = util.EUR

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
			},
			buildStubs: func(store *mockdb.MockStore) {

				arg := db.TransferTxParams{
					FromAccountID: account1.ID,
					ToAccountID:   account2.ID,
					Amount:        amount,
					Username:      account1.Owner,
				}

				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Eq(arg)).
					Times(1)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, middlewares.AuthorizationTypeBearer, account1.Owner, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "TryToReachAnotherAccount",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
			},
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).
					Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "NoAuthHeader",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
			},
			buildStubs: func(store *mockdb.MockStore) {

				err := errors.New("account doesn't belong to the sender")
				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).
					Times(1).Return(db.TransferTxResult{}, err)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, middlewares.AuthorizationTypeBearer, account2.Owner, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "NegativeAmount",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          -amount,
			},
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).
					Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, middlewares.AuthorizationTypeBearer, account2.Owner, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "DbError",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
			},
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).
					Times(1).Return(db.TransferTxResult{}, sql.ErrConnDone)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, middlewares.AuthorizationTypeBearer, account2.Owner, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tt.buildStubs(store)

			server := server.NewTestServer(t, store)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tt.body)
			require.NoError(t, err)

			url := "/transfers"
			req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			tt.setupAuth(t, req, server.TokenMaker)
			server.Router.ServeHTTP(recorder, req)
			tt.checkResponse(t, recorder)
		})
	}
}
