package pdb

import "futureEVChronJobs/pkg/entity"

type PlanI interface {
	GetPlans(pType, city string) ([]entity.PlanDB, error)
	GetPlan(id string) (entity.PlanDB, error)
}
