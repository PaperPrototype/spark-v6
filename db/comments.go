package db

func GetComments(postID string, limit int) ([]Comment, int64, error) {
	comments := []Comment{}
	var count int64

	err := gormDB.Model(&Comment{}).Where("post_id = ?", postID).Count(&count).Error
	// if there was an error
	if err != nil {
		return comments, count, err
	}

	// if no comments
	if count == 0 {
		return comments, count, nil
	}

	err1 := gormDB.Model(&Comment{}).Where("post_id = ?", postID).Preload("User").Limit(limit).Find(&comments).Error

	return comments, count, err1
}

func GetNewComments(postID string, newestCommentDate string) ([]Comment, int64, error) {
	comments := []Comment{}
	var count int64

	err := gormDB.Model(&Comment{}).Where("post_id = ?", postID).Where("created_at > ?", newestCommentDate).Count(&count).Error
	if err != nil {
		return comments, count, err
	}

	// if no comments
	if count == 0 {
		return comments, count, nil
	}

	err1 := gormDB.Model(&Comment{}).Where("post_id = ?", postID).Where("created_at > ?", newestCommentDate).Preload("User").Find(&comments).Error

	return comments, count, err1
}
