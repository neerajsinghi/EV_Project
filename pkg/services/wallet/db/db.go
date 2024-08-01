package db

import (
	"bikeRental/pkg/entity"
	"bikeRental/pkg/repo/bikeDevice"
	"bikeRental/pkg/repo/generic"
	"bikeRental/pkg/repo/wallet"
	bookingDB "bikeRental/pkg/services/booking/db"
	"bikeRental/pkg/services/city"
	"bikeRental/pkg/services/motog"
	"bikeRental/pkg/services/notifications/notify"
	pdb "bikeRental/pkg/services/plan/pDB"
	predefnotification "bikeRental/pkg/services/predefNotification"
	"bikeRental/pkg/services/users/udb"
	utils "bikeRental/pkg/util"
	"context"
	"fmt"
	"sort"
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

func (s *service) FindForPlan(plan string) ([]entity.WalletS, error) {
	return getWalletTotalForPlan(plan)
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
	if document.PaymentID != "" && document.DepositedMoney > 0 && document.RefundableMoney == 0 {
		capture, err := utils.Capture(document.PaymentID, int(document.DepositedMoney))
		if err != nil {
			return WalletTotal{}, err
		}
		document.CaptureData = capture
	}
	if document.PaymentID != "" && document.RefundableMoney > 0 && document.DepositedMoney == 0 {
		capture, err := utils.Capture(document.PaymentID, int(document.RefundableMoney))
		if err != nil {
			return WalletTotal{}, err
		}
		document.CaptureData = capture
	}
	if document.PaymentID != "" && document.RefundableMoney > 0 && document.DepositedMoney > 0 {
		capture, err := utils.Capture(document.PaymentID, int(document.RefundableMoney+document.DepositedMoney))
		if err != nil {
			return WalletTotal{}, err
		}
		document.CaptureData = capture
	}
	if document.RefundedMoney > 0 && document.PaymentID != "" {
		refund, err := utils.Refund(document.PaymentID, int(document.RefundedMoney))
		if err != nil {
			return WalletTotal{}, err
		}
		document.CaptureData = refund
	}
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
	if document.DepositedMoney > 0 && document.PaymentID != "" {
		s.CheckMyBooking(document.UserID)
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
	refundIDMap := make(map[string]bool)
	for _, w := range wallets {
		totalBalance += w.DepositedMoney
		totalBalance -= w.UsedMoney
		refundableMoney += w.RefundableMoney - w.RefundedMoney
		if refundableMoney > 0 {
			refundIDMap[w.PaymentID] = true
		}
	}
	for _, w := range wallets {
		if w.RefundedMoney > 0 {
			if refundIDMap[w.PaymentID] {
				refundIDMap[w.PaymentID] = false
			}
		}
	}
	if refundableMoney <= 0 {
		refundIDMap = nil
	}
	refundID := ""
	for key, value := range refundIDMap {
		if value {
			refundID = key
			break
		}
	}
	if refundID == "" {
		refundableMoney = 0
	}
	walletL := WalletTotal{
		Wallets:         wallets,
		TotalBalance:    totalBalance,
		RefundableMoney: refundableMoney,
		RefundPaymentID: refundID,
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
func getWalletTotalForPlan(plan string) ([]entity.WalletS, error) {
	pipelineL := bson.A{
		bson.D{{Key: "$match", Value: bson.D{{Key: "plan_id", Value: bson.D{{Key: "$ne", Value: ""}}}}}},
		bson.D{{Key: "$addFields", Value: bson.D{{Key: "userId", Value: bson.D{{Key: "$toObjectId", Value: "$user_id"}}}}}},
		bson.D{{Key: "$lookup", Value: bson.D{{Key: "from", Value: "users"}, {Key: "localField", Value: "userId"}, {Key: "foreignField", Value: "_id"}, {Key: "as", Value: "userData"}}}},
		bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$userData"}}}},
		bson.D{
			{Key: "$lookup",
				Value: bson.D{
					{Key: "from", Value: "booking"},
					{Key: "localField", Value: "plan._id"},
					{Key: "foreignField", Value: "plan._id"},
					{Key: "as", Value: "bookings"},
				},
			},
		},
	}
	if plan != "" {
		pipeline := bson.A{bson.D{{Key: "$match", Value: bson.D{{Key: "plan_id", Value: plan}}}}}
		pipeline = append(pipeline, pipelineL...)
		pipelineL = pipeline
	}
	var repoG = generic.NewRepository("wallet")
	cursor, err := repoG.Aggregate(pipelineL)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	var wallets []entity.WalletS
	for cursor.Next(context.Background()) {
		var wallet entity.WalletS
		if err = cursor.Decode(&wallet); err != nil {
			return nil, err
		}
		wallet.EndTime = new(primitive.DateTime)
		//use plan validity for end time
		validity, _ := strconv.Atoi(wallet.Plan.Validity)
		*wallet.EndTime = primitive.NewDateTimeFromTime(wallet.CreatedTime.Time().Add(time.Hour * time.Duration(validity*24)))
		wallet.BookingCount = new(int)
		*wallet.BookingCount = 0
		for _, booking := range wallet.Bookings {
			if booking.Plan.ID.Hex() == wallet.PlanID && booking.ProfileID == wallet.UserID && booking.CreatedTime.Time().Unix() >= wallet.CreatedTime.Time().Unix() && booking.CreatedTime.Time().Unix() <= wallet.EndTime.Time().Unix() {
				*wallet.BookingCount += 1
			}
		}
		wallet.Bookings = nil
		wallets = append(wallets, wallet)
	}
	return wallets, nil
}

func (w *service) CheckMyBooking(userId string) {
	booking, err := bookingDB.NewService().GetMyLatestBooking(userId)
	if err != nil {
		fmt.Println(err)
		return
	}

	planList, err := pdb.NewService().GetPlans("hourly", booking.City)
	if err != nil {
		fmt.Println(err)
		return
	}

	wallet, err := w.FindMy(booking.ProfileID)
	if err != nil {
		fmt.Println(err)
		return
	}
	timeSpent := 0
	walletAmount := wallet.TotalBalance
	extendedPrice := 0.0
	extendedTime := 0
	city, err := city.InCity(booking.BikeWithDevice.Location.Coordinates[1], booking.BikeWithDevice.Location.Coordinates[0])
	if err != nil || city.Name != booking.City {
		fmt.Println(err)
		if booking.Profile.FirebaseToken != nil {
			predef, err := predefnotification.Get("outOfGeofence")
			if err == nil && predef.Name == "outOfGeofence" {
				notify.NewService().SendNotification(predef.Title, predef.Body, booking.Profile.ID.Hex(), predef.Type, *booking.Profile.FirebaseToken)
			}
		}
		if booking.BikeWithDevice.Type == "moto" {
			motog.ImmoblizeDevice(1, booking.BikeWithDevice.Name)
		} else {
			motog.ImmoblizeDeviceRoadcast(booking.DeviceID, "engineStop")
		}
		filter := bson.M{"device_id": booking.DeviceID}
		repoBike := bikeDevice.NewRepository("bikeDevice")

		set := bson.M{"$set": bson.M{"immobilizeds": true}}
		repoBike.UpdateOne(filter, bson.M{"$set": set})

	}
	sort.Slice(planList, func(i, j int) bool {
		return planList[i].EndingMinutes < planList[j].EndingMinutes
	})
	for i, plan := range planList {
		if plan.EndingMinutes != 0 {
			if walletAmount == plan.Price {
				// Calculate the time this plan can provide

				// Add the time to the total timeSpent
				timeSpent = plan.EndingMinutes
				// Deduct the plan's price from the wallet

				walletAmount -= plan.Price

				break

			} else if i > 0 && (planList[i-1].Price < walletAmount && walletAmount < plan.Price) {
				timeSpent = planList[i-1].EndingMinutes
				walletAmount -= planList[i-1].Price
				break
			}
		} else if plan.EveryXMinutes != 0 {
			extendedPrice = plan.Price
			extendedTime = plan.EveryXMinutes
		}
	}
	if walletAmount > 0 {
		timeEx := int(walletAmount/extendedPrice) * extendedTime
		if timeEx >= 1 {
			timeSpent += timeEx
			walletAmount -= float64(timeSpent/extendedTime) * (extendedPrice)
		}
	}

	bookingDB.AddTimeRemaining(booking.ID.Hex(), timeSpent-int(time.Now().Unix()/60)+int(booking.StartTime/60))

}
