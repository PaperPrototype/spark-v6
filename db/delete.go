package db

import (
	"log"
	"time"
)

func DeleteExpiredSessions() error {
	// when delete_at is before time.Now()
	return gormDB.Delete(&Session{}, "delete_at < ?", time.Now()).Error
}

func DeleteSession(TokenUUID string) error {
	return gormDB.Delete(&Session{}, "token_uuid = ?", TokenUUID).Error
}

func DeleteVersion(versionID string) error {
	// delete media first
	err1 := DeleteVersionMedias(versionID)
	if err1 != nil {
		return err1
	}

	// delete sections
	err2 := DeleteVersionSections(versionID)
	if err2 != nil {
		return err2
	}

	// delete version
	err := gormDB.Where("id = ?", versionID).Delete(&Version{}).Error
	if err != nil {
		log.Println("db ERROR deleting Version:", err)
		return err
	}

	return nil
}

func DeleteVersionSections(versionID string) error {
	section := Section{}
	err := gormDB.Where("id = ?", versionID).Delete(&section).Error
	if err != nil {
		log.Println("db ERROR deleting Version Section:", err)
		return err
	}

	// delte contents of section
	err1 := gormDB.Where("section_id = ?", section.ID).Delete(&Content{}).Error
	if err1 != nil {
		log.Println("db ERROR deleting Version Section:", err1)
		return err1
	}

	return nil
}

func DeleteVersionMedias(versionID string) error {
	media := Media{}
	err := gormDB.Where("version_id = ?", versionID).Delete(&media).Error
	if err != nil {
		log.Println("db ERROR deleting Version Media:", err)
		return err
	}

	// also delete media chunks
	err1 := gormDB.Where("media_id = ?", media.ID).Model(&MediaChunk{}).Error
	if err1 != nil {
		log.Println("db ERROR deleting Version MediaChunk:", err1)
		return err1
	}

	return nil
}

func DeleteBuyRelease(buyReleaseID string) error {
	return gormDB.Where("stripe_session_id = ?", buyReleaseID).Delete(&BuyRelease{}).Error
}

func DeleteRelease(releaseID string) error {
	err := gormDB.Where("id = ?", releaseID).Delete(&Release{}).Error
	if err != nil {
		return err
	}

	err1 := DeleteReleaseVersions(releaseID)
	if err1 != nil {
		return err1
	}

	return nil
}

func DeleteReleaseVersions(releaseID string) error {
	return gormDB.Model(&Version{}).Where("release_id = ?", releaseID).Delete(&Version{}).Error
}
