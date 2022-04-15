package db

import (
	"main/githubapi"
	"time"

	"github.com/google/uuid"
)

func CreateUser(user *User) error {
	return gormDB.Create(user).Error
}

func CreateCourse(course *Course) error {
	return gormDB.Create(course).Error
}

func CreateSession(userID uint64) (string, error) {
	session := Session{
		TokenUUID: uuid.NewString(),
		DeleteAt:  time.Now().Add(time.Hour * 100),
		UserID:    userID,
	}
	err := gormDB.Create(&session).Error
	return session.TokenUUID, err
}

func CreateRelease(release *Release) error {
	return gormDB.Create(release).Error
}

func CreateVersion(version *Version) error {
	return gormDB.Create(version).Error
}

func CreatePost(post *Post) error {
	return gormDB.Create(post).Error
}

func CreatePostToCourse(relation *PostToCourse) error {
	return gormDB.Create(relation).Error
}

func CreatePurchase(purchase *Purchase) error {
	return gormDB.Create(purchase).Error
}

func CreateBuyRelease(attemptBuyRelease *AttemptBuyRelease) error {
	return gormDB.Create(attemptBuyRelease).Error
}

func CreateStripeConnection(stripeConnection *StripeConnection) error {
	return gormDB.Create(stripeConnection).Error
}

func CreateVerify(verify *Verify) error {
	return gormDB.Create(verify).Error
}

func CreateGithubConnection(githubConnection *githubapi.GithubConnection) error {
	return gormDB.Create(githubConnection).Error
}

func CreateGithubRelease(githubRelease *GithubRelease) error {
	return gormDB.Create(githubRelease).Error
}

func CreateGithubVersion(githubVersion *GithubVersion) error {
	return gormDB.Create(githubVersion).Error
}

func CreateComment(comment *Comment) error {
	return gormDB.Create(comment).Error
}

func CreateChannel(channel *Channel) error {
	return gormDB.Create(channel).Error
}

func CreateMessage(message *Message) error {
	return gormDB.Create(&message).Error
}
