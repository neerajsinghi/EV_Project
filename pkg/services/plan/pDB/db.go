package pdb

import (
	"bikeRental/pkg/entity"
	"bikeRental/pkg/repo/plan"
	"errors"
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
		dat, err := s.GetPlans("rental", plan.City)
		if err == nil && len(dat) > 0 {
			for _, v := range dat {
				if dat[0].IsActive != nil && *dat[0].IsActive && v.VehicleType == plan.VehicleType {
					if plan.Validity == v.Validity {
						return "", errors.New("rental plan already exist for this city")
					}
				}
			}
		}
	}

	plan.IsActive = new(bool)
	*plan.IsActive = true
	if plan.Deposit == nil {
		plan.Deposit = new(float64)
		*plan.Deposit = 0
	}
	return repo.InsertOne(plan)
}

// DeletePlan implements PlanI.
func (s *service) DeletePlan(id string) error {
	idObject, _ := primitive.ObjectIDFromHex(id)
	return repo.DeleteOne(bson.M{"_id": idObject})
}

// GetPlan implements PlanI.
func (s *service) GetPlan(id string) (entity.PlanDB, error) {
	idObject, _ := primitive.ObjectIDFromHex(id)
	return repo.FindOne(bson.M{"_id": idObject}, bson.M{})

}

// GetPlans implements PlanI.
func (s *service) GetPlans(pType, city string) ([]entity.PlanDB, error) {
	filter := bson.M{}
	if pType != "" {
		filter["type"] = pType
	}
	if city != "" {
		filter["city"] = city
	}
	return repo.Find(filter, bson.M{})

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

	return repo.UpdateOne(bson.M{"_id": idObject}, bson.M{"$set": set})

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
