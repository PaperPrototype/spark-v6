package db

import (
	"html/template"
	"log"
	"main/markdown"
)

func (course *Course) GetCourseReleasesLogError() []Release {
	releases := []Release{}
	err := gormDB.Model(&Release{}).Where("course_id = ?", course.ID).Find(&releases).Error
	if err != nil {
		log.Println("db/methods ERROR getting course releases:", err)
	}

	return releases
}

func (course *Course) GetNewestCourseReleaseNumLogError() uint16 {
	release := Release{}
	err := gormDB.Model(&Release{}).Where("course_id = ?", course.ID).Order("num DESC").First(&release).Error
	if err != nil {
		log.Println("db/methods ERROR getting newest course release num:", err)
	}

	return release.Num
}

func (course *Course) GetNewestVersionIDLogError() uint64 {
	release := Release{}
	err := gormDB.Model(&Release{}).Where("course_id = ?", course.ID).Order("num DESC").First(&release).Error
	if err != nil {
		log.Println("db/methods ERROR getting newest course release id:", err)
	}

	version := Version{}
	err1 := gormDB.Model(&Version{}).Where("release_id = ?", release.ID).Order("num DESC").First(&version).Error
	if err1 != nil {
		log.Println("db/methods ERROR getting newest version id:", err1)
	}
	return version.ID
}

func (release *Release) GetVersionsLogError() []Version {
	versions := []Version{}
	err := gormDB.Model(&Version{}).Where("release_id = ?", release.ID).Order("num DESC").Find(&versions).Error
	if err != nil {
		log.Println("db/methods ERROR getting versions for release:", err)
	}
	return versions
}

func (release *Release) GetNewestVersionNumLogError() uint16 {
	version := Version{}
	err := gormDB.Model(&Version{}).Where("release_id = ?", release.ID).Order("num DESC").First(&version).Error
	if err != nil {
		log.Println("db/methods ERROR getting newest version num:", err)
	}
	return version.Num
}

func (version *Version) GetSectionsLogError() []Section {
	sections := []Section{}
	err := gormDB.Model(&Section{}).Where("version_id = ?", version.ID).Find(&sections).Error
	if err != nil {
		log.Println("db/methods ERROR getting sections for version:", err)
	}
	return sections
}

func (version *Version) GetFirstSectionPreload() (*Section, error) {
	section := Section{}
	err := gormDB.Model(&Section{}).Preload("Contents").Where("version_id = ?", version.ID).First(&section).Error
	return &section, err
}

func (section *Section) GetChildrenSectionsLogError() []Section {
	sections := []Section{}
	err := gormDB.Model(&Section{}).Where("parent_id = ?", section.ID).Find(&sections).Error
	if err != nil {
		log.Println("db/methods EROOR getting children sections for section:", err)
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

// based on the number of posts in a course version
func (release *Release) UserPostsCountLogError(userID uint64) int64 {
	postIDs := gormDB.Model(&PostToRelease{}).Select("post_id").Where("release_id = ?", release.ID)

	var count int64
	err := gormDB.Model(&Post{}).Where("user_id = ?", userID).Where("id IN (?)", postIDs).Count(&count).Error
	if err != nil {
		log.Println("db ERROR counting posts of a user for course release:", err)
	}

	return count
}

func (version *Version) SectionsCountLogError() int64 {
	var count int64
	err := gormDB.Model(&Section{}).Where("version_id = ?", version.ID).Count(&count).Error
	if err != nil {
		log.Println("db ERROR counting sections for version:", err)
	}

	return count
}

func (course *Course) GetNewestCourseReleaseLogError() *Release {
	release := Release{}
	err := gormDB.Model(&Release{}).Where("course_id = ?", course.ID).Order("num DESC").First(&release).Error
	if err != nil {
		log.Println("db ERROR getting newest course release:", err)
	}

	return &release
}

// get course that the user has posted a post to
// VERY EXPENSIVE QUERY
func (user *User) GetCourses() ([]Course, error) {
	releaseIDs := gormDB.Model(&Purchase{}).Select("release_id").Where("user_id = ?", user.ID)

	courseIDs := gormDB.Model(&Release{}).Select("course_id").Where("id IN (?)", releaseIDs)

	courses := []Course{}

	err := gormDB.Model(&Course{}).Where("id IN (?)", courseIDs).Find(&courses).Error
	return courses, err
}

func (release *Release) GetNewestVersionLogError() *Version {
	version := Version{}
	err := gormDB.Model(&Version{}).Where("release_id = ?", release.ID).Order("num DESC").First(&version).Error
	if err != nil {
		log.Println("db/methods ERROR getting newest version:", err)
	}
	return &version
}
