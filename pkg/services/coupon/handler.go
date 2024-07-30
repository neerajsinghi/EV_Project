package coupon

import (
	"bikeRental/pkg/entity"
	"bikeRental/pkg/services/coupon/cdb"
	utils "bikeRental/pkg/util"
	"encoding/json"
	"net/http"

	trestCommon "github.com/Trestx-technology/trestx-common-go-lib"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
)

var coupon = cdb.NewCoupon()

// AddCoupon adds a new coupon
func AddCoupon(w http.ResponseWriter, r *http.Request) {
	utils.SetOutput(w,r)
	couponD, err := getCoupon(r)
	if err != nil {
		trestCommon.ECLog1(errors.Wrapf(err, "unable to get coupon"))
		w.WriteHeader(http.StatusUnsupportedMediaType)
		json.NewEncoder(w).Encode(bson.M{"status": false, "error": "Something Went wrong"})
		return
	}
	data, err := coupon.AddCoupon(couponD)
	utils.SendOutput(err, w, r, data, couponD, "AddCoupon")
}

// UpdateCoupon updates a coupon
func UpdateCoupon(w http.ResponseWriter, r *http.Request) {
	utils.SetOutput(w,r)
	couponD, err := getCoupon(r)
	if err != nil {
		trestCommon.ECLog1(errors.Wrapf(err, "unable to get coupon"))
		w.WriteHeader(http.StatusUnsupportedMediaType)
		json.NewEncoder(w).Encode(bson.M{"status": false, "error": "Something Went wrong"})
		return
	}
	id := mux.Vars(r)["id"]
	data, err := coupon.UpdateCoupon(id, couponD)
	utils.SendOutput(err, w, r, data, couponD, "UpdateCoupon")
}

// DeleteCoupon deletes a coupon
func DeleteCoupon(w http.ResponseWriter, r *http.Request) {
	utils.SetOutput(w,r)
	id := mux.Vars(r)["id"]
	err := coupon.DeleteCoupon(id)
	utils.SendOutput(err, w, r, "Deleted successfully", nil, "DeleteCoupon")
}

// GetCoupon gets a coupon
func GetCoupons(w http.ResponseWriter, r *http.Request) {
	utils.SetOutput(w,r)
	data, err := coupon.GetCoupon()
	utils.SendOutput(err, w, r, data, nil, "GetCoupon")
}

func getCoupon(r *http.Request) (entity.CouponDB, error) {
	var coupon entity.CouponDB
	err := json.NewDecoder(r.Body).Decode(&coupon)
	if err != nil {
		return entity.CouponDB{}, errors.Wrapf(err, "unable to decode request body")
	}
	return coupon, nil
}
