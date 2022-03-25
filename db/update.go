package db

func UpdateCourse(courseID, title, name, desc string) error {
	return gormDB.Model(&Course{}).Where("id = ?", courseID).Update("name", name).Update("title", title).Update("desc", desc).Error
}

func UpdateRelease(releaseID, desc string) error {
	return gormDB.Model(&Release{}).Where("id = ?", releaseID).Update("desc", desc).Error
}

func UpdatePost(postID, markdown string) error {
	return gormDB.Model(&Post{}).Where("id = ?", postID).Update("markdown", markdown).Error
}
