package pdb

import "bikeRental/pkg/entity"

type PlanI interface {
	AddPlan(plan entity.PlanDB) (string, error)
	GetPlan(id string) (entity.PlanDB, error)
	GetPlans(pType, city string) ([]entity.PlanDB, error)
	GetPlansAdmin(pType, city string) ([]entity.PlanDB, error)
	UpdatePlan(id string, plan entity.PlanDB) (string, error)
	DeletePlan(id string) error
}
