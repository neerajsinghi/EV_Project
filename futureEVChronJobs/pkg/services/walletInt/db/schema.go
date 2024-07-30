package wdb

import (
	"futureEVChronJobs/pkg/entity"
)

type Wallet interface {
	FindMy(userId string) (WalletTotal, error)
}

type WalletTotal struct {
	Wallets         []entity.WalletS
	TotalBalance    float64
	RefundableMoney float64
	UserData        *entity.ProfileDB
}
