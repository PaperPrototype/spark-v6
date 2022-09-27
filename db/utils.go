package db

import (
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func UserCourseNameAvailable(username string, name string) (bool, error) {
	user := User{}
	err1 := GormDB.Model(&User{}).Where("username = ?", username).First(&user).Error
	if err1 != nil {
		return false, err1
	}

	var count int64 = 0
	err := GormDB.Model(&Course{}).Where("user_id = ?", user.ID).Where("name = ?", name).Count(&count).Error
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
	err1 := GormDB.Model(&User{}).Where("username = ?", username).First(&user).Error
	if err1 != nil {
		return false, err1
	}

	var count int64 = 0
	err := GormDB.Model(&Course{}).Where("user_id = ?", user.ID).Where("name = ?", name).Where("id != ?", courseID).Count(&count).Error

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
	err := GormDB.Model(&User{}).Where("username = ?", username).Count(&count).Error
	// if err then taken
	if err != nil {
		log.Println("db/utils ERROR checking if username is available in UsernameAvailableLogError:", err)
		return false, err
	}

	// if there is another user with that name taken
	if count != 0 {
		return false, nil
	}

	return true, nil
}

func UsernameAvailableIgnoreError(username string) bool {
	var count int64 = 0
	err := GormDB.Model(&User{}).Where("username = ?", username).Count(&count).Error
	// if err then taken
	if err != nil {
		log.Println("db/utils ERROR checking if username is available in UsernameAvailableLogError:", err)
		return false
	}

	// if there is another user with that name taken
	if count != 0 {
		return false
	}

	return true
}

func EmailAvailable(email string) (bool, error) {
	var count int64 = 0
	err := GormDB.Model(&User{}).Where("email = ?", email).Count(&count).Error
	// if err then taken
	if err != nil {
		log.Println("db/utils ERROR checking if email is available in EmailAvailable:", err)
		return false, err
	}

	// if there is another user with that email
	if count != 0 {
		return false, nil
	}

	return true, nil
}

func EmailAvailableIgnoreError(email string) bool {
	var count int64 = 0
	err := GormDB.Model(&User{}).Where("email = ?", email).Count(&count).Error
	// if err then taken
	if err != nil {
		log.Println("db/utils ERROR checking if email is available in EmailAvailable:", err)
		return false
	}

	// if there is another user with that email
	if count != 0 {
		return false
	}

	return true
}

func SessionExists(tokenUUID string) bool {
	var count int64 = 0
	err := GormDB.Model(&Session{}).Where("token_uuid = ?", tokenUUID).Count(&count).Error

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
	err := GormDB.Model(&User{}).Where("username = ?", username).First(&user).Error

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

// returns true if successfully entered password
func TryEmailPassword(email string, password string) (*User, bool) {
	user := User{}
	err := GormDB.Model(&User{}).Where("email = ?", email).First(&user).Error

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
	err1 := GormDB.Model(&Course{}).Where("id = ?", version.CourseID).First(&course).Error
	if err1 != nil {
		return false
	}

	// FIRST check if it is the author...
	// if they are the owner of the course
	if userID == course.UserID {
		return true
	}

	// THEN if not author... check if they own the course
	var count int64 = 0
	err := GormDB.Model(&Ownership{}).Where("user_id = ?", userID).Where("release_id = ?", version.ReleaseID).Count(&count).Error

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

func (release *Release) HasVersions() bool {
	var count int64 = 0
	err := GormDB.Model(&Version{}).Where("release_id = ?", release.ID).Count(&count).Error

	// if err then not valid
	if err != nil {
		log.Println("db/utils ERROR getting versions in HasVersions:", err)
		return false
	}

	// if nothing exists
	if count == 0 {
		return false
	}

	return true
}

func (release *Release) HasGithubRelease() bool {
	var count int64 = 0
	err := GormDB.Model(&GithubRelease{}).Where("release_id = ?", release.ID).Count(&count).Error

	// if err then not valid
	if err != nil {
		log.Println("db/utils ERROR getting github release in HasGithubRelease:", err)
		return false
	}

	// if nothing exists
	if count == 0 {
		return false
	}

	return true
}

func (version *Version) HasGithubVersion() bool {
	var count int64 = 0
	err := GormDB.Model(&GithubVersion{}).Where("version_id = ?", version.ID).Count(&count).Error

	// if err then not valid
	if err != nil {
		log.Println("db/utils ERROR getting github release in HasGithubRelease:", err)
		return false
	}

	// if nothing exists
	if count == 0 {
		return false
	}

	return true
}

func CountPublicCourseReleasesLogError(courseID string) int64 {
	var count int64 = 0
	err := GormDB.Model(&Release{}).Where("course_id = ?", courseID).Where("public = ?", true).Count(&count).Error
	if err != nil {
		log.Println("db/utils ERROR getting github release in HasGithubRelease:", err)
	}

	return count
}

func (version *Version) GetAuthorUser() (*User, error) {
	userIDs := GormDB.Model(&Course{}).Select("user_id").Where("id = ?", version.CourseID)
	user := User{}
	err := GormDB.Model(&User{}).Where("id IN (?)", userIDs).First(&user).Error
	return &user, err
}

func (release *Release) GetAuthorUser() (*User, error) {
	userIDs := GormDB.Model(&Course{}).Select("user_id").Where("id = ?", release.CourseID)
	user := User{}
	err := GormDB.Model(&User{}).Where("id IN (?)", userIDs).First(&user).Error
	return &user, err
}

func CountCourseReviewsLogError(courseID uint64) int64 {
	var count int64
	err := GormDB.Model(&PostToCourseReview{}).Where("course_id = ?", courseID).Count(&count).Error
	if err != nil {
		log.Println("db/utils ERROR counting PostToCourseReviews in CountCourseReviewsLogError:", err)
	}
	return count
}

func CountUserReviewsLogError(userID uint64, courseID uint64) int64 {
	var count int64
	err := GormDB.Model(&PostToCourseReview{}).Where("course_id = ?", courseID).Where("user_id = ?", userID).Count(&count).Error
	if err != nil {
		log.Println("db/utils ERROR counting PostToCourseReviews in CountUserReviewsLogError:", err)
	}
	return count
}

func orderByNewestCourseRelease(db *gorm.DB) *gorm.DB {
	return db.Order("created_at DESC")
}
