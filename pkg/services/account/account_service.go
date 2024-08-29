package account_service

import (
	"bikeRental/pkg/repo/profile"
	db "bikeRental/pkg/services/account/dbs"
	utils "bikeRental/pkg/util"
	"io"
	"log"

	"encoding/json"
	"net/http"
	"strings"
	"time"

	trestCommon "github.com/Trestx-technology/trestx-common-go-lib"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/bson"
)

var (
	accountService = db.NewSignUpService(profile.NewProfileRepository("users"))
)

// SignUp godoc
// @Summary SignUp
// @Description SignUp with the input payload
// @Tags SignUp
// @Accept  json
// @Produce  json
// @Param SignUp body db.Credentials true "SignUp"
// @Success 200
// @Router /register [post]
func SignUp(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	trestCommon.DLogMap("sign up email sent", logrus.Fields{
		"start_time": startTime})
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	user, err := GetCredentials(r)
	if user.Password == nil || *user.Password == "" {
		pass := trestCommon.GetRandomString(10)
		user.Password = new(string)
		user.Password = &pass
	}
	if err != nil {
		trestCommon.ECLog1(errors.Wrapf(err, "unable to get credentials"))
		w.WriteHeader(http.StatusUnsupportedMediaType)
		json.NewEncoder(w).Encode(bson.M{"status": false, "error": "Something Went wrong"})
		return

	}
	data, err := accountService.SignUp(user)
	if err != nil || data == "" {
		trestCommon.ECLog1(errors.Wrapf(err, "unable to sent singup email"))

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(bson.M{"status": false, "error": err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bson.M{"status": true, "error": "", "message": "sign up email sent successfully", "token": data})
	endTime := time.Now()
	duration := endTime.Sub(startTime)
	trestCommon.DLogMap("sign up email sent successfully", logrus.Fields{"duration": duration})
}

func LoginUsingPassword(w http.ResponseWriter, r *http.Request) {
	utils.SetOutput(w)

	user, err := GetCredentials(r)
	if err != nil {
		trestCommon.ECLog1(errors.Wrapf(err, "unable to parse credentials"))
		w.WriteHeader(http.StatusUnsupportedMediaType)
		json.NewEncoder(w).Encode(bson.M{"status": false, "error": "Something Went wrong"})
		return

	}
	data, token, err := accountService.LoginUsingPassword(user)
	if err != nil {
		if err.Error() == "user not verified" {
			trestCommon.ECLog1(errors.Wrapf(err, "unable to login"))
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(bson.M{"status": false, "error": "Email not Verified"})
			return
		}
		trestCommon.ECLog1(errors.Wrapf(err, "unable to login"))
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(bson.M{"status": false, "error": "invalid credentials"})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bson.M{"status": true, "error": "", "token": token, "data": data})
}
func LoginUsingPhone(w http.ResponseWriter, r *http.Request) {
	utils.SetOutput(w)

	user, err := GetCredentials(r)
	if err != nil {
		trestCommon.ECLog1(errors.Wrapf(err, "unable to parse credentials"))
		w.WriteHeader(http.StatusUnsupportedMediaType)
		if err.Error() != "user already logged in from another device" {
			json.NewEncoder(w).Encode(bson.M{"status": false, "error": "Something Went wrong"})
		} else {
			json.NewEncoder(w).Encode(bson.M{"status": false, "error": err.Error()})
		}
		return

	}

	response, err := accountService.LoginUsingPhone(user.PhoneNo, user.LastLoginDeviceID)
	if err != nil {
		trestCommon.ECLog1(errors.Wrapf(err, "unable to login"))
		w.WriteHeader(http.StatusOK)
		if err.Error() != "user already logged in from another device" {
			json.NewEncoder(w).Encode(bson.M{"status": false, "error": "Something Went wrong"})
		} else {
			json.NewEncoder(w).Encode(bson.M{"status": false, "error": err.Error()})
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bson.M{"status": true, "error": "", "data": response})
}
func ForgetPasswordOTPLink(w http.ResponseWriter, r *http.Request) {
	utils.SetOutput(w)

	user, err := GetCredentials(r)
	if err != nil {
		trestCommon.ECLog1(errors.Wrapf(err, "unable to parse credentials"))
		w.WriteHeader(http.StatusUnsupportedMediaType)
		json.NewEncoder(w).Encode(bson.M{"status": false, "error": "Something Went wrong"})
		return

	}
	data, err := accountService.SendEmailOTP(user.Email)
	utils.SendOutput(err, w, r, data, nil, "ForgetPasswordOTPLink")
}
func VerifyAndUpdatePassword(w http.ResponseWriter, r *http.Request) {
	utils.SetOutput(w)

	user, err := GetCredentials(r)
	if err != nil {
		trestCommon.ECLog1(errors.Wrapf(err, "unable to parse credentials"))
		w.WriteHeader(http.StatusUnsupportedMediaType)
		json.NewEncoder(w).Encode(bson.M{"status": false, "error": "Something Went wrong"})
		return

	}
	data, _, err := accountService.VerifyResetLink(user)
	utils.SendOutput(err, w, r, data, nil, "VerifyAndUpdatePassword")
}
func VerifyOTPAndSendToken(w http.ResponseWriter, r *http.Request) {
	utils.SetOutput(w)

	user, err := GetCredentials(r)
	if err != nil {
		trestCommon.ECLog1(errors.Wrapf(err, "unable to parse credentials"))
		w.WriteHeader(http.StatusUnsupportedMediaType)
		json.NewEncoder(w).Encode(bson.M{"status": false, "error": "Something Went wrong"})
		return

	}
	if user.PhoneOTP == nil {
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(bson.M{"status": false, "error": "Wrong OTP"})
		return
	}
	token := ""
	log.Println(user.FirebaseToken)
	if user.FirebaseToken != nil {
		token = *user.FirebaseToken
	}
	data, token, err := accountService.VerifyOTP(user.PhoneNo, *user.PhoneOTP, token, user.LastLoginDeviceID)
	if err != nil {
		if err.Error() == "user not verified" {
			trestCommon.ECLog1(errors.Wrapf(err, "unable to login"))
			w.WriteHeader(http.StatusAccepted)
			json.NewEncoder(w).Encode(bson.M{"status": false, "error": "Wrong OTP"})
			return
		}
		trestCommon.ECLog1(errors.Wrapf(err, "unable to login"))
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(bson.M{"status": false, "error": "Wrong OTP"})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bson.M{"status": true, "error": "", "data": data, "token": token})
}
func ChangePassword(w http.ResponseWriter, r *http.Request) {
	utils.SetOutput(w)
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
	user, err := GetCredentials(r)
	if err != nil {
		trestCommon.ECLog1(errors.Wrapf(err, "unable to parse credentials"))
		w.WriteHeader(http.StatusUnsupportedMediaType)
		json.NewEncoder(w).Encode(bson.M{"status": false, "error": "Something Went wrong"})
		return

	}
	data, err := accountService.ChangePassword(user)
	utils.SendOutput(err, w, r, data, user, "ChangePassword")
}

func GetCredentials(r *http.Request) (db.Credentials, error) {
	var user db.Credentials

	body, _ := io.ReadAll(r.Body)
	err := json.Unmarshal(body, &user)
	if err != nil {
		return user, err
	}
	user.Email = strings.TrimSpace(user.Email)
	return user, err
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	utils.SetOutput(w)
	id := mux.Vars(r)["id"]
	data, err := accountService.Logout(id)
	utils.SendOutput(err, w, r, data, nil, "Logout")
}
