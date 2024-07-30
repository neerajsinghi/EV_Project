package users

import (
	"bikeRental/pkg/entity"
	"bikeRental/pkg/services/users/udb"
	utils "bikeRental/pkg/util"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

var service = udb.NewService()

func GetUsers(w http.ResponseWriter, r *http.Request) {
	if utils.SetOutput(w, r) {
		return
	}
	userType := r.URL.Query().Get("type")
	data, err := service.GetUsers(userType)
	utils.SendOutput(err, w, r, data, nil, "GetUsers")
}

func GetUserById(w http.ResponseWriter, r *http.Request) {
	if utils.SetOutput(w, r) {
		return
	}
	id := mux.Vars(r)["id"]
	data, err := service.GetUserById(id)
	utils.SendOutput(err, w, r, data, nil, "GetUserById")
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	if utils.SetOutput(w, r) {
		return
	}
	id := mux.Vars(r)["id"]
	user, err := getUser(r)
	shouldReturn := utils.CheckError(err, w)
	if shouldReturn {
		return
	}
	data, err := service.UpdateUser(id, user)
	utils.SendOutput(err, w, r, data, user, "UpdateUser")

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
	if utils.SetOutput(w, r) {
		return
	}
	id := mux.Vars(r)["id"]
	_, err := service.DeleteUser(id)
	utils.SendOutput(err, w, r, nil, nil, "DeleteUser")
}

func DeleteUserPermanently(w http.ResponseWriter, r *http.Request) {
	if utils.SetOutput(w, r) {
		return
	}
	id := mux.Vars(r)["id"]
	err := service.DeleteUserPermanently(id)
	utils.SendOutput(err, w, r, nil, nil, "DeleteUserPermanently")
}

func RemovePlan(w http.ResponseWriter, r *http.Request) {
	if utils.SetOutput(w, r) {
		return
	}
	id := mux.Vars(r)["id"]
	planID := mux.Vars(r)["plan_id"]
	data, err := service.RemovePlan(id, planID)
	utils.SendOutput(err, w, r, data, nil, "RemovePlan")
}
