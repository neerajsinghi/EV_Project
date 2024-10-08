package statistics

import (
	utils "bikeRental/pkg/util"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// Logic is a function that returns the statistics of the users, stations and chargers.
func Statistics(w http.ResponseWriter, r *http.Request) {
	utils.SetOutput(w)
	startDate := r.URL.Query().Get("startDate")
	endDate := r.URL.Query().Get("endDate")
	city := r.URL.Query().Get("city")
	service := r.URL.Query().Get("service")
	start, _ := time.Parse(time.RFC3339, startDate)

	end, _ := time.Parse(time.RFC3339, endDate)

	data, err := Logic(start, end, city, service)
	utils.SendOutput(err, w, r, data, nil, "Statistics")
}

func GetVehicleDataHand(w http.ResponseWriter, r *http.Request) {
	utils.SetOutput(w)
	id := mux.Vars(r)["id"]
	idi, _ := strconv.Atoi(id)
	data, err := GetVehicleData(idi)
	utils.SendOutput(err, w, r, data, nil, "GetVehicleData")
}

func ImmobilizeDevHand(w http.ResponseWriter, r *http.Request) {
	utils.SetOutput(w)
	id := mux.Vars(r)["id"]
	idi, _ := strconv.Atoi(id)
	err := ImmobilizeDevice(idi)
	utils.SendOutput(err, w, r, "immobilized successfully", id, "ImmobilizeDevice")
}
