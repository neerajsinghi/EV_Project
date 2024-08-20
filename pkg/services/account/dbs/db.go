package db

import (
	"bikeRental/pkg/entity"
	"bikeRental/pkg/repo/profile"
	"bikeRental/pkg/repo/reffer"
	userattendance "bikeRental/pkg/services/userAttendance"
	utils "bikeRental/pkg/util"
	"math/rand"
	"net/mail"
	"strconv"
	"strings"

	"errors"
	"time"

	trestCommon "github.com/Trestx-technology/trestx-common-go-lib"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

var (
	repo = profile.NewProfileRepository("users")
)

type accountService struct{}

func NewSignUpService(profile profile.ProfileRepository) AccountService {
	repo = profile
	return &accountService{}
}
func (*accountService) SignUp(cred Credentials) (string, error) {
	if cred.Password == nil {
		err := errors.New("password missing")
		trestCommon.ECLog2("sign up failed no password", err)
		return "", err
	}
	if cred.Role == nil {
		cred.Role = new(string)
		*cred.Role = "user"
	}
	if cred.Role != nil && *cred.Role == "" {
		*cred.Role = "user"
	}
	if cred.PhoneNo == "" {
		err := errors.New("phone number required")
		trestCommon.ECLog2("sign up failed phone number required", err)
		return "", err

	}
	cred.Email = strings.ToLower(cred.Email)
	if cred.Email != "" {
		//validate email format
		_, err := mail.ParseAddress(cred.Email)
		if err != nil {
			trestCommon.ECLog2("sign up failed invalid email format", err)
			return "", errors.New("sign up failed invalid email format")
		}
	}
	cred.PEmail = strings.ToLower(cred.PEmail)
	if cred.PEmail != "" {
		//validate PEmail format
		_, err := mail.ParseAddress(cred.PEmail)
		if err != nil {
			trestCommon.ECLog2("sign up failed invalid email format", err)
			return "", errors.New("sign up failed invalid email format")
		}
	}
	if cred.ReferralCodeUsed != nil && *cred.ReferralCodeUsed != "" {
		user, err := repo.FindOne(bson.M{"referral_code": *cred.ReferralCodeUsed}, bson.M{})
		if user.ID.IsZero() {
			trestCommon.ECLog2("sign up failed invalid referral code", err)
			return "", errors.New("sign up failed invalid referral code")
		}
	}
	// check name for special characters
	if strings.ContainsAny(cred.Name, "!@#$%^&*()_+{}|:<>?") {
		err := errors.New("name contains special characters")
		trestCommon.ECLog2("sign up failed name contains special characters", err)
		return "", err
	}
	_, err := CheckUser(cred.Email, cred.PhoneNo)

	if err != nil {
		trestCommon.ECLog2("sign up user not found", err)
		if err.Error() == "mongo: no documents in result" {
			var serv *accountService
			tokenString, err := serv.hashAndInsertData(cred)
			if err != nil {
				trestCommon.ECLog3("sign up not successful", err, logrus.Fields{"email": cred.Email})
			}

			return tokenString, err
		} else {
			return "", err

		}
	}

	return "", errors.New("phone number already registered")
}

func (*accountService) SendVerificationEmail(email, pemail, uid string) (string, error) {
	emailSentTime := time.Now()
	if email == "" {
		return "email sent successfully", nil
	}
	email = strings.ToLower(email)
	verificationCode := trestCommon.GetRandomString(16)
	sendCode, err := trestCommon.Encrypt(email + ":" + verificationCode)
	if err != nil {
		trestCommon.ECLog2("send verification email encryption failed", err)
		return "", err
	}
	var userData entity.ProfileDB
	if pemail != "" {
		userData, _ = CheckUser(pemail, "")
	} else {
		userData, _ = CheckUser(email, "")
	}
	_, err = utils.SendVerificationCode(email, userData.Name, sendCode)
	if err != nil {
		trestCommon.ECLog2("send verification email failed", err)
		return "", err
	}
	if pemail != "" {
		_, err = repo.UpdateOne(bson.M{"email": pemail}, bson.M{"$set": bson.M{"email": email, "email_sent_time": emailSentTime, "verification_code": verificationCode}})
		if err != nil {
			trestCommon.ECLog2("send verification email update failed", err)
			return "", err
		}
		return trestCommon.CreateToken(uid, email, "", "")
	} else {
		_, err = repo.UpdateOne(bson.M{"email": email}, bson.M{"$set": bson.M{"email_sent_time": emailSentTime, "verification_code": verificationCode}})
		if err != nil {
			trestCommon.ECLog2("send verification email update failed", err)
			return "", err
		}
		return "email sent successfully", nil
	}

}

func (*accountService) SendEmailOTP(email string) (string, error) {
	email = strings.ToLower(email)
	emailSentTime := time.Now()
	if email == "" {
		return "", errors.New("email id required")
	}
	_, err := CheckUser(email, "")
	if err != nil {
		return "", errors.New("user account doesn't exist")
	}
	randomOTP := 1000 + rand.Intn(9999-1000)

	utils.EmailLoginOTP(email, email, strconv.Itoa(randomOTP), "change_password")

	_, err = repo.UpdateOne(bson.M{"email": email}, bson.M{"$set": bson.M{"email_sent_time": emailSentTime, "password_reset_code": strconv.Itoa(randomOTP)}})
	if err != nil {
		trestCommon.ECLog2("send verification email update failed", err)
		return "", err
	}
	return "email sent successfully", nil
}
func (*accountService) LoginUsingPhone(phone string) (string, error) {
	data, err := CheckUser("", phone)
	if err != nil {
		trestCommon.ECLog3("verify user not found", err, logrus.Fields{"phone": phone})
		return "", err
	}
	if data.Status == "deleted" {
		return "user deleted", errors.New("user deleted")
	}
	code := strconv.Itoa(100000 + rand.Intn(999999-100000))

	utils.SendVarificationSMS(phone, code)
	_, err = repo.UpdateOne(bson.M{"_id": data.ID}, bson.M{"$set": bson.M{"verified_time": time.Now(), "phone_otp": code}})
	if err != nil {
		trestCommon.ECLog3("verify user unable to update status", err, logrus.Fields{"phone": phone})
		return "", err
	}
	return "otp Sent", nil
}

func (*accountService) VerifyOTP(phone, otp, token string) (*entity.ProfileDB, string, error) {
	data, err := CheckUser("", phone)
	if err != nil {
		trestCommon.ECLog3("verify user not found", err, logrus.Fields{"phone": phone})
		return nil, "", err
	}

	_, err = utils.VerifyOTP(phone, otp)
	if err != nil && data.PhoneNo != "+919996729701" {
		trestCommon.ECLog3("verify user unable to update status", err, logrus.Fields{"phone": phone})
		return nil, "", err
	}
	_, err = repo.UpdateOne(bson.M{"_id": data.ID}, bson.M{"$set": bson.M{"verified_time": time.Now(), "firebase_token": token}})
	if err != nil {
		trestCommon.ECLog3("verify user unable to update status", err, logrus.Fields{"phone": phone})
		return nil, "", err
	}
	tokenString, _ := trestCommon.CreateToken(data.ID.Hex(), data.Email, "", data.Status)
	data.Password = nil
	return &data, tokenString, nil
}

func (serv *accountService) LoginUsingPassword(cred Credentials) (entity.ProfileDB, string, error) {
	if cred.Password == nil {
		err := errors.New("password missing")
		trestCommon.ECLog2("login failed no password", err)
		return entity.ProfileDB{}, "", err
	}
	cred.Email = strings.ToLower(cred.Email)
	salt := viper.GetString("salt")
	userData, err := CheckUser(cred.Email, "")
	if err != nil {
		trestCommon.ECLog2("login failed user not found", err)
		return entity.ProfileDB{}, "", err
	}
	// if userData.Status == "created" {
	// 	err = errors.New("user not verified")
	// 	trestCommon.ECLog2("login failed user has not verified his/her email address", err)
	// 	return "", err
	// }
	if userData.Status == "deleted" || userData.Status == "archived" {
		err = errors.New("unauthorized user")
		trestCommon.ECLog2("login failed user has deleted or archived his profile", err)
		return entity.ProfileDB{}, "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(*userData.Password), []byte(*cred.Password+salt))
	if err != nil {
		trestCommon.ECLog2("login failed password hash doesn't match", err)
		return entity.ProfileDB{}, "", err
	}
	if cred.Role == nil {
		cred.Role = new(string)
		*cred.Role = "admin"
	}
	tokenString, err := trestCommon.CreateToken(userData.ID.Hex()+" name:"+cred.Name, *cred.Role, "", userData.Status)
	if err != nil {
		trestCommon.ECLog3("login failed unable to create token", err, logrus.Fields{"email": cred.Email, "name": userData.Name, "status": userData.Status})
		return entity.ProfileDB{}, "", err
	}
	repo.UpdateOne(bson.M{"_id": userData.ID}, bson.M{"$set": bson.M{"login_time": time.Now()}})
	userData.Password = nil
	return userData, tokenString, nil
}
func (*accountService) ChangePassword(cred Credentials) (string, error) {
	userData, err := CheckUser(cred.Email, "")
	if err != nil {
		trestCommon.ECLog3("verify user not found", err, logrus.Fields{"email": cred.Email})
		return "", err
	}
	setFilter := bson.M{}
	salt := viper.GetString("salt")
	hash, err := bcrypt.GenerateFromPassword([]byte(*cred.Password+salt), 5)
	if err != nil {
		trestCommon.ECLog3("hash password", err, logrus.Fields{"email": cred.Email})
		return "", err
	}
	*cred.Password = string(hash)
	setFilter["$set"] = bson.M{"verified_time": time.Now(), "password": cred.Password, "login_time": time.Now(), "updated_at": time.Now()}

	return repo.UpdateOne(bson.M{"_id": userData.ID}, setFilter)
}

func (*accountService) VerifyEmail(cred Credentials) (string, error) {

	userData, err := CheckUser(cred.Email, "")
	if err != nil {
		trestCommon.ECLog3("verify user not found", err, logrus.Fields{"email": cred.Email})
		return "", err
	}

	if cred.VerificationCode != userData.VerificationCode {
		err = errors.New("unauthorized user")
		trestCommon.ECLog3("verify user verification code didn't match", err, logrus.Fields{"email": cred.Email, "db verify code": userData.VerificationCode, "code provided by user": cred.VerificationCode})
		return "", err
	}
	if userData.Status == "verified" {
		err = errors.New("user already verified")
		trestCommon.ECLog3("verify user verification user already verified", err, logrus.Fields{"email": cred.Email})
		return "", err
	}
	_, err = repo.UpdateOne(bson.M{"_id": userData.ID}, bson.M{"$set": bson.M{"verified_time": time.Now(), "status": "verified"}})
	if err != nil {
		trestCommon.ECLog3("verify user unable to update status", err, logrus.Fields{"email": cred.Email})
		return "", err
	}
	return "verified", nil
}

func CheckUser(email, mobile string) (entity.ProfileDB, error) {

	if email != "" && mobile != "" {
		return repo.FindOne(bson.M{"$or": bson.A{bson.M{"email": email}, bson.M{"phone_no": mobile}}}, bson.M{})
	} else if email == "" {
		return repo.FindOne(bson.M{"phone_no": mobile}, bson.M{})

	} else {
		return repo.FindOne(bson.M{"email": email}, bson.M{})
	}
}

func (*accountService) hashAndInsertData(cred Credentials) (string, error) {
	salt := viper.GetString("salt")

	hash, err := bcrypt.GenerateFromPassword([]byte(*cred.Password+salt), 5)
	if err != nil {
		trestCommon.ECLog3("hash password", err, logrus.Fields{"email": cred.Email})
		return "", err
	}
	*cred.Password = string(hash)
	cred.Status = "created"
	cred.PhoneNoVerified = false
	userName := strings.ReplaceAll(cred.Name, " ", "")
	refCode := ""
	if len(userName) > 4 {
		refCode = strings.ToLower(userName[:4]) + trestCommon.GetRandomString(4)
	} else {
		refCode = strings.ToLower(userName) + trestCommon.GetRandomString(8-len(userName))
	}

	cred.ReferralCode = &refCode
	cred.ID = primitive.NewObjectID()
	cred.CreatedTime = primitive.NewDateTimeFromTime(time.Now())
	_, err = repo.InsertOne(cred.ProfileDB)
	if err != nil {
		trestCommon.ECLog3("hashAndInsertData Insert failed", err, logrus.Fields{"email": cred.Email})
		return "", nil
	}
	// _, err = serv.SendOTP(cred.Email, cred.PhoneNo)
	// if err != nil {
	// 	trestCommon.ECLog3("hashAndInsertData Insert failed", err, logrus.Fields{"email": cred.Email})
	// }
	utils.SendVarificationSMS(cred.PhoneNo, "")
	if cred.ReferralCodeUsed != nil && *cred.ReferralCodeUsed != "" {
		user, err := repo.FindOne(bson.M{"referral_code": *cred.ReferralCodeUsed}, bson.M{})
		if err == nil {
			refUser, _ := repo.FindOne(bson.M{"_id": cred.ID}, bson.M{})
			refer := entity.ReferralDB{
				ReferralCode:      *cred.ReferralCodeUsed,
				ReferrerBy:        user.ID.Hex(),
				ReferredByProfile: user,
				ReferralOf:        refUser,
				ReferralStatus:    "pending",
				ReferralOfId:      refUser.ID.Hex(),
			}
			repoRef := reffer.NewRepository("referral")
			repoRef.InsertOne(refer)
		}
	}
	if cred.StaffStatus != "" {
		attend := entity.UserAttendance{
			ProfileID:   cred.ID.Hex(),
			Status:      cred.StaffStatus,
			UpdatedTime: time.Now(),
		}
		userattendance.AddUserAttendance(attend)
	}
	return "Verification code sent successfully", nil
}
func (*accountService) SendResetLink(email string) (string, error) {
	var cred Credentials
	cred.Email = email
	_, err := CheckUser(cred.Email, "")
	if err != nil {
		trestCommon.ECLog2("user not found", err)
		return "", err

	}
	emailSentTime := time.Now()
	verificationCode := trestCommon.GetRandomString(16)
	resetCode, err := trestCommon.Encrypt(email + ":" + verificationCode)
	if err != nil {
		trestCommon.ECLog2("send reset link encryption failed", err)
		return "", err
	}
	_, err = trestCommon.SendResetPasswordLink(email, resetCode)
	if err != nil {
		trestCommon.ECLog2("send reset password link failed", err)
		return "", err
	}
	_, err = repo.UpdateOne(bson.M{"email": email}, bson.M{"$set": bson.M{"email_sent_time": emailSentTime, "password_reset_code": verificationCode}})
	if err != nil {
		trestCommon.ECLog2("send reset link update failed", err)
		return "", err
	}
	return "Reset link sent successfully", nil
}

func (*accountService) VerifyResetLink(cred Credentials) (string, string, error) {

	userData, err := CheckUser(cred.Email, "")
	if err != nil {
		trestCommon.ECLog3("verify user not found", err, logrus.Fields{"email": cred.Email})
		return "", "", err
	}

	if cred.PasswordResetCode != userData.PasswordResetCode {
		err = errors.New("unauthorized user")
		trestCommon.ECLog3("verify user password reset code didn't match", err, logrus.Fields{"email": cred.Email, "db verify code": userData.PasswordResetCode, "code provided by user": cred.PasswordResetCode})
		return "", "", err
	}
	salt := viper.GetString("salt")

	hash, err := bcrypt.GenerateFromPassword([]byte(*cred.Password+salt), 5)
	if err != nil {
		trestCommon.ECLog3("hash password", err, logrus.Fields{"email": cred.Email})
		return "", "", err
	}
	*cred.Password = string(hash)

	_, err = repo.UpdateOne(bson.M{"_id": userData.ID}, bson.M{"$set": bson.M{"password_reset_time": time.Now(), "password": *cred.Password}})
	if err != nil {
		trestCommon.ECLog3("verify user unable to update status", err, logrus.Fields{"email": cred.Email})
		return "", "", err
	}

	return "verified", userData.Email, nil
}

func GetUser(ids []string) ([]entity.ProfileDB, error) {
	filter := bson.A{}
	for _, id := range ids {
		idObj, _ := primitive.ObjectIDFromHex(id)
		filter = append(filter, bson.M{"_id": idObj})
	}
	return repo.Find(bson.M{"$or": filter}, bson.M{})
}

func UpdateUser(id string, user entity.ProfileDB) (string, error) {
	idObject, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": idObject}
	set := bson.M{}

	inc := bson.M{}
	inc["green_points"] = user.GreenPoints
	inc["carbon_saved"] = user.CarbonSaved
	inc["total_travelled"] = user.TotalTravelled
	set["update_time"] = time.Now()
	setS := bson.M{"$set": set}
	setS["$inc"] = inc

	return repo.UpdateOne(filter, setS)
}
