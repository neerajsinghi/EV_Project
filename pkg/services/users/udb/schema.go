package udb

import "bikeRental/pkg/entity"

type UserI interface {
	GetUsers(typeU string) ([]entity.ProfileOut, error)
	GetUserById(id string) (entity.ProfileOut, error)
	UpdateUser(id string, user entity.ProfileDB) (string, error)
	DeleteUser(id string) (string, error)
	DeleteUserPermanently(id string) error
	RemovePlan(id, planId string) (string, error)
}
