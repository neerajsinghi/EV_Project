package iotbike

import (
	db "bikeRental/pkg/services/iotBike/db"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	trestCommon "github.com/Trestx-technology/trestx-common-go-lib"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

var service = db.NewService()

// FindAll returns all the bikes
func GetAll(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	trestCommon.DLogMap("setting brand", logrus.Fields{
		"start_time": startTime})
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	tokenString := strings.Split(r.Header.Get("Authorization"), " ")
	if len(tokenString) < 2 {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(bson.M{"status": false, "error": "authorization failed"})
		return
	}
	_, err := trestCommon.DecodeToken(tokenString[1])
	if err != nil {
		trestCommon.ECLog1(errors.Wrapf(err, "failed to authenticate token"))
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(bson.M{"status": false, "error": "authorization failed"})
		return
	}
	data, err := service.FindAll()
	if err != nil {
		trestCommon.ECLog1(errors.Wrapf(err, "unable to set brand"))

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(bson.M{"status": false, "error": "Unable to find brand"})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bson.M{"status": true, "error": "", "data": data})
	endTime := time.Now()
	duration := endTime.Sub(startTime)
	trestCommon.DLogMap("brand updated", logrus.Fields{
		"duration": duration,
	})

}

// get nearest based on lat and long
func GetNearest(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	defer func() {
		trestCommon.DLogMap("brand updated", logrus.Fields{
			"duration": time.Since(startTime),
		})
	}()
	trestCommon.DLogMap("setting brand", logrus.Fields{
		"start_time": startTime})
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
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
	if err != nil {
		trestCommon.ECLog1(errors.Wrapf(err, "unable to set brand"))

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(bson.M{"status": false, "error": "Unable to find brand"})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bson.M{"status": true, "error": "", "data": data})

}
