package db

import (
	"html/template"
	"time"
)

func migrate() {
	// updates to database for v7

	GormDB.AutoMigrate(
		// auth
		&User{},
		&Session{},
		&Verify{},
		&Notif{},

		// course
		&Course{},
		&Release{},
		&Section{},
		&GithubSection{},
		// &Content{},

		// chat
		&Channel{},
		&Message{},

		// hierarchy
		&Prerequisite{},

		// github based courses
		&GithubRelease{},

		// TODO DELETE
		&GithubVersion{},

		// purchases
		&Purchase{},
		&AttemptBuyRelease{},

		// course ownership and progress
		&Ownership{},
		&FinishedSection{},

		// third party apps
		&StripeConnection{},
		&GithubConnection{},

		// posts
		&Post{},
		&Comment{},
		&PostToCourse{},
		&PostToCourseReview{},
	)
}

// github access tokens never expire
type GithubConnection struct {
	UserID uint64 `gorm:"not null"`

	// the token for accessing the users github repos etc
	AccessToken string

	// the type of token
	TokenType string
}

// temporary payment attempt
type AttemptBuyRelease struct {
	StripeSessionID string `gorm:"primaryKey"`
	StripePaymentID string // the id of the payment intent
	ReleaseID       uint64
	UserID          uint64
	AmountPaying    uint64
	ExpiresAt       time.Time // WE DON'T DELETE the buy release if it is expired. It just gets filtered out
}

// the user has purchased the course
type Purchase struct {
	ID     uint64 `gorm:"primaryKey"`
	UserID uint64 `gorm:"not null"` // the student buying the course

	// id of the successful payment session
	StripeSessionID string

	// id of successful payment
	StripePaymentID string

	CreatedAt  time.Time `gorm:"not null"`
	AmountPaid uint64    `gorm:"default:0"` // amount the student paid for the course
	AuthorsCut uint64    `gorm:"default:0"` // record of what was payed out to the author

	Desc string // description or error or reason for purchase (eg. given to user for free since author payout failed)

	// the purchases courseID
	// course ID
	CourseID uint64 `gorm:"not null"`

	// a specific course release
	ReleaseID uint64 `gorm:"not null"`

	// optional parameter but used to keep track of version purchased
	// then in the "home" page of the website for logged in users we can show them their courses, and take them to the version they are currently taking
	// also set to newest version when user first buys a course
	VersionID uint64

	// Preloading the user (from UserID) who owns this pruchase
	User User // don't add tag for cascading on delete cause it will delete the user when trying to delete the purchase
}

// course ownership and access to a course
// also caches course progress
type Ownership struct {
	ID        uint64
	UserID    uint64
	CourseID  uint64
	ReleaseID uint64
	VersionID uint64

	CreatedAt time.Time

	Desc string // describe any errors or info

	Completed   bool
	CompletedAt time.Time
	Progress    float32
	PostsCount  uint32

	User    User
	Course  Course
	Release Release
}

// existence of this signifies the
// user has completed this section
type FinishedSection struct {
	VersionID uint64
	SectionID string
}

// claimable coupon
type Coupon struct {
	ID           uint64
	OwnsCourseID uint64
	CourseID     uint64
	ReleaseID    uint64
	Public       bool // visible as a sale on course page vs private link the author will send to students

	CreatedAt uint64
	ExpiresAt uint64
	Available uint32
	Claimed   uint32
	// TODO
	/*
		- CreateCoupon(ExpiresBy)
		- Use methods to check coupons abaility
			- GetNumAvailable
			- Available
			- GetStripeData
	*/

	CouponClaims []CouponClaim
}

// coupon claims
type CouponClaim struct {
	ID        uint64
	CouponID  uint64
	UserID    uint64
	CourseID  uint64
	ReleaseID uint64
}

type User struct {
	ID       uint64 `gorm:"primaryKey"`
	Username string `gorm:"unique"` // unique for each user (identifer used in the url of courses as well)
	Name     string // real name
	Hash     string // password hash
	Email    string `gorm:"unique"`
	Bio      string

	CreatedAt time.Time `gorm:"not null"` // date when the account was created

	Verified bool `gorm:"not null; default:f"`

	Notifs    []Notif    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Purchases []Purchase `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// general purpose notification
type Notif struct {
	ID        uint64 `gorm:"primaryKey"`
	UserID    uint64
	Message   string
	URL       string
	CreatedAt time.Time

	// if the user has read the notification
	Read bool `gorm:"default:f"`
}

type Verify struct {
	UserID     uint64
	VerifyUUID string
	ExpiresAt  time.Time
}

// stripe connection
// existence of this means they can make courses
// stripe connection cannot be made until they verify their email
type StripeConnection struct {
	StripeAccountID string
	UserID          uint64 `gorm:"not null; unique"`
}

type Session struct {
	TokenUUID string `gorm:"primaryKey; unique"` // this is the session id
	DeleteAt  time.Time
	UserID    uint64 `gorm:"not null"`
}

/* COURSE */
type Course struct {
	ID       uint64 `gorm:"primaryKey"`
	Title    string `gorm:"not null"` // a short title of the course
	Name     string `gorm:"not null"` // the courses unique url name (eg. spark.com/username/minecraftcourse)
	Subtitle string
	Public   bool   `gorm:"default:f"`
	Markdown string // course landing page markdown

	Level uint32 `gorm:"default:0"`

	UserID uint64 `gorm:"not null"`

	// ORM preloadable property
	User    User
	Release Release // can be preloaded with the newest release

	Channels      []Channel      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`  // channels for the courses chat
	Releases      []Release      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`  // course releases
	Prerequisites []Prerequisite `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE; "` // prerequisites for this course
}

