package entity

type ReferralDB struct {
	ReferralCode      string    `json:"referralCode" bson:"referral_code"`
	ReferrerBy        string    `json:"referrerBy" bson:"referrer_by"`
	ReferredByProfile ProfileDB `json:"referredByProfile" bson:"referred_by_profile"`
	ReferralOf        ProfileDB `json:"referralOf" bson:"referral_of"`
}
