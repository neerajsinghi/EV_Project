package reffer

import (
	"encoding/json"
	"net/http"
	"time"

	trestCommon "github.com/Trestx-technology/trestx-common-go-lib"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

type Referral struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func GetReferralsHandler(w http.ResponseWriter, r *http.Request) {
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
	data, err := GetReffer()
	if err != nil {
		trestCommon.ECLog1(err)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(bson.M{"status": false, "error": err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bson.M{"status": true, "error": "", "data": data})

}
