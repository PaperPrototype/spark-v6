package db

import (
	"log"
	"main/markdown"
)

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

func GetCourse(name string) (*Course, error) {
	course := Course{}
	err := gormDB.Model(&Course{}).Where("name = ?", name).First(&course).Error
	return &course, err
}

func GetAllCourses() ([]Course, error) {
	courses := []Course{}
	err := gormDB.Model(&Course{}).Find(&courses).Error
	return courses, err
}

func GetCourseWithIDStr(courseID string) (*Course, error) {
	course := Course{}
	err := gormDB.Model(&Course{}).Where("id = ?", courseID).First(&course).Error
	return &course, err
}

func GetCoursewithID(courseID uint64) (*Course, error) {
	course := Course{}
	err := gormDB.Model(&Course{}).Where("id = ?", courseID).First(&course).Error
	return &course, err
}

func GetReleaseWithIDStr(releaseID string) (*Release, error) {
	release := Release{}
	err := gormDB.Where("id = ?", releaseID).First(&release).Error
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

func GetReleasePosts(releaseID uint64) ([]Post, error) {
	postIDs := gormDB.Model(&PostToRelease{}).Select("post_id").Where("release_id = ?", releaseID)

	posts := []Post{}
	err := gormDB.Model(&Post{}).Where("id IN (?)", postIDs).Find(&posts).Error

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

func GetRelease(releaseID uint64) (*Release, error) {
	release := Release{}
	err := gormDB.Model(&Release{}).Where("id = ?", releaseID).First(&release).Error
	return &release, err
}
func GetPost(postID string) (*Post, error) {
	post := Post{}
	err := gormDB.Model(&Post{}).Where("id = ?", postID).First(&post).Error
	return &post, err
}

func GetMedia(versionID string, mediaName string) (*Media, error) {
	media := Media{}
	err := gormDB.Model(&Media{}).Where("version_id = ?", versionID).Where("name = ?", mediaName).First(&media).Error
	return &media, err
}

func GetNewestCourseRelease(courseID uint64) (*Release, error) {
	release := Release{}
	err := gormDB.Model(&Release{}).Where("course_id = ?", courseID).Order("num DESC").First(&release).Error
	return &release, err
}

func GetBuyRelease(buyReleaseID string) (*BuyRelease, error) {
	err := DeleteExpiredBuyReleases()
	if err != nil {
		log.Println("db ERROR deleting expired releases:", err)
		return nil, err
	}

	buyRelease := BuyRelease{}
	err1 := gormDB.Model(&BuyRelease{}).Where("id = ?", buyReleaseID).First(&buyRelease).Error
	return &buyRelease, err1
}

func GetNewestReleaseVersion(releaseID uint64) (*Version, error) {
	version := Version{}
	err := gormDB.Model(&Version{}).Where("release_id = ?", releaseID).Order("num DESC").First(&version).Error
	return &version, err
}
