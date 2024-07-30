package userattendance

import (
	utils "futureEVChronJobs/pkg/util"
	"net/http"

	"github.com/gorilla/mux"
)

func GetUserAttendanceHandler(w http.ResponseWriter, r *http.Request) {
	utils.SetOutput(w)
	data, err := GetUserAttendance()
	utils.SendOutput(err, w, r, data, nil, "GetUserAttendance")
}

func GetUserAttendanceByIDHandler(w http.ResponseWriter, r *http.Request) {
	utils.SetOutput(w)
	id := mux.Vars(r)["id"]
	data, err := GetUserAttendanceByID(id)
	utils.SendOutput(err, w, r, data, nil, "GetUserAttendanceByID")
}
