package db

// README!
/*
	This package contains all sensitive payments and stripe related db code
*/

import (
	"log"

	"github.com/stripe/stripe-go/v72/account"
)

func (stripeConnection *StripeConnection) DetailsSubmitted() (bool, error) {
	connectedAccount, err := account.GetByID(stripeConnection.StripeAccountID, nil)
	return connectedAccount.DetailsSubmitted, err
}

func (stripeConnection *StripeConnection) DetailsSubmittedLogError() bool {
	connectedAccount, err := account.GetByID(stripeConnection.StripeAccountID, nil)
	if err != nil {
		log.Println("db/methods_stripe_payments ERROR getting DetailsSubmitted param:", err)
	}
	return connectedAccount.DetailsSubmitted
}

func (stripeConnection *StripeConnection) ChargesEnabled() (bool, error) {
	connectedAccount, err := account.GetByID(stripeConnection.StripeAccountID, nil)
	return connectedAccount.ChargesEnabled, err
}

func (stripeConnection *StripeConnection) ChargesEnabledLogError() bool {
	connectedAccount, err := account.GetByID(stripeConnection.StripeAccountID, nil)
	if err != nil {
		log.Println("db/methods_stripe_payments ERROR getting ChargesEnabled param:", err)
	}
	return connectedAccount.ChargesEnabled
}

func (course *Course) GetPurchasesLogError() []Purchase {
	purchases := []Purchase{}
	err := gormDB.Model(&Purchase{}).Where("course_id = ?", course.ID).Find(&purchases).Error
	if err != nil {
		log.Println("db/methods ERROR getting purchases from GetPurchasesLogError:", err)
	}

	return purchases
}

func (user *User) HasStripeConnection() bool {
	var count int64 = 0
	err := gormDB.Model(&StripeConnection{}).Where("user_id = ?", user.ID).Count(&count).Error

	// if err then not valid
	if err != nil {
		log.Println("db/stripe ERROR getting stripe connection:", err)
		return false
	}

	// if nothing exists
	if count == 0 {
		return false
	}

	return true
}

func (user *User) HasPurchasedRelease(releaseID uint64) bool {
	var count int64 = 0
	err := gormDB.Model(&Purchase{}).Where("user_id = ?", user.ID).Where("release_id = ?", releaseID).Count(&count).Error

	// if err then not valid
	if err != nil {
		return false
	}

	// if nothing exists
	if count == 0 {
		return false
	}

	return true
}

func GetStripeConnection(userID interface{}) (*StripeConnection, error) {
	stripeConnection := StripeConnection{}
	err := gormDB.Model(&StripeConnection{}).Where("user_id = ?", userID).First(&stripeConnection).Error
	return &stripeConnection, err
}
