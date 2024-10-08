package pdb

import (
	"bikeRental/pkg/entity"
	"bikeRental/pkg/repo/plan"
	"errors"
	"sort"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type service struct{}

var (
	repo = plan.NewRepository("plan")
)

func NewService() PlanI {
	return &service{}
}

// AddPlan implements PlanI.
func (s *service) AddPlan(plan entity.PlanDB) (string, error) {
	plan.ID = primitive.NewObjectID()
	plan.CreatedTime = primitive.NewDateTimeFromTime(time.Now())
	if plan.City == "" {
		return "", errors.New("city is required")
	}
	if plan.VehicleType == "" {
		return "", errors.New("vehicle type is required")
	}
	if plan.Price == 0 && plan.Deposit == nil {
		return "", errors.New("price is required")
	}
	if plan.Type == "" {
		return "", errors.New("type is required")
	}
	if plan.Type == entity.Hourly {
		dat, err := s.GetPlans("hourly", plan.City)
		if err == nil && len(dat) > 0 {
			for _, v := range dat {
				if dat[0].IsActive != nil && *dat[0].IsActive && v.VehicleType == plan.VehicleType {
					if plan.EveryXMinutes == 0 && v.EveryXMinutes == 0 && ((plan.Deposit == nil || *plan.Deposit == 0) && (v.Deposit == nil || *v.Deposit == 0)) {
						if plan.StartingMinutes != 0 && (plan.StartingMinutes > v.StartingMinutes && plan.StartingMinutes < v.EndingMinutes) || (plan.EndingMinutes > v.StartingMinutes && plan.EndingMinutes < v.EndingMinutes) || (plan.StartingMinutes <= v.StartingMinutes && plan.EndingMinutes >= v.EndingMinutes) {
							return "", errors.New("hourly plan already exist for this city")
						}
						if plan.StartingMinutes == 0 && plan.EndingMinutes == 0 {
							return "", errors.New("hourly plan already exist for this city")
						}
						if plan.StartingMinutes == 0 && v.StartingMinutes == 0 {
							return "", errors.New("hourly plan already exist for this city")
						}
					}
					if plan.Deposit != nil && v.Deposit != nil {
						return "", errors.New("hourly plan already exist for this city")
					}
					if plan.EveryXMinutes != 0 && v.EveryXMinutes != 0 {
						return "", errors.New("hourly plan already exist for this city")
					}
				}
			}
		}
	}
	if plan.Type == entity.Rental {
		dat, err := repo.FindOne(bson.M{"city": plan.City, "vehicle_type": plan.VehicleType, "validity": plan.Validity}, bson.M{})
		if err == nil && dat.City != "" {
			return "", errors.New("rental plan already exist for this city")
		}
	}

	plan.IsActive = new(bool)
	*plan.IsActive = true
	if plan.Deposit == nil {
		plan.Deposit = new(float64)
		*plan.Deposit = 0
	}
	data, err := repo.InsertOne(plan)
	if err != nil {
		return "", errors.New("error in inserting plan")
	}
	return data, nil
}

// DeletePlan implements PlanI.
func (s *service) DeletePlan(id string) error {
	idObject, _ := primitive.ObjectIDFromHex(id)
	err := repo.DeleteOne(bson.M{"_id": idObject})
	if err != nil {
		return errors.New("error in deleting plan")
	}
	return nil
}

// GetPlan implements PlanI.
func (s *service) GetPlan(id string) (entity.PlanDB, error) {
	idObject, _ := primitive.ObjectIDFromHex(id)
	data, err := repo.FindOne(bson.M{"_id": idObject}, bson.M{})
	if err != nil {
		return entity.PlanDB{}, errors.New("error in finding plan")
	}
	return data, nil
}

// GetPlans implements PlanI.
func (s *service) GetPlans(pType, city string) ([]entity.PlanDB, error) {
	filter := bson.M{"is_active": true}
	if pType != "" {
		filter["type"] = pType
	}
	if city != "" {
		filter["city"] = city
	}
	data, err := repo.Find(filter, bson.M{})
	if err != nil {
		return nil, errors.New("error in finding plans")
	}
	if pType == string(entity.Rental) {
		sort.Slice(data, func(i, j int) bool {
			validityI, _ := strconv.Atoi(data[i].Validity)
			validityJ, _ := strconv.Atoi(data[j].Validity)
			return validityI < validityJ
		})
	}
	if pType == string(entity.Hourly) {
		sort.Slice(data, func(i, j int) bool {
			return data[i].StartingMinutes < data[j].StartingMinutes
		})
	}
	return data, nil
}

// GetPlans implements PlanI.
func (s *service) GetPlansAdmin(pType, city string) ([]entity.PlanDB, error) {
	filter := bson.M{}
	if pType != "" {
		filter["type"] = pType
	}
	if city != "" {
		filter["city"] = city
	}
	data, err := repo.Find(filter, bson.M{})
	if err != nil {
		return nil, errors.New("error in finding plans")
	}
	if pType == string(entity.Rental) {
		sort.Slice(data, func(i, j int) bool {
			validityI, _ := strconv.Atoi(data[i].Validity)
			validityJ, _ := strconv.Atoi(data[j].Validity)
			return validityI < validityJ
		})
	}
	if pType == string(entity.Hourly) {
		sort.Slice(data, func(i, j int) bool {
			return data[i].StartingMinutes < data[j].StartingMinutes
		})
	}
	return data, nil
}

// UpdatePlan implements PlanI.
func (s *service) UpdatePlan(id string, plan entity.PlanDB) (string, error) {
	idObject, _ := primitive.ObjectIDFromHex(id)
	set := bson.M{}

	if plan.Name != "" {
		set["name"] = plan.Name
	}
	if plan.Price != 0 {
		set["price"] = plan.Price
	}
	if plan.Validity != "" {
		set["validity"] = plan.Validity
	}
	if plan.Description != "" {
		set["description"] = plan.Description
	}
	if plan.Discount != 0 {
		set["discount"] = plan.Discount
	}
	if plan.IsActive != nil {
		set["is_active"] = plan.IsActive
	}
	if plan.StartingMinutes != 0 {
		set["starting_minutes"] = plan.StartingMinutes
	}
	if plan.EndingMinutes != 0 {
		set["ending_minutes"] = plan.EndingMinutes
	}
	if plan.EveryXMinutes != 0 {
		set["every_x_minutes"] = plan.EveryXMinutes
	}
	if plan.City != "" {
		set["city"] = plan.City
	}
	if plan.VehicleType != "" {
		set["vehicle_type"] = plan.VehicleType
	}
	if plan.Deposit != nil {
		set["deposit"] = *plan.Deposit
	}

	set["updated_at"] = primitive.NewDateTimeFromTime(time.Now())

	data, err := repo.UpdateOne(bson.M{"_id": idObject}, bson.M{"$set": set})
	if err != nil {
		return "", errors.New("error in updating plan")
	}
	return data, nil
}

func GetDeposit(city, planType string) (float64, error) {
	filter := bson.M{}
	if city != "" {
		filter["city"] = city
	}
	if planType != "" {
		filter["type"] = planType
	}
	dat, err := repo.Find(filter, bson.M{})
	if err != nil {
		return 0, err
	}
	if len(dat) == 0 {
		return 0, errors.New("no deposit found")
	}
	for _, v := range dat {
		if v.IsActive != nil && *v.IsActive && v.Deposit != nil {
			return *v.Deposit, nil
		}
	}
	return 0, nil
}
