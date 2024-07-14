package users

import (
	"bikeRental/pkg/entity"
	"bikeRental/pkg/services/users/udb"
	"encoding/json"
	"net/http"
	"time"

	trestCommon "github.com/Trestx-technology/trestx-common-go-lib"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

var service = udb.NewService()

func GetUsers(w http.ResponseWriter, r *http.Request) {
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
	userType := r.URL.Query().Get("type")
	data, err := service.GetUsers(userType)
	if err != nil {
		trestCommon.ECLog1(err)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(bson.M{"status": false, "error": "Unable to get plan"})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bson.M{"status": true, "error": "", "data": data})
}

func GetUserById(w http.ResponseWriter, r *http.Request) {
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
	id := mux.Vars(r)["id"]
	data, err := service.GetUserById(id)
	if err != nil {
		trestCommon.ECLog1(err)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(bson.M{"status": false, "error": "Unable to get plan"})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bson.M{"status": true, "error": "", "data": data})
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
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
	id := mux.Vars(r)["id"]
	user, err := getUser(r)
	if err != nil {
		trestCommon.ECLog1(err)
		w.WriteHeader(http.StatusUnsupportedMediaType)
		json.NewEncoder(w).Encode(bson.M{"status": false, "error": "Something Went wrong"})
		return
	}
	data, err := service.UpdateUser(id, user)
	if err != nil {
		trestCommon.ECLog1(err)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(bson.M{"status": false, "error": "Unable to update user"})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bson.M{"status": true, "error": "", "data": data})
}

func getUser(r *http.Request) (entity.ProfileDB, error) {
	var user entity.ProfileDB
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		return entity.ProfileDB{}, err
	}
	return user, nil
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
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
	id := mux.Vars(r)["id"]
	_, err := service.DeleteUser(id)
	if err != nil {
		trestCommon.ECLog1(err)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(bson.M{"status": false, "error": "Unable to delete user"})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bson.M{"status": true, "error": ""})
}
