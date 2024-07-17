package db

import (
	"bikeRental/pkg/entity"
	"bikeRental/pkg/repo/wallet"
	bookingDB "bikeRental/pkg/services/booking/db"
	"bikeRental/pkg/services/notifications/notify"
	pdb "bikeRental/pkg/services/plan/pDB"
	"bikeRental/pkg/services/users/udb"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type service struct{}

var (
	repo = wallet.NewRepository("wallet")
)

// NewService creates a new service
func NewService() Wallet {
	return &service{}
}

// DeleteOne implements Wallet.
func (s *service) DeleteOne(filter primitive.M) error {
	panic("unimplemented")
}

// Find implements Wallet.
func (s *service) Find() ([]WalletTotal, error) {
	return getWalletTotalAll()
}

// FindMy implements Wallet.
func (s *service) FindMy(userId string) (WalletTotal, error) {
	return getWalletTotal(userId)
}

// InsertOne implements Wallet.
func (s *service) InsertOne(document entity.WalletS) (WalletTotal, error) {
	document.ID = primitive.NewObjectID()
	user, err := udb.NewService().GetUserById(document.UserID)
	if err != nil && user.Name != "" {
		return WalletTotal{}, err
	}
	if document.BookingID != "" {
		booking, err := bookingDB.GetBooking(document.BookingID)
		if err != nil {
			return WalletTotal{}, err
		}
		document.Booking = booking
	}
	if document.PlanID != "" {
		plan, err := pdb.NewService().GetPlan(document.PlanID)
		if err != nil {
			return WalletTotal{}, err
		}
		document.Plan = &plan
		profile := entity.ProfileDB{
			PlanID: &document.PlanID,
		}
		udb.NewService().UpdateUser(document.UserID, profile)
	}
	document.CreatedTime = primitive.NewDateTimeFromTime(time.Now())

	_, err = repo.InsertOne(document)

	if err != nil {
		return WalletTotal{}, err
	}
	if document.RefundedMoney != 0 {
		if user.FirebaseToken != nil {
			refund := strconv.FormatFloat(document.RefundedMoney, 'f', -1, 64)
			notify.NewService().SendNotification("Refund", "Refund of "+refund+" has been credited to your wallet", document.UserID, "refund", *user.FirebaseToken)
		}
	}
	return getWalletTotal(document.UserID)
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
