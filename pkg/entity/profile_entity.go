package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Profile represents the model for an profile
type ProfileDB struct {
	ID                  primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Email               string             `bson:"email" json:"email,omitempty"`
	Status              string             `bson:"status" json:"status,omitempty"`
	StatusBool          *bool              `bson:"status_bool" json:"statusBool,omitempty"`
	JoiningDate         string             `bson:"joining_date" json:"joiningDate,omitempty"`
	Name                string             `bson:"name" json:"name,omitempty"`
	DOB                 string             `bson:"dob" json:"dob,omitempty"`
	Designation         *string            `bson:"designation" json:"designation,omitempty"`
	Gender              *string            `bson:"gender" json:"gender,omitempty"`
	PhoneNo             string             `bson:"phone_no" json:"phoneNumber,omitempty"`
	PhoneOTP            *string            `bson:"phone_otp" json:"phoneOtp,omitempty"`
	Role                *string            `bson:"roles" json:"role,omitempty"`
	PhoneNoVerified     bool               `bson:"phone_no_verified" json:"phoneNumberVerified,omitempty"`
	Address             *AddressDB         `bson:"address" json:"address,omitempty"`
	About               *string            `bson:"about" json:"about,omitempty"`
	UrlToProfileImage   *string            `bson:"url_to_profile_image" json:"url_to_profile_image,omitempty"`
	Password            *string            `bson:"password" json:"password,omitempty"`
	CreatedTime         primitive.DateTime `bson:"created_time" json:"createdTime,omitempty"`
	EmailLoginOTP       *string            `bson:"email_login_otp" json:"emailLoginOtp,omitempty"`
	OTP                 *string            `bson:"otp_code" json:"otp_code,omitempty"`
	UpdateTime          time.Time          `bson:"update_time" json:"updateTime,omitempty"`
	EmailSentTime       *time.Time         `bson:"email_sent_time" json:"emailSentTime,omitempty"`
	VerificationCode    *string            `bson:"verification_code" json:"verificationCode,omitempty"`
	PasswordResetCode   *string            `bson:"password_reset_code" json:"passwordResetCode,omitempty"`
	CountryCode         *string            `bson:"country_code" json:"countryCode,omitempty"`
	PasswordResetTime   time.Time          `bson:"password_reset_time" json:"passwordResetTime,omitempty"`
	LastLoginDeviceID   *string            `bson:"last_login_device_id" json:"lastLoginDeviceID,omitempty"`
	LastLoginDeviceName *string            `bson:"last_login_device_name" json:"lastLoginDeviceName,omitempty"`
	LastLoginLocation   *string            `bson:"last_login_location" json:"lastLoginLocation,omitempty"`
	Online              *bool              `bson:"online" json:"online,omitempty"`
	DLVerified          *bool              `bson:"dl_verified" json:"dlVerified,omitempty"`
	DLFrontImage        string             `bson:"dl_front_image" json:"dlFrontImage,omitempty"`
	DLBackImage         string             `bson:"dl_back_image" json:"dlBackImage,omitempty"`
	IDFrontImage        string             `bson:"id_front_image" json:"idFrontImage,omitempty"`
	IDBackImage         string             `bson:"id_back_image" json:"idBackImage,omitempty"`
	IDVerified          *bool              `bson:"id_verified" json:"idVerified,omitempty"`
	PlanID              *string            `bson:"plan_id" json:"planId,omitempty"`
	Plan                *PlanDB            `bson:"plan" json:"plan,omitempty"`
	PlanStartTime       int64              `bson:"plan_start_time" json:"planStartTime,omitempty"`
	PlanActive          *bool              `bson:"plan_active" json:"planActive,omitempty"`
	UserBlocked         *bool              `bson:"user_blocked" json:"userBlocked,omitempty"`
	PlanEndTime         int64              `bson:"plan_end_time" json:"planEndTime,omitempty"`
	PlanRemainingTime   int64              `bson:"plan_remaining_time" json:"planRemainingTime,omitempty"`
	ReferralCode        *string            `bson:"referral_code" json:"referralCode,omitempty"`
	ReferralCodeUsed    *string            `bson:"referral_code_used" json:"referralCodeUsed,omitempty"`
	Access              interface{}        `bson:"access" json:"access,omitempty"`
	TermsChecked        *bool              `bson:"terms_and_condition" json:"termsAndConditions,omitempty"`
	AllowPromotions     *bool              `bson:"allow_promotions" json:"allowPromotions,omitempty"`
	FirebaseToken       *string            `bson:"firebase_token" json:"firebaseToken,omitempty"`
	TotalBalance        float64            `bson:"total_balance" json:"totalBalance,omitempty"`
	TotalRides          int64              `bson:"total_rides" json:"totalRides,omitempty"`
	GreenPoints         int64              `bson:"green_points" json:"greenPoints"`
	CarbonSaved         float64            `bson:"carbon_saved" json:"carbonSaved"`
	TotalTravelled      float64            `bson:"total_travelled" json:"totalTravelled,omitempty"`
	BlockedBy           string             `bson:"blocked_by" json:"blockedBy,omitempty"`
	BlockReason         string             `bson:"block_reason" json:"blockReason,omitempty"`
	ServiceType         string             `bson:"service_type" json:"serviceType,omitempty"`
	StaffId             string             `bson:"staff_id" json:"staffId,omitempty"`
	Stations            int64              `bson:"stations" json:"stations,omitempty"`
	StaffStatus         string             `bson:"staff_status" json:"staffStatus,omitempty"`
	StaffShiftStartTime time.Time          `bson:"staff_shift_start_time" json:"staffShiftStartTime,omitempty"`
	StaffShiftEndTime   time.Time          `bson:"staff_shift_end_time" json:"staffShiftEndTime,omitempty"`
	StaffVerificationId string             `bson:"staff_verification_id" json:"staffVerificationId,omitempty"`
}
type AddressDB struct {
	Address string `bson:"address" json:"address,omitempty"`
	Country string `bson:"country" json:"country,omitempty"`
	Pin     string `bson:"pin" json:"pin,omitempty"`
	City    string `bson:"city" json:"city,omitempty"`
	State   string `bson:"state" json:"state,omitempty"`
}

type ProfileOut struct {
	ProfileDB `json:",inline"`
	Plan      *PlanDB     `json:"plan"`
	Booking   []BookingDB `json:"booking"`
	Wallet    []WalletS   `json:"wallet"`
	Station   []StationDB `json:"station"`
}

type UserAttendance struct {
	ProfileID   string    `bson:"profile_id" json:"profileId,omitempty"`
	Status      string    `bson:"status" json:"status,omitempty"`
	UpdatedTime time.Time `bson:"updated_time" json:"updatedTime,omitempty"`
}
