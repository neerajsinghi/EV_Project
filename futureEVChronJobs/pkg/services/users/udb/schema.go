package udb

import "futureEVChronJobs/pkg/entity"

type UserI interface {
	GetUsers(typeU string) ([]entity.ProfileOut, error)
	GetUserById(id string) (entity.ProfileOut, error)
}
