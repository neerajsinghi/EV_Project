package db

import (
	"bikeRental/pkg/entity"
	"bikeRental/pkg/repo/booking"
	"bikeRental/pkg/repo/generic"
	"bikeRental/pkg/repo/wallet"
	db "bikeRental/pkg/services/account/dbs"
	bookedlogic "bikeRental/pkg/services/bookedBike/logic"
	"bikeRental/pkg/services/coupon/cdb"
	"bikeRental/pkg/services/motog"
	"bikeRental/pkg/services/notifications/notify"
	predefnotification "bikeRental/pkg/services/predefNotification"
	"bikeRental/pkg/services/users/udb"
	wdb "bikeRental/pkg/services/walletInt/db"
	"context"
	"fmt"
	"log"
	"sort"

	bDB "bikeRental/pkg/services/bikeDevice/db"
	bikedb "bikeRental/pkg/services/iotBike/db"
	pdb "bikeRental/pkg/services/plan/pDB"
	sdb "bikeRental/pkg/services/station/sDB"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type service struct{}

var (
	repo = booking.NewRepository("booking")
)

func NewService() Booking {

	return &service{}
}

// AddBooking implements Booking.
func (s *service) AddBooking(document entity.BookingDB) (string, error) {
	document.ID = primitive.NewObjectID()
	if document.StartTime == 0 {
		document.StartTime = time.Now().Unix()
	}
	bDB.DeviceBooked(document.DeviceID)
	user, err := repo.Find(bson.M{"status": "started", "profile_id": document.ProfileID}, bson.M{})
	if err == nil && len(user) > 0 {
		return "", errors.New("user already has a booking")
	}
	userData, err := db.GetUser([]string{document.ProfileID})
	if err != nil {
		return "", errors.New("user not found")
	}
	if len(userData) == 0 {
		return "", errors.New("user not found")
	}
	if userData[0].UserBlocked != nil && *userData[0].UserBlocked {
		return "", errors.New("user is blocked")
	}
	device, err := repo.Find(bson.M{"status": "started", "device_id": document.DeviceID}, bson.M{})
	if err == nil && len(device) > 0 {
		return "", errors.New("device already booked")
	}
	deviceData, err := bikedb.GetBike([]int{document.DeviceID})
	if err != nil {
		return "", errors.New("device not found")
	}
	if len(deviceData) == 0 {
		return "", errors.New("device not found")
	}
	startKM := deviceData[0].TotalDistanceFloat
	if err == nil && startKM != 0 {
		document.StartKM = startKM
	}
	document.CreatedTime = primitive.NewDateTimeFromTime(time.Now())
	if document.StartingStationID == "" {
		return "", errors.New("starting station not found")
	}
	station, err := sdb.NewService().GetStationByID(document.StartingStationID)
	if err != nil {
		return "", errors.New("station not found")
	}
	if station.ID.Hex() == "" {
		return "", errors.New("station not found")
	}

	document.StartingStation = &station

	if document.Plan == nil || document.Plan.ID.Hex() == "" {
		return "", errors.New("plan not found")
	}

	plan, err := pdb.NewService().GetPlan(document.Plan.ID.Hex())
	if err == nil {
		document.Plan = &plan
		document.VehicleType = plan.VehicleType
		document.City = plan.City
		document.BookingType = string(plan.Type)
		if plan.Type == "hourly" {
			planList, _ := pdb.NewService().GetPlans("hourly", document.City)

			wallet, err := wdb.NewService().FindMy(document.ProfileID)
			if err != nil {
				fmt.Println(err)
				return "", err
			}
			timeSpent := 0
			walletAmount := wallet.TotalBalance
			extendedPrice := 0.0
			extendedTime := 0
			for _, plan := range planList {
				if plan.EndingMinutes != 0 {
					if walletAmount <= plan.Price {
						// Calculate the time this plan can provide

						// Add the time to the total timeSpent
						timeSpent = plan.EndingMinutes
						// Deduct the plan's price from the wallet
						walletAmount -= plan.Price

						// If there's no money left, stop iterating
						if walletAmount <= 0 {
							break
						}
					}
				} else {
					extendedPrice = plan.Price
					extendedTime = plan.EveryXMinutes
				}
			}
			if walletAmount > 0 {
				timeSpent += int(walletAmount/extendedPrice) * extendedTime
			}

			document.RideTimeRemaining = timeSpent - int(time.Now().Unix()/60) + int(document.StartTime/60)
		}
	}
	udb.ChangeServiceType(document.ProfileID, document.BookingType)
	if userData[0].FirebaseToken != nil {
		predef, err := predefnotification.Get("bookingStarted")
		if err == nil && predef.Name == "bookingStarted" {
			notify.NewService().SendNotification(predef.Title, predef.Body, document.ProfileID, predef.Type, *userData[0].FirebaseToken)
		}
	}
	success, err := repo.InsertOne(document)
	if err != nil {
		return "", errors.New("unable to  booking")
	}
	return success, nil
}

// DeleteBooking implements Booking.
func (s *service) DeleteBooking(id string) error {
	idObject, _ := primitive.ObjectIDFromHex(id)
	err := repo.DeleteOne(bson.M{"_id": idObject})
	if err != nil {
		return errors.New("unable to delete booking")
	}
	return nil
}

// GetAllBookings implements Booking.
func (s *service) GetAllBookings(status, bType, vType string) ([]entity.BookingOut, error) {
	filter := bson.D{}
	if status != "" && status != "all" {
		filter = append(filter, bson.E{Key: "status", Value: status})

	}
	if bType != "" && bType != "all" {
		filter = append(filter, bson.E{Key: "booking_type", Value: bType})
	}
	if vType != "" && vType != "all" {
		filter = append(filter, bson.E{Key: "vehicle_type", Value: vType})
	}

	pipeline := createPipeline(filter)

	data, err := repo.Aggregate(pipeline)
	if err != nil {
		return nil, errors.New("unable to get booking")
	}
	return data, nil
}

func createPipeline(filter bson.D) primitive.A {
	pipeline := bson.A{}
	if len(filter) > 0 {
		pipeline = append(pipeline, bson.D{{Key: "$match", Value: filter}})
	}
	pipeline = append(pipeline, bson.A{
		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "userObjectId", Value: bson.D{{Key: "$toObjectId", Value: "$profile_id"}}},
		}}},

		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "iotBike"},
			{Key: "localField", Value: "device_id"},
			{Key: "foreignField", Value: "deviceId"},
			{Key: "as", Value: "bikeWithDevice"},
		}}},

		bson.D{{Key: "$unwind", Value: "$bikeWithDevice"}},

		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "users"},
			{Key: "localField", Value: "userObjectId"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "profile"},
		}}},

		bson.D{{Key: "$unwind", Value: "$profile"}},

		bson.D{{Key: "$project", Value: bson.D{
			{Key: "bookingDb", Value: "$$ROOT"},
			{Key: "bikeWithDevice", Value: 1},
			{Key: "profile", Value: 1},
		}}},

		bson.D{{Key: "$replaceRoot", Value: bson.D{{Key: "newRoot", Value: bson.D{
			{Key: "bookingDb", Value: "$bookingDb"},
			{Key: "bikeWithDevice", Value: "$bikeWithDevice"},
			{Key: "profile", Value: "$profile"},
		}}}}},
	}...)
	return pipeline
}
func GetAllHourlyBookings() ([]entity.BookingOut, error) {
	filter := bson.D{
		{Key: "status", Value: "started"},
		{Key: "booking_type", Value: "hourly"},
	}
	pipeline := createPipeline(filter)

	data, err := repo.Aggregate(pipeline)
	if err != nil {
		return nil, errors.New("unable to get booking")
	}
	return data, nil
}
func GetAllStartedBookings() ([]entity.BookingOut, error) {
	filter := bson.D{
		{Key: "status", Value: "started"},
	}
	pipeline := createPipeline(filter)

	data, err := repo.Aggregate(pipeline)
	if err != nil {
		return nil, errors.New("unable to get booking")
	}
	return data, nil
}
func GetMyBookingCount(userID string) (int64, error) {
	filter := bson.M{"profile_id": userID}
	return repo.Count(filter)
}

