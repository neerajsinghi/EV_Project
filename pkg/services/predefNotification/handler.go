package predefnotification

import (
	"bikeRental/pkg/entity"
	utils "bikeRental/pkg/util"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func GetPredef(w http.ResponseWriter, r *http.Request) {
	utils.SetOutput(w)
	name := r.URL.Query().Get("name")
	var data interface{}
	var err error
	if name == "" {
		data, err = GetAll()
	} else {
		data, err = Get(name)
	}
	utils.SendOutput(err, w, r, data, "GetPreddef")
}
func AddPredef(w http.ResponseWriter, r *http.Request) {
	utils.SetOutput(w)
	notification, err := parseNotification(r)
	if err != nil {
		utils.SendOutput(err, w, r, nil, "AddPredef")
		return
	}
	data, err := InsertOne(notification)
	utils.SendOutput(err, w, r, data, "AddPredef")
}
func UpdatePredef(w http.ResponseWriter, r *http.Request) {
	utils.SetOutput(w)
	name := r.URL.Query().Get("name")
	if name == "" {
		utils.SendOutput(errors.New("name missing"), w, r, nil, "UpdatePredef")
		return
	}
	notification, err := parseNotification(r)
	if err != nil {
		utils.SendOutput(err, w, r, nil, "UpdatePredef")
		return
	}
	data, err := UpdateOne(name, notification)
	utils.SendOutput(err, w, r, data, "UpdatePredef")
}
func DeletePredef(w http.ResponseWriter, r *http.Request) {
	utils.SetOutput(w)
	name := r.URL.Query().Get("name")
	if name == "" {
		utils.SendOutput(errors.New("name missing"), w, r, nil, "DeletePredef")
		return
	}
	err := DeleteOne(name)
	utils.SendOutput(err, w, r, nil, "DeletePredef")
}
func parseNotification(r *http.Request) (entity.PreDefNotification, error) {
	var notification entity.PreDefNotification
	body, _ := io.ReadAll(r.Body)
	err := json.Unmarshal(body, &notification)
	return notification, err
}
