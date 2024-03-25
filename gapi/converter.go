package gapi

import (
	db "github.com/dbracic21-foi/simplebank/db/sqlc"
	"github.com/dbracic21-foi/simplebank/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertUser(user db.User) *pb.User {
	return &pb.User{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangetAt: timestamppb.New(user.PasswordChangedAt),
		CreatedAt:         timestamppb.New(user.CreatedAt),
	}

}
func convertAccount(account db.Account) *pb.Account {
	return &pb.Account{
		Owner:    account.Owner,
		Balance:  account.Balance,
		Currency: account.Currency,
	}
}
func convertListAccounts(accounts []db.Account) []*pb.Account {
	result := make([]*pb.Account, len(accounts))
	for i, account := range accounts {
		result[i] = convertAccount(account)
	}
	return result
}
func convertTransfers(transfer db.Transfer) *pb.Transfer {
	return &pb.Transfer{
		FromAccountId: transfer.FromAccountID,
		ToAccountId:   transfer.ToAccountID,
		Amount:        transfer.Amount,
	}
}
