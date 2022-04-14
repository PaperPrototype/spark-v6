package db

import (
	"main/markdown"
	"time"
)

func GetCourse(courseID interface{}) (*Course, error) {
	course := Course{}
	err := gormDB.Model(&Course{}).Where("id = ?", courseID).First(&course).Error
	return &course, err
}

func GetUserFromSession(token string) (*User, error) {
	session, err := GetSession(token)
	if err != nil {
		return &User{}, err
	}

	user, err1 := GetUser(session.UserID)
	return user, err1
}

func GetSession(token string) (*Session, error) {
	session := Session{}
	err := gormDB.Model(&Session{}).Where("token_uuid = ?", token).First(&session).Error
	return &session, err
}

func GetUser(userID uint64) (*User, error) {
	user := User{}
	err := gormDB.Model(&User{}).Where("id = ?", userID).First(&user).Error
	return &user, err
}

func GetUserCoursePreloadUser(username string, courseName string) (*Course, error) {
	user := User{}
	course := Course{}

	err := gormDB.Model(&User{}).Where("username = ?", username).First(&user).Error
	if err != nil {
		return &course, err
	}

	err1 := gormDB.Model(&Course{}).Where("user_id = ?", user.ID).Where("name = ?", courseName).First(&course).Error

	// fill in user
	course.User = user

	// return
	return &course, err1
}

func GetUserCourseWithIDPreloadUser(username string, courseID interface{}) (*Course, error) {
	user := User{}
	course := Course{}

	err := gormDB.Model(&User{}).Where("username = ?", username).First(&user).Error
	if err != nil {
		return &course, err
	}

	err1 := gormDB.Model(&Course{}).Where("user_id = ?", user.ID).Where("id = ?", courseID).First(&course).Error

	// fill in user
	course.User = user

	// return
	return &course, err1
}

func GetWithIDUserCoursePreloadUser(userID interface{}, courseID interface{}) (*Course, error) {
	user := User{}
	course := Course{}

	err := gormDB.Model(&User{}).Where("id = ?", userID).First(&user).Error
	if err != nil {
		return &course, err
	}

	err1 := gormDB.Model(&Course{}).Where("user_id = ?", user.ID).Where("id = ?", courseID).First(&course).Error

	// fill in user
	course.User = user

	// return
	return &course, err1
}

func GetCourseWithIDPreloadUser(courseID interface{}) (*Course, error) {
	course := Course{}
	err := gormDB.Model(&Course{}).Preload("User").Where("id = ?", courseID).First(&course).Error
	return &course, err
}

func GetPublicReleaseWithID(releaseID interface{}) (*Release, error) {
	release := Release{}
	err := gormDB.Where("id = ?", releaseID).Where("public = ?", true).First(&release).Error
	return &release, err
}

func GetVersion(versionID string) (*Version, error) {
	version := Version{}
	err := gormDB.Model(&Version{}).Where("id = ?", versionID).First(&version).Error
	return &version, err
}

func GetCourseReleaseNumString(courseID uint64, releaseNum string) (*Release, error) {
	release := Release{}
	err := gormDB.Model(&Release{}).Where("course_id = ?", courseID).Where("num = ?", releaseNum).First(&release).Error
	return &release, err
}

func GetSectionPreloadConvertMarkdown(sectionID string) (*Section, error) {
	section := Section{}
	err := gormDB.Model(&Section{}).Preload("Contents").Where("id = ?", sectionID).First(&section).Error
	for i := range section.Contents {
		buf, err := markdown.Convert([]byte(section.Contents[i].Markdown))
		if err != nil {
			return &section, err
		}
		section.Contents[i].Markdown = buf.String()
	}
	return &section, err
}

func GetSectionPreload(sectionID string) (*Section, error) {
	section := Section{}
	err := gormDB.Model(&Section{}).Preload("Contents").Where("id = ?", sectionID).First(&section).Error
	return &section, err
}

