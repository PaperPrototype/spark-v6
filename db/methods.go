package db

import (
	"html/template"
	"log"
	"main/markdown"

	"gorm.io/gorm"
)

func (user *User) GetPublicAuthoredCourses() ([]Course, error) {
	courses := []Course{}
	err := GormDB.Model(&Course{}).Where("user_id = ?", user.ID).Where("public = ?", true).Preload("Release", orderByNewestCourseRelease).Preload("User").Find(&courses).Error
	return courses, err
}

func (user *User) GetPublicAndPrivateAuthoredCourses() ([]Course, error) {
	courses := []Course{}
	err := GormDB.Model(&Course{}).Preload("User").Where("user_id = ?", user.ID).Preload("Release", orderByNewestCourseRelease).Find(&courses).Error
	return courses, err
}

func (course *Course) GetPublicCourseReleasesLogError() []Release {
	releases := []Release{}
	err := GormDB.Model(&Release{}).Where("course_id = ?", course.ID).Where("public = ?", true).Order("num DESC").Find(&releases).Error
	if err != nil {
		log.Println("db/methods ERROR getting course releases:", err)
	}

	return releases
}

// don't filter out private releases
func (course *Course) GetAllCourseReleasesLogError() []Release {
	releases := []Release{}
	err := GormDB.Model(&Release{}).Where("course_id = ?", course.ID).Order("num DESC").Find(&releases).Error
	if err != nil {
		log.Println("db/methods ERROR getting course releases in GetAllCourseReleasesLogError:", err)
	}

	return releases
}

func (course *Course) GetAllPrerequisitesLogError() []Prerequisite {
	prerequisites := []Prerequisite{}
	err := GormDB.Model(&Prerequisite{}).Where("course_id = ?", course.ID).Preload("PrerequisiteCourse", func(db *gorm.DB) *gorm.DB {
		return db.Preload("User").Preload("Release", func(db *gorm.DB) *gorm.DB {
			return db.Order("num DESC")
		})
	}).Find(&prerequisites).Error
	if err != nil {
		log.Println("db/methods ERROR getting course prerequisites in GetAllPrerequisitesLogError:", err)
	}

	return prerequisites
}

func (course *Course) GetNewestPublicCourseReleaseNumLogError() uint16 {
	release := Release{}
	err := GormDB.Model(&Release{}).Where("course_id = ?", course.ID).Where("public = ?", true).Order("num DESC").First(&release).Error
	if err != nil {
		log.Println("db/methods ERROR getting newest course release num:", err)
	}

	return release.Num
}

func (course *Course) GetAllNewestCourseReleaseNumLogError() uint16 {
	release := Release{}
	err := GormDB.Model(&Release{}).Where("course_id = ?", course.ID).Order("num DESC").First(&release).Error
	if err != nil {
		log.Println("db/methods ERROR getting newest course release num:", err)
	}

	return release.Num
}

func (course *Course) GetVersion(versionID interface{}) (*Version, error) {
	version := Version{}
	err := GormDB.Model(&Version{}).Where("course_id = ?", course.ID).Where("id = ?", versionID).First(&version).Error
	return &version, err
}

func (release *Release) IsFree() bool {
	if release.Price <= 0 {
		return true
	}

	return false
}

func (release *Release) GetVersionsLogError() []Version {
	versions := []Version{}
	err := GormDB.Model(&Version{}).Where("release_id = ?", release.ID).Order("num DESC").Find(&versions).Error
	if err != nil {
		log.Println("db/methods ERROR getting versions for release:", err)
	}
	return versions
}

func (section *Section) GetChildrenSectionsLogError() []Section {
	sections := []Section{}
	err := GormDB.Model(&Section{}).Where("parent_id = ?", section.ID).Find(&sections).Error
	if err != nil {
		log.Println("db/methods ERROR getting sections for version in GetBaseSectionsLogError:", err)
	}
	return sections
}

func (content *Content) GetMarkdownHTMLLogError() template.HTML {
	buf, err := markdown.Convert([]byte(content.Markdown))
	if err != nil {
		log.Println("db/methods ERROR parsing markdown into html:", err)
	}

	return template.HTML(buf.Bytes())
}

// get the number of posts in a course version
func (release *Release) UserPostsCountLogError(userID uint64) int64 {
	postIDs := GormDB.Model(&PostToCourse{}).Select("post_id").Where("release_id = ?", release.ID)

	var count int64
	err := GormDB.Model(&Post{}).Where("user_id = ?", userID).Where("id IN (?)", postIDs).Count(&count).Error
	if err != nil {
		log.Println("db ERROR counting posts of a user for course release:", err)
	}

	return count
}

func (course *Course) GetNewestPublicCourseReleaseLogError() *Release {
	release := Release{}
	err := GormDB.Model(&Release{}).Where("course_id = ?", course.ID).Order("num DESC").Where("public = ?", true).First(&release).Error
	if err != nil {
		log.Println("db ERROR getting newest course release:", err)
	}

	return &release
}

