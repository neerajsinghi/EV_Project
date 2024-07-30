package userattendance

import (
	utils "bikeRental/pkg/util"
	"net/http"

	"github.com/gorilla/mux"
)

func GetUserAttendanceHandler(w http.ResponseWriter, r *http.Request) {
	if utils.SetOutput(w, r) {
		return
	}
	data, err := GetUserAttendance()
	utils.SendOutput(err, w, r, data, nil, "GetUserAttendance")
}

func GetUserAttendanceByIDHandler(w http.ResponseWriter, r *http.Request) {
	if utils.SetOutput(w, r) {
		return
	}
	id := mux.Vars(r)["id"]
	data, err := GetUserAttendanceByID(id)
	utils.SendOutput(err, w, r, data, nil, "GetUserAttendanceByID")
}
