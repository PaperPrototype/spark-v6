package db

import "log"

func GetUnreadNotifs(userID uint64, num int64) ([]Notif, int64, error) {
	notifs := []Notif{}
	var count int64

	err := gormDB.Model(&Notif{}).Where("user_id = ?", userID).Where("read", false).Count(&count).Error
	// if there was an error
	if err != nil {
		return notifs, count, err
	}

	// if no notifications
	if count == 0 {
		return notifs, count, nil
	}

	err1 := gormDB.Model(&Notif{}).Where("user_id = ?", userID).Where("read", false).Limit(int(num)).Offset(int(count - num)).Order("created_at ASC").Find(&notifs).Error
	return notifs, count, err1
}

func GetNewUnreadNotifs(userID uint64, newestDate string) ([]Notif, int64, error) {
	notifs := []Notif{}
	var count int64

	err := gormDB.Model(&Notif{}).Where("user_id = ?", userID).Where("created_at > ?", newestDate).Where("read", false).Count(&count).Error
	if err != nil {
		return notifs, count, err
	}

	// if no comments
	if count == 0 {
		return notifs, count, nil
	}

	err1 := gormDB.Model(&Notif{}).Where("user_id = ?", userID).Where("created_at > ?", newestDate).Where("read", false).Order("created_at ASC").Find(&notifs).Error

	return notifs, count, err1
}

func NotifyUsers(usernames []string, message string, url string) error {
	var userIDs []uint64

	err := gormDB.Model(&User{}).Where("username IN (?)", usernames).Pluck("id", &userIDs).Error
	if err != nil {
		log.Println("db/utils ERROR getting userIDs in NotifyAllUsers:", err)
		return err
	}

	for _, id := range userIDs {
		notif := Notif{
			UserID:  id,
			Message: message,
			URL:     url,
		}
		err1 := gormDB.Create(&notif).Error
		if err1 != nil {
			log.Println("db/utils ERROR creating notif in NotifyAllUsers:", err1)
			return err1
		}
	}

	return nil
}

func NotifSetRead(notifID string) error {
	return gormDB.Model(&Notif{}).Where("id = ?", notifID).Update("read", true).Error
}
