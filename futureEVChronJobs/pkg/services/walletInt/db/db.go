package wdb

import (
	"futureEVChronJobs/pkg/entity"
	"futureEVChronJobs/pkg/repo/wallet"

	"go.mongodb.org/mongo-driver/bson"
)

type service struct{}

var (
	repo = wallet.NewRepository("wallet")
)

// NewService creates a new service
func NewService() Wallet {
	return &service{}
}

// FindMy implements Wallet.
func (s *service) FindMy(userId string) (WalletTotal, error) {
	return getWalletTotal(userId)
}

func getWalletTotal(userID string) (WalletTotal, error) {
	wallets, err := repo.Find(bson.M{"user_id": userID}, nil)
	if err != nil {
		return WalletTotal{}, err
	}
	var totalBalance float64
	var refundableMoney float64
	for _, w := range wallets {
		totalBalance += w.DepositedMoney
		totalBalance -= w.UsedMoney
		refundableMoney += w.RefundableMoney - w.RefundedMoney
	}
	walletL := WalletTotal{
		Wallets:         wallets,
		TotalBalance:    totalBalance,
		RefundableMoney: refundableMoney,
	}
	return walletL, nil
}
func getWalletTotalAll() ([]WalletTotal, error) {
	wallets, err := repo.Find(nil, nil)
	if err != nil {
		return nil, err
	}
	var walletListmap = make(map[string]WalletTotal)

	for _, w := range wallets {
		if _, ok := walletListmap[w.UserID]; !ok {
			walletListmap[w.UserID] = WalletTotal{
				Wallets: make([]entity.WalletS, 0),
			}
			wallet := walletListmap[w.UserID]
			wallet.Wallets = append(wallet.Wallets, w)
			wallet.TotalBalance += w.DepositedMoney - w.UsedMoney
			wallet.RefundableMoney += w.RefundableMoney - w.RefundedMoney
			walletListmap[w.UserID] = wallet
		} else {
			wallet := walletListmap[w.UserID]
			wallet.Wallets = append(wallet.Wallets, w)
			wallet.TotalBalance += w.DepositedMoney - w.UsedMoney
			wallet.RefundableMoney += w.RefundableMoney - w.RefundedMoney
			walletListmap[w.UserID] = wallet
		}
	}

	var walletList []WalletTotal
	for _, v := range walletListmap {
		walletList = append(walletList, v)
	}
	return walletList, nil
}
