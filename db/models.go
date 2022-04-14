package db

import (
	"html/template"
	"main/githubapi"
	"time"
)

func migrate() {
	gormDB.AutoMigrate(
		// auth
		&User{},
		&Session{},
		&Verify{},
		&Notif{},

		// course
		&Course{},
		&Release{},
		&Version{},
		&Section{},
		&Content{},
		&Channel{},

		// github based courses
		&GithubRelease{},
		&GithubVersion{},

		// purchases
		&Purchase{},
		&AttemptBuyRelease{},

		// third party apps
		&StripeConnection{},
		&githubapi.GithubConnection{},

		// posts
		&Post{},
		&Comment{},
		&PostToCourse{},
	)

	// get rid of old outdated tables
	gormDB.Migrator().DropTable(&PostToRelease{})
	gormDB.Migrator().DropTable("media")
	gormDB.Migrator().DropTable("media_chunks")
}

type AttemptBuyRelease struct {
	StripeSessionID       string `gorm:"primaryKey"`
	StripePaymentIntentID string
	ReleaseID             uint64
	UserID                uint64
	AmountPaying          uint64
	ExpiresAt             time.Time
}

type Purchase struct {
	ID     uint64 `gorm:"primaryKey"`
	UserID uint64 `gorm:"not null"`

	// id of the successful payment session
	StripeSessionID string

	// id of successful payment
	StripePaymentIntentID string

	CreatedAt  time.Time `gorm:"not null"`
	AmountPaid uint64    `gorm:"default:0"`
	AuthorsCut uint64    `gorm:"default:0"`

	Desc string // description or error or reason for purchase (eg. given to user for free since author payout failed)

	// the purchases courseID
	// course ID
	CourseID uint64 `gorm:"not null"`

	// a specific course release
	ReleaseID uint64

	// not a required parameter but used to keep track of version user is currently taking
	// then in the "home" page of the website for logged in users we can show them their courses, and take them to the version they are currently taking
	// also set to newest version when user first buys a course
	VersionID uint64

	// Preloading the user (from UserID) who owns this pruchase
	User User // don't add tag for cascading on delete cause it will delete the user when trying to delete the purchase
}

type User struct {
	ID       uint64 `gorm:"primaryKey"`
	Username string `gorm:"unique"` // unique identifer used in the url
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
	Public   bool `gorm:"default:f"`

	UserID uint64 `gorm:"not null"`

	// ORM preloadable property
	User User

	Channels []Channel `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // channels for the courses chat
	Releases []Release `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // course releases
	Versions []Version `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // course versions (have a parent release)
}

type Coupon struct {
	StripeCouponID string `gorm:"primaryKey"`
	CourseID       uint64
	// TODO
	/*
		- CreateCoupon(ExpiresBy)
		- Use methods to check coupons abaility
			- GetNumAvailable
			- Available
			- GetStripeData
	*/
}

type Release struct {
	ID             uint64 `gorm:"primaryKey"`
	Price          uint64 `gorm:"default:0"`
	Num            uint16 `gorm:"default:0"`
	Markdown       template.HTML
	CourseID       uint64 `gorm:"not null"`
	Public         bool   `gorm:"default:f"`
	Level          uint32 `gorm:"default:0; not null"`
	PostsNeededNum uint16 `gorm:"default:2;"`

	GithubRelease GithubRelease `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // githbu repo info

	Versions  []Version  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Purchases []Purchase `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type GithubRelease struct {
	ReleaseID uint64 `gorm:"not null"`
	RepoID    int64  `gorm:"not null"`
	RepoName  string `gorm:"default:woops;"`
	Branch    string `gorm:"default:main; not null"`
}

// points to a parent course
type Hierarchy struct {
	ReleaseID    uint64
	NextCourseID uint64
}

type Version struct {
	ID        uint64 `gorm:"primaryKey"`
	Num       uint16
	Patch     uint16 `gorm:"not null; default:0"`
	CourseID  uint64 `gorm:"not null"`
	ReleaseID uint64 `gorm:"not null"`
	CreatedAt time.Time
	Preview   bool `gorm:"default:f"`

	// relation for github based versions
	GithubVersion GithubVersion `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	// if using github version instead of manual upload with sections
	UsingGithub bool `gorm:"default:f;"`

	// if using manual uploading option with sections
	Error    string    `gorm:"default:null"`
	Sections []Section `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type GithubVersion struct {
	VersionID uint64 `gorm:"not null; primaryKey"`
	RepoID    int64  `gorm:"not null"`
	SHA       string `gorm:"not null"`
	Branch    string `gorm:"not null"`
	UpdatedAt time.Time
}

type Section struct {
	ID        uint64 `gorm:"primaryKey"`
	Name      string
	UpdatedAt time.Time

	// version this section is connected to
	VersionID uint64 `gorm:"not null"`
	ParentID  uint64 `gorm:"not null"` // parent section ID

	// children contents
	// special ORM parameter that can be preloaded with data
	Contents []Content `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
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
	Markdown  string

	User     User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Comments []Comment `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
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

	Course   Course    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Messages []Message `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type Message struct {
	ID        uint64 `gorm:"primaryKey"`
	UserID    uint64 `gorm:"not null"`
	ChannelID uint64 `gorm:"not null"`
	CreatedAt time.Time
	Markdown  string
}

/* RELATIONS */
type PostToRelease struct {
	PostID    uint64 `gorm:"not null"`
	ReleaseID uint64 `gorm:"not null"`
}

// relate posts to a course release and section
type PostToCourse struct {
	PostID    uint64 `gorm:"not null"`
	CourseID  uint64 `gorm:"not null"`
	ReleaseID uint64 `gorm:"not null"`
}

// maybe?
// allow for "thread-like" conversations to continue from messages?
type Thread struct {
	// the new channel ID
	ChannelID uint64 `gorm:"not null"`

	// the parent message
	ParentMessageID uint64 `gorm:"not null"`
	// the parent channel ID
	ParentChannelID uint64 `gorm:"not null"`
}
