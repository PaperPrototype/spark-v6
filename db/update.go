package db

func UpdateUser(userID uint64, username, name, bio, email string) error {
	return GormDB.Model(&User{}).Where("id = ?", userID).Update("bio", bio).Update("name", name).Update("username", username).Update("email", email).Error
}

func UpdateCourse(courseID, title, name, subtitle string, public bool) error {
	return GormDB.Model(&Course{}).Where("id = ?", courseID).Update("name", name).Update("title", title).Update("subtitle", subtitle).Update("public", public).Error
}

func UpdateRelease(releaseID string, price uint64, public bool, postsNeededNum uint16, imageURL string, githubEnabled bool) error {
	return GormDB.Model(&Release{}).Where("id = ?", releaseID).Update("price", price).Update("public", public).Update("posts_needed_num", postsNeededNum).Update("image_url", imageURL).Update("github_enabled", githubEnabled).Error
}

func UpdatePost(postID, title, markdown string) error {
	return GormDB.Model(&Post{}).Where("id = ?", postID).Update("title", title).Update("markdown", markdown).Error
}

func UpdateSection(sectionID, name, desc string, isFree bool, num uint16) error {
	return GormDB.Model(&Section{}).Where("id = ?", sectionID).Update("name", name).Update("description", desc).Update("free", isFree).Update("num", num).Error
}

func UpdateGithubRelease(releaseID uint64, branch string, repoID int64, repoName string) error {
	return GormDB.Model(&GithubRelease{}).Where("release_id = ?", releaseID).Update("branch", branch).Update("repo_id", repoID).Update("repo_name", repoName).Error
}

func UpdateGithubSectionMarkdownCache(sectionID uint64, markdownCache string) error {
	return GormDB.Model(&GithubSection{}).Where("section_id = ?", sectionID).Update("markdown_cache", markdownCache).Error
}

func UpdateGithubSection(sectionID string, path string, markdownCache string) error {
	return GormDB.Model(&GithubSection{}).Where("section_id = ?", sectionID).Update("path", path).Update("markdown_cache", markdownCache).Error
}
