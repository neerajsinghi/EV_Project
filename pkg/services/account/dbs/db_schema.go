package db

import (
	"bikeRental/pkg/entity"
)

type AccountService interface {
	SignUp(cred Credentials) (string, error)
	LoginUsingPhone(phone, deviceID string) (string, error)
	LoginUsingPassword(cred Credentials) (entity.ProfileDB, string, error)
	VerifyOTP(phone, otp, token, deviceID string) (*entity.ProfileDB, string, error)
	VerifyEmail(cred Credentials) (string, error)
	SendVerificationEmail(email, pemail, uid string) (string, error)
	SendEmailOTP(email string) (string, error)
	ChangePassword(cred Credentials) (string, error)
	VerifyResetLink(cred Credentials) (string, string, error)
	Logout(id string) (string, error)
}

// Credentials represents the model for an credentials
type Credentials struct {
	entity.ProfileDB
}

type OTP struct {
	Email string `bson:"email" json:"email"`
	OTP   string `bson:"otp_code" json:"otp_code"`
}

type EmailVerification struct {
	Email string `bson:"email" json:"email"`
}
