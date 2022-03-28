package db

func UpdateCourse(courseID, title, name, subtitle string) error {
	return gormDB.Model(&Course{}).Where("id = ?", courseID).Update("name", name).Update("title", title).Update("subtitle", subtitle).Error
}

func UpdateRelease(releaseID, markdown, price string, public bool) error {
	return gormDB.Model(&Release{}).Where("id = ?", releaseID).Update("markdown", markdown).Update("price", price).Update("public", public).Error
}

func UpdatePost(postID, markdown string) error {
	return gormDB.Model(&Post{}).Where("id = ?", postID).Update("markdown", markdown).Error
}

func UpdateSectionContentAndIncreasePatch(sectionID string, contentID string, contentMarkdown string, versionID string) error {
	err := gormDB.Model(&Content{}).Where("id = ?", contentID).Update("markdown", contentMarkdown).Error
	if err != nil {
		return err
	}

	version := Version{}
	err1 := gormDB.Model(&Version{}).Where("id = ?", versionID).First(&version).Error
	if err1 != nil {
		return err1
	}

	err2 := gormDB.Model(&Version{}).Where("id = ?", versionID).Update("patch", version.Patch+1).Error
	return err2
}
