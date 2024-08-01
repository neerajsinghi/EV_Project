package entity

type ReferralDB struct {
	ReferralCode      string    `json:"referralCode" bson:"referral_code"`
	ReferrerBy        string    `json:"referrerBy" bson:"referrer_by"`
	ReferralOfId      string    `json:"referralOfId" bson:"referral_of_id"`
	ReferredByProfile ProfileDB `json:"referredByProfile" bson:"referred_by_profile"`
	ReferralOf        ProfileDB `json:"referralOf" bson:"referral_of"`
	ReferralStatus    string    `json:"referralStatus" bson:"referral_status"`
}
