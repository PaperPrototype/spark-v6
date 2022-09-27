package db

import (
	"log"
	"time"
)

func DeleteExpiredSessions() error {
	// when delete_at is before time.Now()
	return GormDB.Delete(&Session{}, "delete_at < ?", time.Now()).Error
}

func DeleteSession(TokenUUID string) error {
	return GormDB.Delete(&Session{}, "token_uuid = ?", TokenUUID).Error
}

func DeleteSection(sectionID string) error {
	return GormDB.Delete(&Section{}, "id = ?", sectionID).Error
}

func DeleteVersion(versionID string) error {
	// delete sections
	err2 := DeleteVersionSections(versionID)
	if err2 != nil {
		return err2
	}

	// delete version
	err := GormDB.Where("id = ?", versionID).Delete(&Version{}).Error
	if err != nil {
		log.Println("db ERROR deleting Version:", err)
		return err
	}

	return nil
}

func DeleteVersionSections(versionID string) error {
	section := Section{}
	err := GormDB.Where("id = ?", versionID).Delete(&section).Error
	if err != nil {
		log.Println("db ERROR deleting Version Section:", err)
		return err
	}

	// delte contents of section
	err1 := GormDB.Where("section_id = ?", section.ID).Delete(&Content{}).Error
	if err1 != nil {
		log.Println("db ERROR deleting Version Section:", err1)
		return err1
	}

	return nil
}

func DeleteBuyRelease(stripeSessionID string) error {
	return GormDB.Model(&AttemptBuyRelease{}).Where("stripe_session_id = ?", stripeSessionID).Update("expires_at = ?", time.Now()).Error
}

func DeleteRelease(releaseID string) error {
	err := GormDB.Where("id = ?", releaseID).Delete(&Release{}).Error
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
	return GormDB.Model(&Version{}).Where("release_id = ?", releaseID).Delete(&Version{}).Error
}

func DeleteExpiredVerify() error {
	return GormDB.Delete(&Verify{}, "expires_at < ?", time.Now()).Error
}

func DeletePrerequisite(preqID interface{}) error {
	return GormDB.Delete(&Prerequisite{}, "id = ?", preqID).Error
}
