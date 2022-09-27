package db

// channels hold messages and are the course chat

// get up to `num` messages
func GetMessages(channelID string, num int64) ([]Message, int64, error) {
	messages := []Message{}
	var count int64

	err := GormDB.Model(&Message{}).Where("channel_id = ?", channelID).Count(&count).Error
	// if there was an error
	if err != nil {
		return messages, count, err
	}

	// if no messages
	if count == 0 {
		return messages, count, nil
	}

	err1 := GormDB.Model(&Message{}).Where("channel_id = ?", channelID).Preload("User").Limit(int(num)).Offset(int(count - num)).Order("created_at ASC").Find(&messages).Error

	return messages, count, err1
}

func GetNewMessages(channelID string, newestMessageDate string) ([]Message, int64, error) {
	messages := []Message{}
	var count int64

	err := GormDB.Model(&Message{}).Where("channel_id = ?", channelID).Where("created_at > ?", newestMessageDate).Count(&count).Error
	if err != nil {
		return messages, count, err
	}

	// if no messages
	if count == 0 {
		return messages, count, nil
	}

	err1 := GormDB.Model(&Message{}).Where("channel_id = ?", channelID).Where("created_at > ?", newestMessageDate).Preload("User").Order("created_at ASC").Find(&messages).Error

	return messages, count, err1
}

func GetChannels(courseID uint64) ([]Channel, error) {
	channels := []Channel{}
	err := GormDB.Model(&Channel{}).Where("course_id = ?", courseID).Find(&channels).Error
	return channels, err
}

func GetChannel(channelID string) (*Channel, error) {
	channel := Channel{}
	err := GormDB.Model(&Channel{}).Where("id = ?", channelID).First(&channel).Error
	return &channel, err
}
