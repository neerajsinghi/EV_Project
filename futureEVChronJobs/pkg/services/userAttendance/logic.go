package userattendance

import (
	"futureEVChronJobs/pkg/entity"
	"futureEVChronJobs/pkg/repo/attendance"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var repo = attendance.NewRepository("userAttendance")

func AddUserAttendance(userAttendance entity.UserAttendance) (string, error) {
	return repo.InsertOne(userAttendance)
}

func GetUserAttendance() ([]entity.UserAttendance, error) {
	return repo.Find(nil, nil)
}

func GetUserAttendanceByID(id string) ([]entity.UserAttendance, error) {
	idO, _ := primitive.ObjectIDFromHex(id)
	return repo.Find(bson.M{"_id": idO}, nil)
}
