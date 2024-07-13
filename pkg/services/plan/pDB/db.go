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
	if plan.Price == 0 {
		return "", errors.New("price is required")
	}
	plan.IsActive = new(bool)
	*plan.IsActive = true
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

	set["updated_at"] = primitive.NewDateTimeFromTime(time.Now())

	return repo.UpdateOne(bson.M{"_id": idObject}, bson.M{"$set": set})

}
