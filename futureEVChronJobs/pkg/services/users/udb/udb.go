package udb

import (
	"errors"
	"futureEVChronJobs/pkg/entity"
	"futureEVChronJobs/pkg/repo/booking"
	"futureEVChronJobs/pkg/repo/profile"
	"futureEVChronJobs/pkg/services/notifications/notify"
	pdb "futureEVChronJobs/pkg/services/plan/pDB"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type service struct{}

func NewService() UserI {
	return &service{}
}

var (
	repo = profile.NewProfileRepository("users")
)

// GetUserById implements UserI.
func (s *service) GetUserById(id string) (entity.ProfileOut, error) {
	idObject, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": idObject}
	filter["status"] = bson.M{"$ne": "deleted"}
	res, err := repo.Aggregate(createPipeline(filter))
	if err != nil {
		return entity.ProfileOut{}, err
	}
	if len(res) == 0 {
		return entity.ProfileOut{}, nil
	}

	res[0].Password = nil
	if res[0].PlanID != nil {
		res[0].PlanRemainingTime = res[0].PlanEndTime - time.Now().Unix()
		if res[0].PlanRemainingTime < 0 {
			RemovePlan(res[0].ID.Hex())
		}
	}
	if len(res[0].Wallet) == 0 {
		res[0].TotalBalance = 0
	} else {
		for _, wallet := range res[0].Wallet {
			res[0].TotalBalance += wallet.DepositedMoney - wallet.UsedMoney

		}
	}
	if res[0].Booking != nil {
		res[0].TotalRides = int64(len(res[0].Booking))
	} else {
		res[0].TotalRides = 0
	}

	return res[0], err
}

// GetUsers implements UserI.
func (s *service) GetUsers(userType string) ([]entity.ProfileOut, error) {
	filter := bson.M{}
	if userType != "" && userType != "all" {
		filter["roles"] = userType
	}
	if userType == "admin" || userType == "staff" {
		res, err := repo.Aggregate(createStaffPipeline(filter))
		for i := 0; i < len(res); i++ {
			res[i].Password = nil
		}
		return res, err
	}
	res, err := repo.Aggregate(createPipeline(filter))
	for i := 0; i < len(res); i++ {
		res[i].Password = nil

		if len(res[i].Wallet) == 0 {
			res[i].TotalBalance = 0
		} else {
			for _, wallet := range res[i].Wallet {
				res[i].TotalBalance += wallet.DepositedMoney - wallet.UsedMoney

			}
		}
		if res[i].Booking != nil {
			res[i].TotalRides = int64(len(res[i].Booking))
		} else {
			res[i].TotalRides = 0
		}
		if res[i].PlanID != nil {
			res[i].PlanRemainingTime = res[i].PlanEndTime - time.Now().Unix()
			if res[i].PlanRemainingTime < 0 {
				RemovePlan(res[i].ID.Hex())
			}
		}
	}
	return res, err
}

func createPipeline(filter bson.M) bson.A {
	pipeline := bson.A{
		bson.D{{Key: "$match", Value: filter}},
		bson.D{
			{
				Key: "$sort", Value: bson.D{{Key: "created_time", Value: -1}},
			},
		},
		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "userStringId", Value: bson.D{{Key: "$toString", Value: "$_id"}}},
		}}},
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "wallet"},
			{Key: "localField", Value: "userStringId"},
			{Key: "foreignField", Value: "user_id"},
			{Key: "as", Value: "wallet"},
		}}},
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "booking"},
			{Key: "localField", Value: "userStringId"},
			{Key: "foreignField", Value: "profile_id"},
			{Key: "as", Value: "booking"},
		}}},
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "profileDb", Value: "$$ROOT"},
			{Key: "wallet", Value: 1},
			{Key: "booking", Value: 1},
		}}},

		bson.D{{Key: "$replaceRoot", Value: bson.D{{Key: "newRoot", Value: bson.D{
			{Key: "profileDb", Value: "$profileDb"},
			{Key: "wallet", Value: "$wallet"},
			{Key: "booking", Value: "$booking"},
		}}}}},
	}

	return pipeline
}

func createStaffPipeline(filter bson.M) bson.A {
	pipeline := bson.A{
		bson.D{{Key: "$match", Value: filter}},
		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "userStringId", Value: bson.D{{Key: "$toString", Value: "$_id"}}},
		}}},
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "station"},
			{Key: "localField", Value: "userStringId"},
			{Key: "foreignField", Value: "supervisor_id"},
			{Key: "as", Value: "station"},
		}}},
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "profileDb", Value: "$$ROOT"},
			{Key: "station", Value: 1},
		}}},
		bson.D{{Key: "$replaceRoot", Value: bson.D{{Key: "newRoot", Value: bson.D{
			{Key: "profileDb", Value: "$profileDb"},
			{Key: "station", Value: "$station"},
		}}}}},
	}

	return pipeline
}

