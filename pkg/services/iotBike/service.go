package iotbike

import (
	db "bikeRental/pkg/services/iotBike/db"
	utils "bikeRental/pkg/util"
	"encoding/json"
	"net/http"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
)

var service = db.NewService()

// FindAll returns all the bikes
func GetAll(w http.ResponseWriter, r *http.Request) {

	utils.SetOutput(w)
	data, err := service.FindAll()
	utils.SendOutput(err, w, r, data, "GetAll")

}

// get nearest based on lat and long
func GetNearest(w http.ResponseWriter, r *http.Request) {

	utils.SetOutput(w)
	lat := r.URL.Query().Get("lat")
	long := r.URL.Query().Get("long")
	distance := r.URL.Query().Get("distance")
	bType := r.URL.Query().Get("type")
	if lat == "" || long == "" {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(bson.M{"status": false, "error": "lat and long are required"})
		return
	}
	latFloat, _ := strconv.ParseFloat(lat, 64)
	longFloat, _ := strconv.ParseFloat(long, 64)
	distantInt, _ := strconv.Atoi(distance)
	data, err := service.FindNearByBikes(latFloat, longFloat, distantInt, bType)
	utils.SendOutput(err, w, r, data, "GetNearest")

}
