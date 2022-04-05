package db

import (
	"html/template"
	"time"
)

func migrate() {
	gormDB.AutoMigrate(
		// auth
		&User{},
		&Session{},
		&Verify{},

		// course
		&Course{},
		&Release{},
		&Version{},
		&Section{},
		&Content{},
		&Media{},
		&MediaChunk{},

		// purchases
		&Purchase{},
		&AttemptBuyRelease{},

		// third party apps
		&StripeConnection{},
		&GithubConnection{},

		// posts
		&Post{},
		&PostToRelease{},
	)
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

	Desc string

	// the purchases courseID
	// course ID
	CourseID uint64 `gorm:"not null"`

	// a specific course release
	ReleaseID uint64

	// not a required parameter but used to keep track of version user is currently taking
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

	Verified bool `gorm:"not null, default:f"`

	Purchases        []Purchase `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	StripeConnection StripeConnection
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
	UserID          uint64 `gorm:"not null,unique"`
}

// github access tokens never expire
type GithubConnection struct {
	UserID uint64 `gorm:"not null"`

	// the token for accessing the users github repos etc
	AccessToken string

	// the type of token
	TokenType string
}

type Session struct {
	TokenUUID string `gorm:"primaryKey, unique"` // this is the session id
	DeleteAt  time.Time
	UserID    uint64 `gorm:"not null"`
}

/* COURSE */
type Course struct {
	ID       uint64 `gorm:"primaryKey"`
	Title    string `gorm:"not null"` // a short title of the course
	Name     string `gorm:"not null"` // the courses unique url name (eg. spark.com/username/minecraftcourse)
	Subtitle string

	UserID uint64 `gorm:"not null"`

	// ORM preloadable property
	User User

	Releases []Release `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Versions []Version `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
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
	ID       uint64 `gorm:"primaryKey"`
	Price    uint64 `gorm:"default:0"`
	Num      uint16 `gorm:"default:0"`
	Markdown template.HTML
	CourseID uint64 `gorm:"not null"`
	Public   bool   `gorm:"default:f"`
	Level    uint32 `gorm:"default:0; not null"`

	UsingGithub   bool          `gorm:"default:f"` // defaults to false
	ReleaseGithub ReleaseGithub // githbu repo info

	Versions  []Version  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Purchases []Purchase `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type ReleaseGithub struct {
	ReleaseID   uint64
	TrackLatest bool
	// TODO needed github repo info
}

// points to a parent course
type Hierarchy struct {
	ReleaseID    uint64
	NextCourseID uint64
}

type Version struct {
	ID        uint64 `gorm:"primaryKey"`
	Num       uint16
	Patch     uint16 `gorm:"default:0"`
	CourseID  uint64 `gorm:"not null"`
	ReleaseID uint64 `gorm:"not null"`

	Error    string    `gorm:"default:null"`
	Sections []Section `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	// if using github repo
	CommitHash string
}

type Section struct {
	ID   uint64 `gorm:"primaryKey"`
	Name string

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

// media files for the course
type Media struct {
	ID        uint64 `gorm:"primaryKey"`
	VersionID uint64 `gorm:"not null"`
	Name      string
	Length    uint32
	Type      string // the "type" of file (.zip .png .gif)

	MediaChunks []MediaChunk `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type MediaChunk struct {
	MediaID uint64 `gorm:"not null"`
	// the raw bytes from the file
	Data []byte
	// Order to load chunks from db
	Position uint16
}

/* SOCIAL */
type Post struct {
	ID        uint64    `gorm:"primaryKey"`
	CreatedAt time.Time // special param name gorm automaically sets time
	UpdatedAt time.Time // special param name gorm automaically sets time
	UserID    uint64    `gorm:"not null"`
	Markdown  string

	User User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type UserMedia struct {
	ID        uint64 `gorm:"primaryKey"`
	UserID    uint64
	Name      string
	Length    uint32
	Type      string // the "type" of file (.zip .png .gif)
	CreatedAt time.Time

	UserMediaChunks []UserMediaChunk `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type UserMediaChunk struct {
	UserMediaID uint64 `gorm:"not null"`
	Data        []byte
	Position    uint16
}

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
// relate posts to a course release and section
type PostToRelease struct {
	PostID    uint64 `gorm:"not null"`
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