type Release struct {
	ID             uint64 `gorm:"primaryKey"`
	Price          uint64 `gorm:"default:0"`
	Num            uint16 `gorm:"default:0"`
	Markdown       template.HTML
	CourseID       uint64 `gorm:"not null"`
	Public         bool   `gorm:"default:f"`
	PostsNeededNum uint16 `gorm:"default:2;"`
	CreatedAt      time.Time
	ImageURL       string        `gorm:"default:'';"`
	GithubEnabled  bool          `gorm:"default:f"`                                     // make this release use a github repo
	GithubRelease  GithubRelease `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // github repo info

	Purchases []Purchase `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type GithubRelease struct {
	ReleaseID uint64 `gorm:"not null"`
	RepoID    uint64 `gorm:"not null"`
	RepoName  string `gorm:"default:woops;"`
	Branch    string `gorm:"default:main; not null"`
	SHA       string // TODO use for tracking the SHA of the current git commit
	Patch     uint32 // number of times we updated the github repo details (used to check if section markdown cache is invalid and needs updated)
	UpdatedAt time.Time
}

type Version struct {
}

type GithubVersion struct {
}

type Section struct {
	ID          uint64 `gorm:"primaryKey"` // TODO convert to string UUID or sha
	Name        string
	ReleaseID   uint64
	Num         uint16 `gorm:"default:0"` // what order to put the sections in
	Free        bool   `gorm:"default:f"`
	Description string
	UpdatedAt   time.Time

	// delete the corresponding GithubSection when we delete its section
	GithubSection GithubSection `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type GithubSection struct {
	SectionID             uint64 `gorm:"unique"`
	Path                  string // sha for specific course section
	MarkdownCache         string // cache
	MarkdownCachePatchNum uint32 // what patch number was on the release when we cached
}

type Content struct {
	ID        uint64 `gorm:"primaryKey"`
	Language  string
	SectionID uint64 `gorm:"not null"`
	Markdown  string
}

/* SOCIAL */
type Post struct {
	ID        uint64    `gorm:"primaryKey"`
	CreatedAt time.Time // special param name gorm automaically sets time
	UpdatedAt time.Time // special param name gorm automaically sets time
	UserID    uint64    `gorm:"not null"`
	Title     string
	Markdown  string

	// each Like has a db hook that runs a db query to update its posts likes count
	LikesCount uint64 // cache of the number of likes in a post

	User User

	// delete the likes when we delete the post
	Likes []Like `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	// delete the comments when we delete the post
	Comments []Comment `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type Like struct {
	UserID uint64
	PostID uint64
}

// post comments
// purposefully has no ID since the ID will max out very quickly
// we don't want to waste IO on massive
type Comment struct {
	PostID    uint64
	Markdown  string
	UserID    uint64
	CreatedAt time.Time

	// preloads
	User User
}

// a channel of messages for a course
type Channel struct {
	ID       uint64 `gorm:"primaryKey"`
	CourseID uint64 `gorm:"not null"`
	Name     string

	Course Course
	// 					when we delete the channel the messages will also delete
	Messages []Message `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type Message struct {
	ID        uint64 `gorm:"primaryKey"`
	UserID    uint64 `gorm:"not null"`
	ChannelID uint64 `gorm:"not null"`
	CreatedAt time.Time
	Markdown  string

	// preloads
	User User
}

/* RELATIONS */
// relate posts to a course release and section
type PostToCourse struct {
	PostID    uint64 `gorm:"not null"`
	CourseID  uint64 `gorm:"not null"`
	ReleaseID uint64 `gorm:"not null"`
}

// points to a parent course
type Prerequisite struct {
	ID                   uint64 `gorm:"primaryKey"`
	CourseID             uint64
	PrerequisiteCourseID uint64

	PrerequisiteCourse Course `gorm:"foreignkey:PrerequisiteCourseID;"`
}

// reviews on courses
type PostToCourseReview struct {
	ID        uint64    `gorm:"not null"`
	CourseID  uint64    `gorm:"not null"`
	ReleaseID uint64    `gorm:"not null"`
	PostID    uint64    `gorm:"not null"`
	Rating    uint8     `gorm:"not null"` // range of 0 to 5
	UserID    uint64    `gorm:"not null"` // the person who posted the review
	CreatedAt time.Time `gorm:"not null"`

	// preloadable properties given foreign keys
	User    User
	Post    Post
	Release Release
}

// thread system ideas:
// allow for "thread-like" conversations to continue from messages?
// allow for child channels to be created?
type ThreadTODO struct {
	// the new channel ID
	ChannelID uint64 `gorm:"not null"`

	// the parent message
	ParentMessageID uint64 `gorm:"not null"`
	// the parent channel ID
	ParentChannelID uint64 `gorm:"not null"`
}

// category ideas:
// allow for channels to be grouped (like in discord)
type ChannelCategoryTODO struct {
}