// GetAllMyBooking implements Booking.
func (s *service) GetAllMyBooking(userID, bType string) ([]entity.BookingOut, error) {
	filter := bson.D{{Key: "profile_id", Value: userID}}
	if bType != "" && bType != "all" {
		filter = append(filter, bson.E{Key: "booking_type", Value: bType})
	}
	pipeline := createPipeline(filter)

	data, err := repo.Aggregate(pipeline)
	if err != nil {
		return nil, errors.New("unable to get booking")
	}
	return data, nil
}

// GetMyBooking implements Booking.
func (s *service) GetMyLatestBooking(userID string) (*entity.BookingOut, error) {
	filter := bson.D{{Key: "profile_id", Value: userID}, {Key: "status", Value: bson.D{{Key: "$ne", Value: "completed"}}}}
	pipeline := createPipeline(filter)

	booking, err := repo.Aggregate(pipeline)
	if err != nil {
		return nil, err
	}
	if len(booking) == 0 {
		return nil, errors.New("booking not found")
	}
	if len(booking) > 1 {
		for i := 1; i < len(booking); i++ {
			if booking[i].Status == "started" {
				return &booking[i], nil
			}
		}
	}
	return &booking[0], nil
}
func (*service) GetBookingByID(id string) (*entity.BookingOut, error) {
	return GetBooking(id)
}

