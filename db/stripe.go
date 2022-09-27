package db

// README!
/*
	This package contains all sensitive payments and stripe related db code
*/

import (
	"log"

	"github.com/stripe/stripe-go/v73/account"
)

// if the account has finished onboarding
func (stripeConnection *StripeConnection) DetailsSubmitted() (bool, error) {
	connectedAccount, err := account.GetByID(stripeConnection.StripeAccountID, nil)
	return connectedAccount.DetailsSubmitted, err
}

// if the account has finished onboarding
func (stripeConnection *StripeConnection) DetailsSubmittedLogError() bool {
	connectedAccount, err := account.GetByID(stripeConnection.StripeAccountID, nil)
	if err != nil {
		log.Println("db/stripe.go ERROR getting DetailsSubmitted param:", err)
	}
	return connectedAccount.DetailsSubmitted
}

// if the account can create charges
func (stripeConnection *StripeConnection) PayoutsEnabled() (bool, error) {
	connectedAccount, err := account.GetByID(stripeConnection.StripeAccountID, nil)
	return connectedAccount.PayoutsEnabled, err
}

// if the account can create charges
func (stripeConnection *StripeConnection) PayoutsEnabledLogError() bool {
	connectedAccount, err := account.GetByID(stripeConnection.StripeAccountID, nil)
	if err != nil {
		log.Println("db/stripe.go ERROR getting PayoutsEnabled param:", err)
		return false // probably not and they should update their stripe info!
	}
	return connectedAccount.PayoutsEnabled
}

// TODO
// using paging to load up to 20 purchases at a time
func (course *Course) GetPurchasesLogError() []Purchase {
	purchases := []Purchase{}
	err := GormDB.Model(&Purchase{}).Where("course_id = ?", course.ID).Preload("User").Order("created_at DESC").Limit(20).Find(&purchases).Error
	if err != nil {
		log.Println("db/stripe.go ERROR getting purchases from GetPurchasesLogError:", err)
	}

	return purchases
}

func (user *User) HasStripeConnection() bool {
	var count int64 = 0
	err := GormDB.Model(&StripeConnection{}).Where("user_id = ?", user.ID).Count(&count).Error

	// if err then not valid
	if err != nil {
		log.Println("db/stripe.go ERROR getting stripe connection:", err)
		return false
	}

	// if nothing exists
	if count == 0 {
		return false
	}

	return true
}

func (user *User) OwnsRelease(releaseID uint64) bool {
	var count int64 = 0
	err := GormDB.Model(&Ownership{}).Where("user_id = ?", user.ID).Where("release_id = ?", releaseID).Count(&count).Error

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
	err := GormDB.Model(&StripeConnection{}).Where("user_id = ?", userID).First(&stripeConnection).Error
	return &stripeConnection, err
}
