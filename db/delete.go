package db

import "time"

func DeleteExpiredSessions() error {
	// when delete_at is before time.Now()
	return gormDB.Delete(&Session{}, "delete_at < ?", time.Now()).Error
}

func DeleteSession(TokenUUID string) error {
	return gormDB.Delete(&Session{}, "token_uuid = ?", TokenUUID).Error
}