// UpdateBooking implements Booking.
func (s *service) UpdateBooking(id string, document entity.BookingDB) (string, error) {

	set := bson.M{}
	var devices []entity.IotBikeDB
	booking, err := GetBooking(id)

	if err != nil {
		return "", err
	}
	deviceList := make([]int, 0)
	deviceList = append(deviceList, booking.DeviceID)
	devices, err = bikedb.GetBike(deviceList)
	if err != nil {
		return "", errors.New("device not found")
	}

	if document.Price != nil {
		set["price"] = *document.Price
	}
	if document.Return != nil {
		set["return"] = document.Return
	}
	if document.VehicleType != "" {
		set["vehicle_type"] = document.VehicleType
	}
	if document.Status == "completed" {
		totalDistanceInt := devices[0].TotalDistanceFloat

		set["end_km"] = totalDistanceInt
		set["total_distance"] = totalDistanceInt - booking.StartKM
		userTotalDist := (totalDistanceInt - booking.StartKM)
		greenPoints := int64(userTotalDist * 5)
		carbonSaved := userTotalDist * 80
		profile := entity.ProfileDB{
			GreenPoints:    greenPoints,
			CarbonSaved:    carbonSaved,
			TotalTravelled: userTotalDist,
		}
		set["green_points"] = greenPoints
		set["carbon_saved"] = carbonSaved
		db.UpdateUser(booking.ProfileID, profile)
		endTime := time.Now().Unix()
		set["end_time"] = endTime
		if document.EndingStationID != "" {
			set["ending_station_id"] = document.EndingStationID
			station, err := sdb.NewService().GetStationByID(document.EndingStationID)
			if err == nil {
				set["ending_station"] = &station
			}
			bDB.DeviceReturned(booking.DeviceID, document.EndingStationID)
			bookedlogic.ChangeOnGoing(id)
			timeBooked := int((endTime - booking.StartTime) / 60)
			set["time_booked"] = timeBooked
			if booking.BookingType == "hourly" {
				plan, err := pdb.NewService().GetPlans(booking.BookingType, booking.City)
				if err == nil && len(plan) > 0 {
					price := float64(0)
					maxPrice := float64(0)
					maxTime := 0
					extendedPrice := 0.0
					extendedTime := 0
					sort.Slice(plan, func(i, j int) bool {
						return plan[i].EndingMinutes < plan[j].EndingMinutes
					})
					for _, p := range plan {
						if p.EndingMinutes != 0 {
							if timeBooked <= p.EndingMinutes {
								price = p.Price
								break
							}
							if maxPrice < p.Price {
								maxPrice = p.Price
							}
							if maxTime < p.EndingMinutes {
								maxTime = p.EndingMinutes
							}
						} else if p.EveryXMinutes != 0 {
							extendedPrice = p.Price
							extendedTime = p.EveryXMinutes
						}
					}
					if price == 0 {
						price = maxPrice
						timeBooked -= maxTime
						for timeBooked > 0 {
							price += extendedPrice
							timeBooked -= extendedTime
						}
					}

					if booking.CouponCode != "" {
						coupon, err := cdb.GetCouponByCode(booking.CouponCode)

						if err == nil {

							if coupon.CouponType == "discount" {
								discount := (price * coupon.Discount / 100)
								if discount > coupon.MaxValue {
									discount = coupon.MaxValue
								}
								if discount < coupon.MinValue {
									discount = coupon.MinValue
								}
								price = price - discount
								set["discounted_amount"] = discount
							} else if coupon.CouponType == "freeRide" {
								if price > coupon.MaxValue {
									price = price - coupon.MaxValue
									set["discounted_amount"] = coupon.MaxValue
								} else {
									set["discounted_amount"] = price
									price = 0
								}
							}
						}
					}

					wall := entity.WalletS{
						ID:          primitive.NewObjectID(),
						UserID:      booking.ProfileID,
						UsedMoney:   price,
						BookingID:   id,
						Description: "Booking",
					}
					wallet.NewRepository("wallet").InsertOne(wall)

				}
				udb.ChangeServiceType(document.ProfileID, "")
			}
			referralPath(booking)
			userData, err := db.GetUser([]string{booking.ProfileID})
			if err == nil && len(userData) == 0 {
				return "", errors.New("user not found")
			}
			if userData[0].FirebaseToken != nil {
				predef, err := predefnotification.Get("bookingCompleted")
				if err == nil && predef.Name == "bookingCompleted" {
					notify.NewService().SendNotification(predef.Title, predef.Body, booking.ProfileID, predef.Type, *userData[0].FirebaseToken)
				}
			}
			set["status"] = document.Status
		} else {
			return "", errors.New("ending station not found")
		}
		idObject, _ := primitive.ObjectIDFromHex(id)
		set["updated_at"] = primitive.NewDateTimeFromTime(time.Now())

		return repo.UpdateOne(bson.M{"_id": idObject}, bson.M{"$set": set})
	}
	if document.Status == "resumed" && booking.Status == "stopped" {

		data, err := s.AddBooking(
			entity.BookingDB{
				ProfileID:         booking.ProfileID,
				DeviceID:          booking.DeviceID,
				StartingStationID: booking.StartingStationID,
				Plan:              booking.Plan,
				Status:            "started",
				BookingType:       booking.BookingType,
				VehicleType:       booking.VehicleType,
			},
		)
		if err != nil {
			return "", errors.New("unable to resume booking")
		}
		if booking.BikeWithDevice.Type == "moto" {
			motog.ImmoblizeDevice(0, booking.BikeWithDevice.Name)
		} else {
			motog.ImmoblizeDeviceRoadcast(booking.DeviceID, "engineResume")
		}
		return data, nil

	}
	return "", errors.New("status not completed")
}

