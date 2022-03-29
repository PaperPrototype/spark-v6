package db

// README!
/*
	This package contains all sensitive payments and stripe related db code
*/

import (
	"log"
	"main/payments"

	"github.com/stripe/stripe-go/v72/account"
)

func (stripeConnection *StripeConnection) DetailsSubmitted() (bool, error) {
	connectedAccount, err := account.GetByID(stripeConnection.StripeAccountID, nil)
	return connectedAccount.DetailsSubmitted, err
}

func (stripeConnection *StripeConnection) DetailsSubmittedLogError() bool {
	connectedAccount, err := account.GetByID(stripeConnection.StripeAccountID, nil)
	if err != nil {
		log.Println("db/methods_stripePayouts ERROR getting DetailsSubmitted param:", err)
	}
	return connectedAccount.DetailsSubmitted
}

// get the total amount we owe teacher from a course
func (purchase *Purchase) CalculatePayout() float32 {
	spark3DsCut := float32(purchase.AmountPaid) * payments.PercentageShare
	return float32(purchase.AmountPaid) - spark3DsCut
}

func (course *Course) GetCurrentTotalCoursePayoutAmountLogError() float64 {
	releaseIDs := gormDB.Model(&Release{}).Select("id").Where("course_id = ?", course.ID)

	purchases := []Purchase{}
	err := gormDB.Model(&Purchase{}).Where("release_id IN (?)", releaseIDs).Find(&purchases).Error
	if err != nil {
		log.Println("db ERROR getting GetCurrentTotalCoursePayoutAmount:", err)
	}

	var total float64 = 0
	for _, purchase := range purchases {
		total += float64(purchase.CalculatePayout())
	}

	return total
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

func GetCurrentTotalCoursePayoutAmount(courseID uint64) (float64, error) {
	releaseIDs := gormDB.Model(&Release{}).Select("id").Where("course_id = ?", courseID)

	purchases := []Purchase{}
	err := gormDB.Model(&Purchase{}).Where("release_id IN (?)", releaseIDs).Where("payed_out = ?", false).Find(&purchases).Error

	var total float64 = 0
	for _, purchase := range purchases {
		total += float64(purchase.CalculatePayout())
	}

	return total, err
}

func GetStripeConnection(userID interface{}) (*StripeConnection, error) {
	stripeConnection := StripeConnection{}
	err := gormDB.Model(&StripeConnection{}).Where("user_id = ?", userID).First(&stripeConnection).Error
	return &stripeConnection, err
}