func GetReleasePosts(releaseID uint64, courseID uint64) ([]Post, error) {
	postIDs := gormDB.Model(&PostToCourse{}).Select("post_id").Where("release_id = ?", releaseID)

	posts := []Post{}
	err := gormDB.Model(&Post{}).Where("id IN (?)", postIDs).Order("created_at DESC").Find(&posts).Error

	for i := range posts {
		buf, err := markdown.Convert([]byte(posts[i].Markdown))
		if err != nil {
			return posts, err
		}
		posts[i].Markdown = buf.String()
	}

	return posts, err
}

func GetPostPreloadUser(postID string) (*Post, error) {
	post := Post{}
	err := gormDB.Model(&Post{}).Where("id = ?", postID).Preload("User").First(&post).Error
	return &post, err
}

func GetUserWithUsername(username string) (*User, error) {
	user := User{}
	err := gormDB.Model(&User{}).Where("username = ?", username).First(&user).Error
	return &user, err
}

func GetPublicRelease(releaseID interface{}) (*Release, error) {
	release := Release{}
	err := gormDB.Model(&Release{}).Where("id = ?", releaseID).Where("public = ?", true).First(&release).Error
	return &release, err
}

func GetAllRelease(releaseID interface{}) (*Release, error) {
	release := Release{}
	err := gormDB.Model(&Release{}).Where("id = ?", releaseID).First(&release).Error
	return &release, err
}

func GetPost(postID string) (*Post, error) {
	post := Post{}
	err := gormDB.Model(&Post{}).Where("id = ?", postID).First(&post).Error
	return &post, err
}

func GetNewestPublicCourseRelease(courseID uint64) (*Release, error) {
	release := Release{}
	err := gormDB.Model(&Release{}).Where("course_id = ?", courseID).Where("public = ?", true).Order("num DESC").First(&release).Error
	return &release, err
}

func GetAllNewestCourseRelease(courseID uint64) (*Release, error) {
	release := Release{}
	err := gormDB.Model(&Release{}).Where("course_id = ?", courseID).Order("num DESC").First(&release).Error
	return &release, err
}

func GetBuyRelease(stripeSessionID string) (*AttemptBuyRelease, error) {
	buyRelease := AttemptBuyRelease{}
	err1 := gormDB.Model(&AttemptBuyRelease{}).Where("stripe_session_id = ?", stripeSessionID).Where("expires_at > ?", time.Now()).First(&buyRelease).Error
	return &buyRelease, err1
}

func GetNewestReleaseVersion(releaseID uint64) (*Version, error) {
	version := Version{}
	err := gormDB.Model(&Version{}).Where("release_id = ?", releaseID).Order("num DESC").First(&version).Error
	return &version, err
}

func GetPurchases(courseID uint64) ([]Purchase, error) {
	releaseIDs := gormDB.Model(&Release{}).Select("id").Where("course_id = ?", courseID)

	purchases := []Purchase{}
	err := gormDB.Model(&Purchase{}).Where("release_id IN (?)", releaseIDs).Find(&purchases).Error
	return purchases, err
}

func GetAuthorPublicCourses(userID uint64) ([]Course, error) {
	courses := []Course{}
	err := gormDB.Model(&Course{}).Where("user_id = ?", userID).Where("public = ?", true).Find(&courses).Error
	return courses, err
}

// careful! this is public AND private courses
func GetAuthorPublicAndPrivateCourses(userID uint64) ([]Course, error) {
	courses := []Course{}
	err := gormDB.Model(&Course{}).Where("user_id = ?", userID).Find(&courses).Error
	return courses, err
}

func GetAllPublicCoursesPreloadUser() ([]Course, error) {
	courses := []Course{}
	err := gormDB.Model(&Course{}).Preload("User").Where("public = ?", true).Find(&courses).Error
	return courses, err
}

func GetVerify(verifyUUID string) (*Verify, error) {
	err1 := DeleteExpiredVerify()
	if err1 != nil {
		return nil, err1
	}

	verify := Verify{}
	err := gormDB.Model(&Verify{}).Where("verify_uuid = ?", verifyUUID).First(&verify).Error
	return &verify, err
}
