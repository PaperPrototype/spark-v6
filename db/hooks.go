package db

import (
	"log"

	"gorm.io/gorm"
)

// update a posts likes_count
func (like *Like) AfterSave(tx *gorm.DB) error {
	var count int64
	err := gormDB.Model(&Like{}).Where("post_id = ?", like.PostID).Count(&count).Error
	if err != nil {
		log.Println("db/hooks ERROR counting likes for post:", err)
		return err
	}

	err1 := gormDB.Model(&Post{}).Where("post_id = ?", like.PostID).Update("likes_count", count).Error
	if err1 != nil {
		log.Println("db/hooks ERROR counting likes for post:", err1)
		return err1
	}

	return nil
}

// update posts likes when a delete a like
func (like *Like) AfterDelete(tx *gorm.DB) error {
	var count int64
	err := gormDB.Model(&Like{}).Where("post_id = ?", like.PostID).Count(&count).Error
	if err != nil {
		log.Println("db/hooks ERROR counting likes for post:", err)
		return err
	}

	err1 := gormDB.Model(&Post{}).Where("post_id = ?", like.PostID).Update("likes_count", count).Error
	if err1 != nil {
		log.Println("db/hooks ERROR counting likes for post:", err1)
		return err1
	}

	return nil
}
