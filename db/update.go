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