// get courses that the user has purchased
func (user *User) GetPublicPurchasedCourses() ([]Course, error) {
	releaseIDs := GormDB.Model(&Purchase{}).Select("release_id").Where("user_id = ?", user.ID)

	courseIDs := GormDB.Model(&Release{}).Select("course_id").Where("id IN (?)", releaseIDs)

	courses := []Course{}

	err := GormDB.Model(&Course{}).Where("id IN (?)", courseIDs).Where("public = ?", true).Preload("Release", orderByNewestCourseRelease).Preload("User").Find(&courses).Error
	return courses, err
}

// get courses that the user has purchased
func (user *User) GetPublicAndPrivatePurchasedCourses() ([]Course, error) {
	releaseIDs := GormDB.Model(&Purchase{}).Select("release_id").Where("user_id = ?", user.ID)

	courseIDs := GormDB.Model(&Release{}).Select("course_id").Where("id IN (?)", releaseIDs)

	courses := []Course{}

	err := GormDB.Model(&Course{}).Preload("User").Where("id IN (?)", courseIDs).Preload("Release", orderByNewestCourseRelease).Find(&courses).Error
	return courses, err
}

func (release *Release) GetNewestVersionLogError() *Version {
	version := Version{}
	err := GormDB.Model(&Version{}).Where("release_id = ?", release.ID).Order("num DESC").First(&version).Error
	if err != nil {
		log.Println("db/methods ERROR getting newest version:", err)
	}
	return &version
}

func (release *Release) GetNewestVersion() (*Version, error) {
	version := Version{}
	err := GormDB.Model(&Version{}).Where("release_id = ?", release.ID).Order("num DESC").First(&version).Error
	return &version, err
}

func (release *Release) GetVersionsCountLogError() int64 {
	var count int64 = 0
	err := GormDB.Model(&Version{}).Where("release_id = ?", release.ID).Count(&count).Error
	if err != nil {
		log.Println("db/GetVersionsCountLogError ERROR counting versions:", err)
	}

	return count
}

func (purchase *Purchase) GetReleaseLogError() *Release {
	release := Release{}
	err := GormDB.Model(&Release{}).Where("id = ?", purchase.ReleaseID).First(&release).Error
	if err != nil {
		log.Println("db ERROR getting release:", err)
	}
	return &release
}

func (purchase *Purchase) GetUserLogError() *User {
	user := User{}
	err := GormDB.Model(&User{}).Where("id = ?", purchase.UserID).First(&user).Error
	if err != nil {
		log.Println("db ERROR getting user:", err)
	}
	return &user
}

func (purchase *Purchase) GetCourseLogError() *Course {
	course := Course{}
	err := GormDB.Model(&Course{}).Where("id = ?", purchase.CourseID).First(&course)
	if err != nil {
		log.Println("db/methods ERROR gettign course GetCourseLogError:", err)
	}

	return &course
}

func (user *User) SetVerified(verified bool) error {
	return GormDB.Model(&User{}).Where("id = ?", user.ID).Update("verified", verified).Error
}

func (release *Release) GetGithubReleaseLogError() *GithubRelease {
	githubRelease := GithubRelease{}
	err := GormDB.Model(&GithubRelease{}).Where("release_id = ?", release.ID).First(&githubRelease).Error
	if err != nil {
		log.Println("db/utils ERROR getting github release in GetGithubReleaseLogError:", err)
	}

	return &githubRelease
}

func GetGithubReleaseWithIDStr(releaseID string) (*GithubRelease, error) {
	githubRelease := GithubRelease{}
	err := GormDB.Model(&GithubRelease{}).Where("release_id = ?", releaseID).First(&githubRelease).Error
	if err != nil {
		log.Println("db/utils ERROR getting github release in GetGithubReleaseLogError:", err)
	}

	return &githubRelease, err
}

func (user *User) NewNotifLogError(message, url string) {
	notif := Notif{
		UserID:  user.ID,
		Message: message,
		URL:     url,
	}
	err := GormDB.Create(&notif).Error
	if err != nil {
		log.Println("db/methods ERROR creating notif in NewNotifLogError:", err)
	}
}

// gives back markdown casted to the template.HTML type
func (post *Post) MarkdownTemplateHTML() template.HTML {
	return template.HTML(post.Markdown)
}

func (user *User) GetReleaseOwnership(releaseID uint64) (*Ownership, error) {
	ownership := Ownership{}
	err := GormDB.Model(&Ownership{}).Where("user_id = ?", user.ID).Where("release_id = ?", releaseID).First(&ownership).Error
	return &ownership, err
}

func (user *User) HasGithubConnection() bool {
	var count int64 = 0
	err := GormDB.Model(&GithubConnection{}).Where("user_id = ?", user.ID).Count(&count).Error

	// if err then not valid
	if err != nil {
		log.Println("db/github ERROR getting github connection:", err)
		return false
	}

	// if nothing exists
	if count == 0 {
		return false
	}

	return true
}