func RemovePlan(userId string) (string, error) {
	idObject, _ := primitive.ObjectIDFromHex(userId)
	filter := bson.M{"_id": idObject}

	booking, err := booking.NewRepository("booking").FindOne(bson.M{"profile_id": userId, "status": "started"}, bson.M{})
	if err == nil || booking.City != "" {
		return "", errors.New("user has ongoing booking")
	}
	set := bson.M{}
	set["plan_id"] = ""
	set["plan"] = nil
	set["plan_active"] = false
	set["plan_start_time"] = 0
	set["plan_end_time"] = 0
	return repo.UpdateOne(filter, bson.M{"$set": set})
}

func (s *service) UpdateUser(id string, user entity.ProfileDB) (string, error) {
	idObject, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": idObject}
	userN, err := repo.FindOne(filter, bson.M{})
	if err != nil {
		return "", err
	}
	set := bson.M{}
	if user.Name != "" {
		set["name"] = user.Name
	}
	if user.CountryCode != nil && *user.CountryCode != "" {
		set["country_code"] = *user.CountryCode
	}

	if user.UserBlocked != nil {
		set["user_blocked"] = user.UserBlocked
		set["blocked_by"] = user.BlockedBy
		set["blocked_time"] = time.Now()
		set["block_reason"] = user.BlockReason
	}
	if user.IDFrontImage != "" {
		set["id_front_image"] = user.IDFrontImage
	}
	if user.IDBackImage != "" {
		set["id_back_image"] = user.IDBackImage
	}
	if user.DLFrontImage != "" {
		set["dl_front_image"] = user.DLFrontImage
	}
	if user.DLBackImage != "" {
		set["dl_back_image"] = user.DLBackImage
	}
	if user.PlanActive != nil {
		set["plan_active"] = user.PlanActive
	}
	if user.Access != nil {
		set["access"] = user.Access
	}
	if user.DOB != "" {
		set["dob"] = user.DOB
	}
	if user.IDVerified != nil {
		set["id_verified"] = *user.IDVerified
		if *user.IDVerified && userN.FirebaseToken != nil {
			notify.NewService().SendNotification("ID Verified", "Your ID has been verified", userN.ID.Hex(), "idVerified", *userN.FirebaseToken)
		}
	}
	if user.DLVerified != nil {
		set["dl_verified"] = *user.DLVerified
		if *user.DLVerified && userN.FirebaseToken != nil {
			notify.NewService().SendNotification("DL Verified", "Your DL has been verified", userN.ID.Hex(), "idVerified", *userN.FirebaseToken)
		}
	}
	if user.Gender != nil && *user.Gender != "" {
		set["gender"] = *user.Gender
	}
	if user.UrlToProfileImage != nil && *user.UrlToProfileImage != "" {
		set["url_to_profile_image"] = user.UrlToProfileImage
	}
	if user.Role != nil {
		set["role"] = *user.Role
	}
	if user.StatusBool != nil {
		set["status_bool"] = *user.StatusBool
	}
	if user.FirebaseToken != nil {
		set["firebase_token"] = *user.FirebaseToken
	}
	if user.ServiceType != "" {
		set["service_type"] = user.ServiceType
	}
	inc := bson.M{}
	inc["green_points"] = user.GreenPoints
	inc["carbon_saved"] = user.CarbonSaved

	if user.PlanID != nil {
		if *user.PlanID == "" {
			set["plan_id"] = ""
			set["plan"] = nil
			set["plan_active"] = false
			set["plan_start_time"] = 0
			set["plan_end_time"] = 0
		}
		set["plan_id"] = *user.PlanID
		plan, _ := pdb.NewService().GetPlan(*user.PlanID)
		set["plan"] = plan
		planValidity, _ := strconv.Atoi(plan.Validity)
		planEnd := int64(planValidity)*24*60*60 + time.Now().Unix()
		set["plan_active"] = true
		set["plan_start_time"] = time.Now().Unix()
		set["plan_end_time"] = planEnd
		set["service_type"] = plan.Type
	}

	set["update_time"] = time.Now()
	setS := bson.M{"$set": set}
	setS["$inc"] = inc

	return repo.UpdateOne(filter, setS)
}
