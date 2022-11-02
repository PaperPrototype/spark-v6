package db

import (
	"time"

	"github.com/google/uuid"
)

func CreateUser(user *User) error {
	return GormDB.Create(user).Error
}

func CreateCourse(course *Course) error {
	return GormDB.Create(course).Error
}

func CreateSession(userID uint64) (string, error) {
	session := Session{
		TokenUUID: uuid.NewString(),
		DeleteAt:  time.Now().Add(time.Hour * 100),
		UserID:    userID,
	}
	err := GormDB.Create(&session).Error
	return session.TokenUUID, err
}

func CreateRelease(release *Release) error {
	return GormDB.Create(release).Error
}

func CreateVersion(version *Version) error {
	return GormDB.Create(version).Error
}

func CreatePost(post *Post) error {
	return GormDB.Create(post).Error
}

func CreatePostToCourse(relation *PostToCourse) error {
	return GormDB.Create(relation).Error
}

func CreatePurchase(purchase *Purchase) error {
	return GormDB.Create(purchase).Error
}

func CreateBuyRelease(attemptBuyRelease *AttemptBuyRelease) error {
	return GormDB.Create(attemptBuyRelease).Error
}

func CreateStripeConnection(stripeConnection *StripeConnection) error {
	return GormDB.Create(stripeConnection).Error
}

func CreateVerify(verify *Verify) error {
	return GormDB.Create(verify).Error
}

func CreateGithubConnection(githubConnection *GithubConnection) error {
	return GormDB.Create(githubConnection).Error
}

func CreateOrUpdateGithubSection(sectionID string, githubSection *GithubSection) error {
	var count int64 = 0
	err := GormDB.Model(&GithubSection{}).Where("section_id = ?", sectionID).Count(&count).Error

	if err != nil {
		return err
	}

	var err1 error = nil
	if count > 0 {
		// update
		err1 = GormDB.Model(&GithubSection{}).Where("section_id = ?", sectionID).Update("path", githubSection.Path).Error

	} else {
		// create

		err1 = GormDB.Create(githubSection).Error // create new record

	}

	return err1
}

func CreateGithubSection(sectionID string, githubSection *GithubSection) error {
	return GormDB.Create(githubSection).Error // create new record
}

func CreateOrUpdateGithubRelease(releaseID string, githubRelease *GithubRelease) error {

	var count int64 = 0
	err := GormDB.Model(&GithubRelease{}).Where("release_id = ?", releaseID).Count(&count).Error

	if err != nil {
		return err
	}

	var err1 error = nil
	if count > 0 {
		err1 = GormDB.Model(&GithubRelease{}).Where("release_id = ?", releaseID).Update("repo_id", githubRelease.RepoID).Update("repo_name", githubRelease.RepoName).Update("branch", githubRelease.Branch).Update("sha", githubRelease.SHA).Error
	} else {
		err1 = GormDB.Create(githubRelease).Error // create new record
	}

	return err1
}

func CreateGithubRelease(githubRelease *GithubRelease) error {
	return GormDB.Create(githubRelease).Error
}

func CreateGithubVersion(githubVersion *GithubVersion) error {
	return GormDB.Create(githubVersion).Error
}

func CreateSection(section *Section) error {
	return GormDB.Create(section).Error
}

func CreateComment(comment *Comment) error {
	return GormDB.Create(comment).Error
}

func CreateChannel(channel *Channel) error {
	return GormDB.Create(channel).Error
}

func CreateMessage(message *Message) error {
	return GormDB.Create(&message).Error
}

func CreateReview(review *PostToCourseReview) error {
	return GormDB.Create(&review).Error
}

func CreatePrerequisite(preq *Prerequisite) error {
	return GormDB.Create(preq).Error
}

func CreateOwnership(ownership *Ownership) error {
	return GormDB.Create(ownership).Error
}
