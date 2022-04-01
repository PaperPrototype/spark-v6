package db

import (
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func UserCourseNameAvailable(username string, name string) (bool, error) {
	user := User{}
	err1 := gormDB.Model(&User{}).Where("username = ?", username).First(&user).Error
	if err1 != nil {
		return false, err1
	}

	var count int64 = 0
	err := gormDB.Model(&Course{}).Where("user_id = ?", user.ID).Where("name = ?", name).Count(&count).Error
	// if err then taken
	if err != nil {
		return false, err
	}

	// if there is another course with that name taken
	if count != 0 {
		return false, nil
	}

	return true, nil
}

func UserCourseNameAvailableNotSelf(username string, name string, courseID interface{}) (bool, error) {
	user := User{}
	err1 := gormDB.Model(&User{}).Where("username = ?", username).First(&user).Error
	if err1 != nil {
		return false, err1
	}

	var count int64 = 0
	err := gormDB.Model(&Course{}).Where("user_id = ?", user.ID).Where("name = ?", name).Where("id != ?", courseID).Count(&count).Error

	log.Println("checking if course name available not self.")

	// if err then taken
	if err != nil {
		return false, err
	}

	// if there is another course with that name taken
	if count != 0 {
		return false, nil
	}

	return true, nil
}

func UsernameAvailable(username string) (bool, error) {
	var count int64 = 0
	err := gormDB.Model(&User{}).Where("username = ?", username).Count(&count).Error
	// if err then taken
	if err != nil {
		return false, err
	}

	// if there is another user with that name taken
	if count != 0 {
		return false, nil
	}

	return true, nil
}

func EmailAvailable(email string) (bool, error) {
	var count int64 = 0
	err := gormDB.Model(&User{}).Where("email = ?", email).Count(&count).Error
	// if err then taken
	if err != nil {
		return false, err
	}

	// if there is another user with that email
	if count != 0 {
		return false, nil
	}

	return true, nil
}

func SessionExists(tokenUUID string) bool {
	var count int64 = 0
	err := gormDB.Model(&Session{}).Where("token_uuid = ?", tokenUUID).Count(&count).Error

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

// returns true if successfully entered password
func TryUserPassword(username string, password string) (*User, bool) {
	user := User{}
	err := gormDB.Model(&User{}).Where("username = ?", username).First(&user).Error

	// err == failed
	if err != nil {
		return &user, false
	}

	// returns error == failed
	if bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(password)) != nil {
		return &user, false
	}

	return &user, true
}

// check if the user is the owner or if the user has purchased the course
func UserCanAccessCourseRelease(userID uint64, version *Version) bool {
	course := Course{}
	err1 := gormDB.Model(&Course{}).Where("id = ?", version.CourseID).First(&course).Error
	if err1 != nil {
		return false
	}

	// if they are the owner of the course
	if userID == course.UserID {
		return true
	}

	var count int64 = 0
	err := gormDB.Model(&Purchase{}).Where("user_id = ?", userID).Where("release_id = ?", version.ReleaseID).Count(&count).Error

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

func DeleteExpiredBuyReleases() error {
	return gormDB.Where("expires_at < ?", time.Now()).Delete(&AttemptBuyRelease{}).Error
}

func (release *Release) HasVersions() bool {
	var count int64 = 0
	err := gormDB.Model(&Version{}).Where("release_id = ?", release.ID).Count(&count).Error

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