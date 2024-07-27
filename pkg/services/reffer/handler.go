package reffer

import (
	utils "bikeRental/pkg/util"
	"net/http"
)

type Referral struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func GetReferralsHandler(w http.ResponseWriter, r *http.Request) {
	utils.SetOutput(w)
	data, err := GetReffer()
	utils.SendOutput(err, w, r, data, nil, "GetReferrals")
}