func referralPath(booking *entity.BookingOut) {
	repoG := generic.NewRepository("referral")
	referral, err := repoG.FindOne(bson.M{"referral_of_id": booking.ProfileID, "referral_status": "pending"}, bson.M{})
	if err == nil {
		defer referral.Close(context.Background())
		var ref entity.ReferralDB
		for referral.Next(context.Background()) {
			err := referral.Decode(&ref)

			if err == nil {
				bonus := 40.0
				value, err := cdb.GetCouponByType("referral")
				if err == nil {
					bonus = value.MaxValue
				}
				wall := entity.WalletS{
					ID:             primitive.NewObjectID(),
					UserID:         ref.ReferrerBy,
					DepositedMoney: float64(bonus),
					Description:    "Referral Bonus",
				}
				wallet.NewRepository("wallet").InsertOne(wall)
				predef, err := predefnotification.Get("referalComplete")
				if err == nil && predef.Name == "referalComplete" {
					notify.NewService().SendNotification(predef.Title, predef.Body, ref.ReferrerBy, predef.Type, *ref.ReferredByProfile.FirebaseToken)
				}
				repoG.UpdateOne(bson.M{"referrer_by": ref.ReferrerBy, "referral_of_id": ref.ReferralOfId}, bson.M{"$set": bson.M{"referral_status": "completed"}})
				break
			}
		}
	}
}

func GetBooking(id string) (*entity.BookingOut, error) {
	idObject, _ := primitive.ObjectIDFromHex(id)
	filter := bson.D{{Key: "_id", Value: idObject}}
	pipeline := createPipeline(filter)

	booking, err := repo.Aggregate(pipeline)
	if err != nil {
		return nil, err
	}
	if len(booking) == 0 {
		return nil, errors.New("booking not found")
	}
	return &booking[0], nil
}

func AddTimeRemaining(id string, timeRemaining int) {
	idObject, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": idObject}
	set := bson.M{}
	set["ride_time_remaining"] = timeRemaining
	set["update_time"] = time.Now()
	repo.UpdateOne(filter, bson.M{"$set": set})
}

func ChangeStatusStopped(id string, price float64, endTime int64, endKm float64) {
	idObject, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": idObject}
	booking, _ := GetBooking(id)
	set := bson.M{}
	set["status"] = "stopped"
	set["price"] = price
	set["end_time"] = endTime
	set["end_km"] = endKm
	set["total_distance"] = (endKm - booking.StartKM)
	userTotalDist := (endKm - booking.StartKM)
	set["updated_time"] = primitive.NewDateTimeFromTime(time.Now())
	greenPoints := int64(userTotalDist * 5)
	carbonSaved := userTotalDist * 80
	profile := entity.ProfileDB{
		GreenPoints:    greenPoints,
		CarbonSaved:    carbonSaved,
		TotalTravelled: userTotalDist,
	}
	set["green_points"] = greenPoints
	set["carbon_saved"] = carbonSaved
	db.UpdateUser(booking.ProfileID, profile)
	_, err := repo.UpdateOne(filter, bson.M{"$set": set})
	log.Println("error", err)
}

func (s *service) GetBookingByPlanAndID(planID, userId string) ([]entity.BookingOut, error) {
	id, _ := primitive.ObjectIDFromHex(planID)
	filter := bson.D{{Key: "plan._id", Value: id}, {Key: "profile_id", Value: userId}}
	pipeline := createPipeline(filter)

	booking, err := repo.Aggregate(pipeline)
	if err != nil {
		return nil, errors.New("booking not found")
	}
	if len(booking) == 0 {
		return nil, errors.New("booking not found")
	}
	return booking, nil
}
