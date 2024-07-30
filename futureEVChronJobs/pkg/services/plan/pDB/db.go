package pdb

import (
	"futureEVChronJobs/pkg/entity"
	"futureEVChronJobs/pkg/repo/plan"

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
