package db

// get up to `num` comments
func GetComments(postID string, num int64) ([]Comment, int64, error) {
	comments := []Comment{}
	var count int64

	err := GormDB.Model(&Comment{}).Where("post_id = ?", postID).Count(&count).Error
	// if there was an error
	if err != nil {
		return comments, count, err
	}

	// if no comments
	if count == 0 {
		return comments, count, nil
	}

	err1 := GormDB.Model(&Comment{}).Where("post_id = ?", postID).Preload("User").Limit(int(num)).Offset(int(count - num)).Order("created_at ASC").Find(&comments).Error

	return comments, count, err1
}

func GetNewComments(postID string, newestCommentDate string) ([]Comment, int64, error) {
	comments := []Comment{}
	var count int64

	err := GormDB.Model(&Comment{}).Where("post_id = ?", postID).Where("created_at > ?", newestCommentDate).Count(&count).Error
	if err != nil {
		return comments, count, err
	}

	// if no comments
	if count == 0 {
		return comments, count, nil
	}

	err1 := GormDB.Model(&Comment{}).Where("post_id = ?", postID).Where("created_at > ?", newestCommentDate).Preload("User").Order("created_at ASC").Find(&comments).Error

	return comments, count, err1
}
