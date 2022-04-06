package db

func UpdateUser(userID uint64, username, name, bio, email string) error {
	return gormDB.Model(&User{}).Where("id = ?", userID).Update("bio", bio).Update("name", name).Update("username", username).Update("email", email).Error
}

func UpdateCourse(courseID, title, name, subtitle string) error {
	return gormDB.Model(&Course{}).Where("id = ?", courseID).Update("name", name).Update("title", title).Update("subtitle", subtitle).Error
}

func UpdateRelease(releaseID string, markdown string, price uint64, public bool) error {
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

func UpdateGithubRelease(releaseID uint64, branch string, repoID int64, repoName string) error {
	return gormDB.Model(&GithubRelease{}).Where("release_id = ?", releaseID).Update("branch", branch).Update("repo_id", repoID).Update("repo_name", repoName).Error
}
